# Makefile for Kairo CLI

.PHONY: default build ls ui stats clear demo dev reset

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

# ⚡ Ultimate 1-Command Dev Reset: Silences color queries & seeds parent/subtasks for tree view testing
dev: reset build
	@echo "🔄 Running fresh dev environment setup..."
	@LIPGLOSS_HAS_DARK_BACKGROUND=true ./kairo add "Submit pull request review" --due today -p HIGH > /dev/null
	@LIPGLOSS_HAS_DARK_BACKGROUND=true ./kairo add "Prepare sprint roadmap" --due tomorrow -p MEDIUM > /dev/null
	@LIPGLOSS_HAS_DARK_BACKGROUND=true ./kairo add "Benchmark SQLite query latency" --due 3d -p LOW > /dev/null
	@LIPGLOSS_HAS_DARK_BACKGROUND=true ./kairo ui

demo: dev

reset:
	@rm -f kairo ~/.kairo.db
	@echo "✨ Cleaned binary and reset local database."