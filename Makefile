# Makefile for Kairo CLI

.PHONY: default build ls ui add add-ai done rm seed reset dev help

# Default command when you just type 'make'
default: ui

# Build binary
build:
	@go build -o kairo main.go

# Run TUI (Rebuilds automatically)
ui: build
	@./kairo ui

# List tasks (Rebuilds automatically)
ls: build
	@./kairo ls

# Add a normal task -> Example: make add TASK="Write unit tests"
add: build
	@if [ -z "$(TASK)" ]; then \
		echo "❌ Please provide a task string. Usage: make add TASK=\"My task description\""; \
	else \
		./kairo add "$(TASK)"; \
	fi

# Add an AI-assisted task -> Example: make add-ai TASK="Refactor auth service"
add-ai: build
	@if [ -z "$(TASK)" ]; then \
		echo "❌ Please provide a task string. Usage: make add-ai TASK=\"My task description\""; \
	else \
		./kairo add "$(TASK)" --ai; \
	fi

# Complete a task by number -> Example: make done N=1
done: build
	@if [ -z "$(N)" ]; then \
		echo "❌ Please provide a task number. Usage: make done N=1"; \
	else \
		./kairo done $(N); \
	fi

# Delete a task by number -> Example: make rm N=1
rm: build
	@if [ -z "$(N)" ]; then \
		echo "❌ Please provide a task number. Usage: make rm N=1"; \
	else \
		./kairo rm $(N); \
	fi

# Seed test tasks
seed: build
	@echo "🌱 Adding test tasks..."
	@./kairo add "Setup Ollama local LLM wrapper"
	@./kairo add "Implement prompt parser for AI"
	@./kairo add "Write unit tests for DB layer"
	@./kairo ls

# Full Dev Workflow: Resets database, seeds demo tasks, and launches UI in 1 step!
dev: build
	@echo "⚡ Running complete dev cycle..."
	@rm -f ~/.kairo.db
	@./kairo add "Configure Cobra CLI flags"
	@./kairo add "Refactor database models"
	@./kairo ui

# Clean binary and database
reset:
	@rm -f kairo ~/.kairo.db
	@echo "✨ Cleaned binary and reset local database."

# Run Go unit tests
test:
	@echo "🧪 Running unit tests..."
	@go test -v ./...