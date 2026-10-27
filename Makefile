# Makefile for Kairo CLI

.PHONY: default build ls ui stats clear export dev reset help

default: ui

build:
	@go build -o kairo main.go

ui: build
	@./kairo ui

ls: build
	@./kairo ls

stats: build
	@./kairo stats

clear: build
	@./kairo clear

export: build
	@./kairo export

# ⚡ Fresh 1-Command Dev Setup: Cleans binary, seeds tasks, and launches TUI
dev: reset build
	@echo "🔄 Running fresh dev environment setup..."
	@LIPGLOSS_HAS_DARK_BACKGROUND=true ./kairo add "Submit pull request review" --due today -p HIGH > /dev/null
	@LIPGLOSS_HAS_DARK_BACKGROUND=true ./kairo add "Prepare sprint roadmap" --due tomorrow -p MEDIUM > /dev/null
	@LIPGLOSS_HAS_DARK_BACKGROUND=true ./kairo add "Benchmark SQLite query latency" --due 3d -p LOW > /dev/null
	@LIPGLOSS_HAS_DARK_BACKGROUND=true ./kairo ui

reset:
	@rm -f kairo ~/.kairo.db
	@echo "✨ Cleaned binary and reset local database."

help:
	@echo "Kairo CLI Targets:"
	@echo "  make dev    - Fresh setup & launch TUI"
	@echo "  make ui     - Open interactive TUI dashboard"
	@echo "  make stats  - Display workspace metrics & analytics"
	@echo "  make export - Export tasks as Markdown PR checklist"
	@echo "  make clear  - Purge completed tasks"
	@echo "  make reset  - Wipe database & clean binary"