// internal/ui/tui.go
package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Raksha-BR/kairo-cli/internal/db"
	"github.com/Raksha-BR/kairo-cli/internal/git"
)

// Styles using Lipgloss
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			PaddingLeft(1).
			PaddingRight(1)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#04B575")).
				Bold(true)

	normalItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	doneItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Strikethrough(true)

	statusHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#EE6FF8")).
				Italic(true)
)

type model struct {
	tasks    []db.Task
	cursor   int
	gitCtx   git.Context
	err      error
}

func InitialModel() model {
	gitCtx := git.GetCurrentContext()
	tasks, err := db.GetTasks(gitCtx.RepoPath, gitCtx.BranchName)
	return model{
		tasks:  tasks,
		cursor: 0,
		gitCtx: gitCtx,
		err:    err,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}

		case " ": // Spacebar toggles task status
			if len(m.tasks) > 0 && m.cursor < len(m.tasks) {
				t := m.tasks[m.cursor]
				if t.Status == "DONE" {
					// Toggle back to TODO (optional future improvement)
				} else {
					_ = db.CompleteTask(t.ID)
					// Refresh task list in state
					m.tasks, _ = db.GetTasks(m.gitCtx.RepoPath, m.gitCtx.BranchName)
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	s := "\n"
	
	// Title / Context Bar
	if m.gitCtx.BranchName != "" {
		s += titleStyle.Render(fmt.Sprintf("⏳ Kairo Workspace [%s]", m.gitCtx.BranchName)) + "\n\n"
	} else {
		s += titleStyle.Render("⏳ Kairo Global Workspace") + "\n\n"
	}

	if len(m.tasks) == 0 {
		s += statusHeaderStyle.Render("🎉 No active tasks in this workspace! Press 'q' to exit.") + "\n\n"
		return s
	}

	s += statusHeaderStyle.Render("Use j/k or ↑/↓ to navigate • Space to mark DONE • q to quit") + "\n\n"

	// Render task list
	for i, task := range m.tasks {
		cursor := " "
		if m.cursor == i {
			cursor = "❯"
		}

		checkbox := "[ ]"
		if task.Status == "DONE" {
			checkbox = "[x]"
		}

		line := fmt.Sprintf("%s %s %s", cursor, checkbox, task.Title)

		if task.Status == "DONE" {
			line = doneItemStyle.Render(line)
		} else if m.cursor == i {
			line = selectedItemStyle.Render(line)
		} else {
			line = normalItemStyle.Render(line)
		}

		s += line + "\n"
	}

	return s
}