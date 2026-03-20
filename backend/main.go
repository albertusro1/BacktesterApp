package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024, // 100MB max limit
	})

	app.Use(logger.New())
	app.Use(recover.New())

	expertsDir := filepath.Join("..", "data", "experts")
	historyDir := filepath.Join("..", "data", "history")
	
	os.MkdirAll(expertsDir, 0755)
	os.MkdirAll(historyDir, 0755)

	api := app.Group("/api")
	api.Post("/backtest", handleUpload)

	app.Static("/", filepath.Join("..", "frontend", "dist"))

	app.Use(func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join("..", "frontend", "dist", "index.html"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	
	log.Printf("Starting backend server on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func handleUpload(c *fiber.Ctx) error {
	ex4File, err := c.FormFile("ex4")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing EX4 file"})
	}

	csvFile, err := c.FormFile("history")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing CSV history file"})
	}

	expertsPath := filepath.Join("..", "data", "experts", ex4File.Filename)
	if err := c.SaveFile(ex4File, expertsPath); err != nil {
		log.Printf("Failed to save EX4: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save EX4 file"})
	}

	historyPath := filepath.Join("..", "data", "history", csvFile.Filename)
	if err := c.SaveFile(csvFile, historyPath); err != nil {
		log.Printf("Failed to save CSV: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save CSV file"})
	}

	// 4. Generate INI config
	iniPath := filepath.Join("..", "data", "config.ini")
	if err := generateINI(iniPath, ex4File.Filename, csvFile.Filename); err != nil {
		log.Printf("Failed to generate INI: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate INI file"})
	}

	// 5. Execute MT4 via Wine and Xvfb
	log.Println("Starting MT4 backtest execution...")
	output, runErr := runBacktest(iniPath)
	if runErr != nil {
		log.Printf("Execution failed: %v\nOutput: %s", runErr, output)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Backtest execution failed",
			"details": string(output),
			"error_msg": runErr.Error(),
		})
	}

	// 6. Parse the HTM report
	possiblePaths := []string{
		filepath.Join("..", "data", "report.htm"),
		filepath.Join("..", "data", "MT4", "report.htm"),
		"report.htm",
	}

	var reportData map[string]string
	var parseErr error
	parsed := false
	
	for _, p := range possiblePaths {
		reportData, parseErr = parseReport(p)
		if parseErr == nil {
			parsed = true
			break
		}
	}

	if !parsed {
		log.Printf("Failed to find/parse report")
		return c.JSON(fiber.Map{
			"message": "Backtest executed, but report generating/parsing failed.",
			"raw_output": string(output),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Backtest completed successfully",
		"status":  "success",
		"results": reportData,
		"raw_output": string(output),
	})
}

func generateINI(iniPath, ex4Name, csvName string) error {
	// A basic MT4 tester configuration template
	// Stripping .ex4 extension for the INI file
	expertName := ex4Name
	if len(expertName) > 4 && expertName[len(expertName)-4:] == ".ex4" {
		expertName = expertName[:len(expertName)-4]
	}

	config := fmt.Sprintf(`[Tester]
Expert=%s
Symbol=EURUSD
Period=H1
Deposit=10000
Model=0
Optimization=0
FromDate=2023.01.01
ToDate=2024.01.01
Report=report.htm
ReplaceReport=1
ShutdownTerminal=1
`, expertName)

	return os.WriteFile(iniPath, []byte(config), 0644)
}

func runBacktest(iniPath string) ([]byte, error) {
	// Target Docker environment architecture execution:
	// xvfb-run -a wine /path/to/terminal.exe /portable "config.ini"
	// We'll formulate a command assuming /data/MT4/terminal.exe
	
	absIniPath, err := filepath.Abs(iniPath)
	if err != nil {
		absIniPath = iniPath // fallback
	}

	// The command expects MT4 terminal in /data/MT4/terminal.exe inside the future container.
	cmd := exec.Command("xvfb-run", "-a", "wine", "/data/MT4/terminal.exe", "/portable", absIniPath)
	
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	
	err = cmd.Run()
	return out.Bytes(), err
}

func parseReport(reportPath string) (map[string]string, error) {
	content, err := os.ReadFile(reportPath)
	if err != nil {
		return nil, err
	}

	html := string(content)
	results := make(map[string]string)

	keys := []string{
		"Total net profit",
		"Gross profit",
		"Gross loss",
		"Profit factor",
		"Expected payoff",
		"Absolute drawdown",
		"Maximal drawdown",
		"Total trades",
	}

	for _, key := range keys {
		val := extractValue(html, key)
		if val != "" {
			results[key] = val
		}
	}

	return results, nil
}

func extractValue(html, key string) string {
	re := regexp.MustCompile(fmt.Sprintf(`(?is)>\s*%s\s*<.*?<td[^>]*>(.*?)</td>`, regexp.QuoteMeta(key)))
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		cleanText := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(match[1], "")
		return strings.TrimSpace(cleanText)
	}
	return ""
}
