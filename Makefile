# Makefile for Kairo CLI

.PHONY: default build ls ui add add-ai done rm seed demo dev reset clean help

# Default target when running 'make'
default: ui

# Build binary
build:
	@go build -o kairo main.go

# Launch interactive TUI (Auto-builds)
ui: build
	@./kairo ui

# View CLI Task List (Auto-builds)
ls: build
	@./kairo ls

# View Analytics & Metrics
stats: build
	@./kairo stats

# Clear Completed Tasks
clear: build
	@./kairo clear

# 🚀 Ultimate 1-Command Demo Flow: Builds, adds sample tasks with priorities/due dates, and launches TUI!
demo: build
	@echo "🌱 Seeding demo tasks with priorities & due dates..."
	@./kairo add "Submit pull request review" --due today -p HIGH
	@./kairo add "Prepare sprint roadmap" --due tomorrow -p MEDIUM
	@./kairo add "Benchmark SQLite query latency" --due 3d -p LOW
	@./kairo ui

# ⚡ Ultimate 1-Command Dev Reset: Resets DB, rebuilds, seeds demo data, and opens TUI!
dev: reset build
	@echo "🔄 Running fresh dev environment setup..."
	@./kairo add "Submit pull request review" --due today -p HIGH
	@./kairo add "Prepare sprint roadmap" --due tomorrow -p MEDIUM
	@./kairo add "Benchmark SQLite query latency" --due 3d -p LOW
	@./kairo ui

# Clean compiled binary and local database file
reset:
	@rm -f kairo ~/.kairo.db
	@echo "✨ Cleaned binary and reset local database."