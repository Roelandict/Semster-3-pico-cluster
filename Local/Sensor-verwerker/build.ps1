# Build script for Windows PowerShell - multi-platform compilation
# Supports: linux/amd64 (Windows/Linux), linux/arm64 (Raspberry Pi)

param(
    [string]$TargetArch = "amd64",
    [string]$Version = "dev"
)

$BinaryName = "sensor-verwerker"

Write-Host "Building sensor-verwerker v${Version}" -ForegroundColor Cyan

switch ($TargetArch.ToLower()) {
    "amd64" {
        Write-Host "Building for AMD64 (Windows/Linux x86_64)..." -ForegroundColor Yellow
        $env:GOOS = "linux"
        $env:GOARCH = "amd64"
        $env:CGO_ENABLED = "0"
        go build -o "${BinaryName}-amd64.exe" -ldflags="-X main.Version=${Version}" .
        Write-Host "✓ Binary: ${BinaryName}-amd64.exe" -ForegroundColor Green
        break
    }
    "arm64" {
        Write-Host "Building for ARM64 (Raspberry Pi)..." -ForegroundColor Yellow
        $env:GOOS = "linux"
        $env:GOARCH = "arm64"
        $env:CGO_ENABLED = "0"
        go build -o "${BinaryName}-arm64" -ldflags="-X main.Version=${Version}" .
        Write-Host "✓ Binary: ${BinaryName}-arm64" -ForegroundColor Green
        break
    }
    "all" {
        Write-Host "Building for all platforms..." -ForegroundColor Yellow
        
        # AMD64
        $env:GOOS = "linux"
        $env:GOARCH = "amd64"
        $env:CGO_ENABLED = "0"
        go build -o "${BinaryName}-amd64.exe" -ldflags="-X main.Version=${Version}" .
        Write-Host "✓ Binary: ${BinaryName}-amd64.exe" -ForegroundColor Green
        
        # ARM64
        $env:GOOS = "linux"
        $env:GOARCH = "arm64"
        $env:CGO_ENABLED = "0"
        go build -o "${BinaryName}-arm64" -ldflags="-X main.Version=${Version}" .
        Write-Host "✓ Binary: ${BinaryName}-arm64" -ForegroundColor Green
        break
    }
    "docker" {
        Write-Host "Building Docker image for multi-platform (amd64 + arm64)..." -ForegroundColor Yellow
        if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
            Write-Host "Error: Docker is not installed or not in PATH" -ForegroundColor Red
            exit 1
        }
        
        # Check if buildx is available (needed for multi-platform builds)
        docker buildx version | Out-Null
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Warning: Docker buildx is not available. Using standard docker build (single platform)." -ForegroundColor Yellow
            docker build -t "sensor-verwerker:${Version}" -t sensor-verwerker:latest .
        } else {
            docker buildx build --platform linux/amd64,linux/arm64 -t "sensor-verwerker:${Version}" -t sensor-verwerker:latest --push .
        }
        Write-Host "✓ Docker image built: sensor-verwerker:${Version}" -ForegroundColor Green
        break
    }
    default {
        Write-Host "Usage: .\build.ps1 -TargetArch [amd64|arm64|all|docker] -Version [version]" -ForegroundColor Cyan
        Write-Host "  amd64  - Build for AMD64 (default)" -ForegroundColor Gray
        Write-Host "  arm64  - Build for ARM64 (Raspberry Pi)" -ForegroundColor Gray
        Write-Host "  all    - Build for all platforms" -ForegroundColor Gray
        Write-Host "  docker - Build and push Docker image for both platforms" -ForegroundColor Gray
        Write-Host ""
        Write-Host "Examples:" -ForegroundColor Gray
        Write-Host "  .\build.ps1                              # Build for AMD64" -ForegroundColor Gray
        Write-Host "  .\build.ps1 -TargetArch arm64            # Build for ARM64" -ForegroundColor Gray
        Write-Host "  .\build.ps1 -TargetArch all              # Build for both" -ForegroundColor Gray
        Write-Host "  .\build.ps1 -TargetArch docker -Version 1.0.0  # Docker build" -ForegroundColor Gray
        exit 1
    }
}

Write-Host "Build complete!" -ForegroundColor Green
