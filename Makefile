.PHONY: build install clean

# Build the killhost binary
build:
	go build -o killhost ./cmd/killhost

# Install killhost to ~/bin (or /usr/local/bin if ~/bin doesn't exist)
install: build
	@mkdir -p ~/bin
	@cp killhost ~/bin/killhost
	@chmod +x ~/bin/killhost
	@echo "✓ killhost installed to ~/bin/killhost"
	@echo "  Make sure ~/bin is in your PATH"

# Install to /usr/local/bin (requires sudo)
install-system: build
	sudo cp killhost /usr/local/bin/killhost
	sudo chmod +x /usr/local/bin/killhost
	@echo "✓ killhost installed to /usr/local/bin/killhost"

# Clean build artifacts
clean:
	rm -f killhost

# Run tests (if any)
test:
	go test ./...

