$DIST = ".dist"
$VERSION = "0.0.1"
$LDFLAGS = "-X main.cliVersion=$VERSION"

$target = $args[0]

if (-not $target) {
    $target = "build-windows"
}

switch ($target) {
    "build" {
        Write-Host "Building current platform..."
        if (-not (Test-Path $DIST)) {
            New-Item -ItemType Directory -Path $DIST | Out-Null
        }
        go build -ldflags $LDFLAGS -o "$DIST/ai.exe" .
        Write-Host "Build complete: $DIST/ai.exe"
    }
    "build-windows" {
        Write-Host "Building Windows version..."
        if (-not (Test-Path $DIST)) {
            New-Item -ItemType Directory -Path $DIST | Out-Null
        }
        $env:GOOS = "windows"
        $env:GOARCH = "amd64"
        go build -ldflags $LDFLAGS -o "$DIST/ai.exe" .
        Write-Host "Build complete: $DIST/ai.exe"
    }
    "build-mac" {
        Write-Host "Building macOS version..."
        if (-not (Test-Path $DIST)) {
            New-Item -ItemType Directory -Path $DIST | Out-Null
        }
        $env:GOOS = "darwin"
        $env:GOARCH = "arm64"
        go build -ldflags $LDFLAGS -o "$DIST/ai" .
        Write-Host "Build complete: $DIST/ai"
    }
    "build-linux" {
        Write-Host "Building Linux version..."
        if (-not (Test-Path $DIST)) {
            New-Item -ItemType Directory -Path $DIST | Out-Null
        }
        $env:GOOS = "linux"
        $env:GOARCH = "amd64"
        go build -ldflags $LDFLAGS -o "$DIST/ai-linux-amd64" .
        Write-Host "Build complete: $DIST/ai-linux-amd64"
    }
    "clean" {
        Write-Host "Cleaning build artifacts..."
        if (Test-Path $DIST) {
            Remove-Item -Recurse -Force $DIST
        }
        Write-Host "Clean complete."
    }
    "test" {
        Write-Host "Running tests..."
        go test ./...
    }
    "version" {
        Write-Host "Version: $VERSION"
    }
    default {
        Write-Host "Usage: ./make-windows.ps1 [target]"
        Write-Host "Targets:"
        Write-Host "  build         - Build current platform"
        Write-Host "  build-windows - Build Windows version (default)"
        Write-Host "  build-mac     - Build macOS version"
        Write-Host "  build-linux   - Build Linux version"
        Write-Host "  clean         - Remove build artifacts"
        Write-Host "  test          - Run tests"
        Write-Host "  version       - Print version info"
    }
}