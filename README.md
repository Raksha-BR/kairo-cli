# ⏳ Kairo CLI — Context-Aware Terminal Task & AI Decomposition Engine

> **Kairo** (derived from *Kairos*, meaning the opportune moment) is a high-performance, offline-first terminal task manager built with **Go**, **Cobra**, **Charm (Bubbletea & Lipgloss)**, **SQLite**, and local **Ollama AI**.

---

## 🌟 Key Features

- 🌲 **Git Context Awareness**: Tasks are automatically bound to your active Git repository path and branch name. Switching branches dynamically filters your workspace!
- 🎨 **Interactive TUI Dashboard (`kairo ui` / `make ui`)**: Full-screen terminal dashboard with Vim navigation (`j`/`k`), priority badges, real-time search, sorting, and filter tabs.
- 🌳 **Collapsible AI Sub-task Tree View (`o` key)**: Expand or collapse nested AI subtasks (`▼ [3 subtasks]` / `▶ [3 subtasks]`) tied via SQLite relational `parent_id` mapping.
- ⏱️ **Focus Stopwatch Timer (`f` key)**: Track focus time per task directly inside the TUI with live status bar ticking and persistent database logging.
- 🎯 **Batch Multi-Selection & Bulk Operations (`x` key)**: Select multiple tasks (`[•]`) to perform bulk completion (`Space`), bulk priority upgrades (`p`), or bulk deletion (`d`).
- 🤖 **Async Local AI Task Decomposition (`--ai` / `Shift+A`)**: Non-blocking background goroutines parse complex tasks into subtasks using local Ollama (`llama3.2`) models without freezing the UI.
- 📅 **Relative Due Dates (`--due`)**: Parse relative inputs (`today`, `tomorrow`, `3d`, `1w`) with real-time `[OVERDUE]` and `[DUE TODAY]` badge alerts.
- 📋 **PR Checklist Exporter (`kairo export`)**: Export active workspace tasks directly into formatted Markdown checklists for GitHub PR descriptions.
- 📊 **Workspace Metrics & Analytics (`kairo stats`)**: Visual progress meters, completion percentages, and priority distribution.
- ⚙️ **Config & Shell Auto-Completion**: Config viewer (`kairo config`) and shell tab completion (`kairo completion bash/zsh`).

---

## 🏗 System Architecture

```mermaid
flowchart TD
    User([Developer / Terminal]) --> CLI[Cobra CLI Router]

    subgraph Core Engine
        CLI --> Git[Git Context Engine\nAuto-detects Repo & Branch]
        CLI --> Config[Viper Config Manager\n~/.kairo.yaml]
    end

    subgraph User Experience & Data
        CLI --> TUI[Interactive Bubbletea TUI\nLipgloss & Live Search]
        CLI --> DB[(SQLite Database\n~/.kairo.db)]
    end

    subgraph AI Intelligence
        TUI --> AI[Ollama Local LLM Wrapper\nAsync Non-Blocking Subtask Breakdown]
        CLI --> AI
    end