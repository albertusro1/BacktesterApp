Write-Host "=== BacktesterApp Local Environment Setup ===" -ForegroundColor Cyan
Write-Host "Checking for required dependencies..."

$deps = @("docker", "go", "node", "npm", "git")
$missing = @()

foreach ($dep in $deps) {
    if (-not (Get-Command $dep -ErrorAction SilentlyContinue)) {
        $missing += $dep
    }
}

if ($missing.Count -gt 0) {
    Write-Host "Missing dependencies detected: $($missing -join ', ')" -ForegroundColor Yellow
    $installDir = Read-Host "Please enter the absolute path for the installation directory (or press Enter for default system paths)"

    Write-Host "Starting installation for missing dependencies..." -ForegroundColor Cyan

    if ([string]::IsNullOrWhiteSpace($installDir)) {
        Write-Host "Using default system paths. Attempting to install via winget..."
        foreach ($dep in $missing) {
            if ($dep -eq "docker") {
                winget install -e --id Docker.DockerDesktop --accept-package-agreements --accept-source-agreements
            }
            elseif ($dep -eq "go") {
                winget install -e --id GoLang.Go --accept-package-agreements --accept-source-agreements
            }
            elseif ($dep -eq "node" -or $dep -eq "npm") {
                winget install -e --id OpenJS.NodeJS --accept-package-agreements --accept-source-agreements
            }
            elseif ($dep -eq "git") {
                winget install -e --id Git.Git --accept-package-agreements --accept-source-agreements
            }
        }
    } else {
        Write-Host "Custom installation to $installDir is requested."
        New-Item -ItemType Directory -Force -Path $installDir | Out-Null
        Write-Host "Please note: Fully automated custom-path installation of Docker/Go on Windows requires downloading specific zip/installers."
        Write-Host "Directory $installDir created/verified. Place your portable binaries here, and ensure it is in your PATH."
    }
} else {
    Write-Host "All required dependencies (Docker, Golang, Node.js, npm, Git) are installed!" -ForegroundColor Green
}

Write-Host "Local environment bootstrap complete." -ForegroundColor Cyan
