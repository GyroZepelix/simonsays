default: build

bin := "simonsays"
bin_path := "bin/"

# Build the binary
build:
    @echo "Building {{bin}}..."
    go build -o {{bin_path}}{{bin}} -ldflags="-s -w" cmd/simonsays/main.go

# Format code
fmt:
    @echo "Formatting code..."
    go fmt ./...

# Lint code
lint:
    @echo "Linting code..."
    golangci-lint run

# Run the binary
run *ARGS:
    @./{{bin_path}}{{bin}} {{ARGS}}

# Clean build artifacts
clean:
    @echo "Cleaning build artifacts..."
    rm -f {{bin_path}}{{bin}}
    rmdir {{bin_path}}

# Install the binary
install:
    @echo "Installing {{bin}}"
    go install

# Create a release build for multiple platforms
release:
    @echo "Creating release builds..."
    mkdir -p dist
    # Linux
    GOOS=linux GOARCH=amd64 go build -o dist/{{bin}}-linux-amd64 -ldflags="-s -w" cmd/simonsays/main.go
    # macOS
    GOOS=darwin GOARCH=amd64 go build -o dist/{{bin}}-darwin-amd64 -ldflags="-s -w" cmd/simonsays/main.go 
    GOOS=darwin GOARCH=arm64 go build -o dist/{{bin}}-darwin-arm64 -ldflags="-s -w" cmd/simonsays/main.go 
    # Windows
    GOOS=windows GOARCH=amd64 go build -o dist/{{bin}}-windows-amd64.exe -ldflags="-s -w" cmd/simonsays/main.go 
    
    # Create checksums
    cd dist && for file in *; do sha256sum "$file" >> checksums.txt; done
