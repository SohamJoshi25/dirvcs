# Build for Linux (amd64)
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o bin/linux_amd64/dirvcs

# Build for Linux (arm64)
$env:GOARCH = "arm64"
go build -o bin/linux_arm64/dirvcs

# Build for macOS (amd64)
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -o bin/mac_amd64/dirvcs

# Build for macOS (arm64)
$env:GOARCH = "arm64"
go build -o bin/mac_arm64/dirvcs

# Clean up
Remove-Item Env:GOOS
Remove-Item Env:GOARCH

# Build for Windows
go build -o bin/windows/dirvcs.exe

Write-Host "Build complete!"