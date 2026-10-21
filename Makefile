# Makefile for Kairo CLI

.PHONY: build run ls seed ui clean reset-db

# Build the binary
build:
	@echo "🔨 Building Kairo binary..."
	@go build -o kairo main.go
	@echo "✅ Build complete!"

# Build and list tasks
ls: build
	@./kairo ls

# Build and launch the TUI
ui: build
	@./kairo ui

# Seed the database with sample tasks for testing
seed: build
	@echo "🌱 Seeding demo tasks..."
	@./kairo add "Setup Ollama local LLM wrapper"
	@./kairo add "Implement prompt parser for AI"
	@./kairo add "Write unit tests for DB layer"
	@./kairo ls

# Clean compiled binary
clean:
	@echo "🧹 Removing compiled binary..."
	@rm -f kairo

# Reset local database (Useful when testing clean states)
reset-db:
	@echo "⚠️ Removing local database ~/.kairo.db..."
	@rm -f ~/.kairo.db
	@echo "✨ Database reset complete."
