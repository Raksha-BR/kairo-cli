# Kairo CLI ⏳

*Eliminate context-switching with an AI-assisted, Git-aware terminal task manager built in Go.*

Kairo is a blazingly fast, offline-first terminal application that seamlessly integrates task management directly into your developer workflow. 

## Features
- **Git-Context Aware:** Automatically detects your active Git repository and branch, instantly filtering tasks for your current workspace.
- **Ultra-Fast:** Powered by an embedded SQLite database for <50ms response times.
- **AI-Powered (Coming Soon):** Decompose complex engineering tickets into sub-tasks using local LLMs (Ollama) or OpenAI.

## Installation
Currently in active development. To build locally:
```bash
go build -o kairo main.go
./kairo --help