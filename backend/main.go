package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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

	// Data directories based on our architecture
	// When running via Docker, these paths might just be /data/...
	// For local dev, we assume we are running within the `backend` folder
	expertsDir := filepath.Join("..", "data", "experts")
	historyDir := filepath.Join("..", "data", "history")
	
	os.MkdirAll(expertsDir, 0755)
	os.MkdirAll(historyDir, 0755)

	api := app.Group("/api")
	api.Post("/backtest", handleUpload)

	// Serve static files from Svelte Dist
	app.Static("/", filepath.Join("..", "frontend", "dist"))

	// SPA fallback
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
	// 1. Get uploaded files
	ex4File, err := c.FormFile("ex4")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing EX4 file"})
	}

	csvFile, err := c.FormFile("history")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing CSV history file"})
	}

	// 2. Save EX4
	expertsPath := filepath.Join("..", "data", "experts", ex4File.Filename)
	if err := c.SaveFile(ex4File, expertsPath); err != nil {
		log.Printf("Failed to save EX4: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save EX4 file"})
	}

	// 3. Save CSV
	historyPath := filepath.Join("..", "data", "history", csvFile.Filename)
	if err := c.SaveFile(csvFile, historyPath); err != nil {
		log.Printf("Failed to save CSV: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save CSV file"})
	}

	// Note: Push 5 will implement the Wine/Xvfb execution logic here

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Successfully processed %s and %s", ex4File.Filename, csvFile.Filename),
		"status":  "upload_complete_pending_execution",
	})
}
