# build-windows.ps1 — 只执行 Windows Release 构建
if ($args.Count -ne 0) {
    Write-Error "Usage: .\build-windows.ps1"
    exit 2
}

$DIST = ".dist"
$VERSION = "0.0.1"
$LDFLAGS = "-s -w -X main.cliVersion=$VERSION"

if (-not (Test-Path $DIST)) {
    New-Item -ItemType Directory -Path $DIST | Out-Null
}
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -trimpath -ldflags $LDFLAGS -o "$DIST/backstage.exe" .
Write-Host "Build complete: $DIST/backstage.exe"
