
echo "Building all binaries"

echo "Building for Windows"
GOOS=windows GOARCH=amd64 go build -a -o build/poppedbit-minesweeper-windows-amd64.exe

echo "Building for MacOS"
GOOS=darwin GOARCH=amd64 go build -a -o build/poppedbit-minesweeper-macos-amd64

echo "Building for Linux"
GOOS=linux GOARCH=amd64 go build -a -o build/poppedbit-minesweeper-linux-amd64