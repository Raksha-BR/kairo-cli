// internal/ui/tui.go
package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Raksha-BR/kairo-cli/internal/ai"
	"github.com/Raksha-BR/kairo-cli/internal/db"
	"github.com/Raksha-BR/kairo-cli/internal/git"
)

type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type aiResultMsg struct {
	parentID string
	priority string
	subtasks []string
	err      error
}

func decomposeTaskCmd(parentID string, title string, priority string) tea.Cmd {
	return func() tea.Msg {
		subtasks, err := ai.DecomposeTask(title)
		return aiResultMsg{
			parentID: parentID,
			priority: priority,
			subtasks: subtasks,
			err:      err,
		}
	}
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			PaddingLeft(1).
			PaddingRight(1)

	progressStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

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

	tabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8BE9FD")).
			Bold(true)

	multiSelectStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF79C6")).
				Bold(true)

	subtaskTagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BD93F9")).
			Italic(true)

	timerActiveStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#50FA7B")).
				Background(lipgloss.Color("#282A36")).
				Padding(0, 1)

	aiStatusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF79C6")).
			Bold(true)

	overdueBadgeStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FF5555")).
				SetString("[OVERDUE]")

	dueTodayBadgeStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFB86C")).
				SetString("[DUE TODAY]")

	dueBadgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BD93F9"))

	helpBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			Width(68)

	highPriorityStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FF5555")).
				SetString("[HIGH]")

	medPriorityStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFB86C")).
				SetString("[MED]")

	lowPriorityStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#8BE9FD")).
				SetString("[LOW]")
)

type model struct {
	allTasks         []db.Task
	tasks            []db.Task
	selectedIDs      map[string]bool
	collapsedParents map[string]bool
	cursor           int
	gitCtx           git.Context
	err              error
	adding           bool
	editing          bool
	searching        bool
	showHelp         bool
	aiLoading        bool
	timerActive      bool
	activeTimerID    string
	timerSeconds     int
	aiMessage        string
	editTaskID       string
	filterStatus     string
	sortBy           string
	textInput        textinput.Model
	searchInput      textinput.Model
}

func InitialModel() model {
	gitCtx := git.GetCurrentContext()
	tasks, _ := db.GetTasks(gitCtx.RepoPath, gitCtx.BranchName)

	ti := textinput.New()
	ti.Placeholder = "Type task description..."
	ti.CharLimit = 150
	ti.Width = 45

	si := textinput.New()
	si.Placeholder = "Search tasks..."
	si.CharLimit = 50
	si.Width = 30

	return model{
		allTasks:         tasks,
		tasks:            tasks,
		selectedIDs:      make(map[string]bool),
		collapsedParents: make(map[string]bool),
		cursor:           0,
		gitCtx:           gitCtx,
		textInput:        ti,
		searchInput:      si,
		adding:           false,
		editing:          false,
		searching:        false,
		showHelp:         false,
		aiLoading:        false,
		timerActive:      false,
		filterStatus:     "ALL",
		sortBy:           "DEFAULT",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) getSelectedIDsList() []string {
	var list []string
	for id, isSel := range m.selectedIDs {
		if isSel {
			list = append(list, id)
		}
	}
	return list
}

func (m model) getSubtaskCount(parentID string) int {
	count := 0
	for _, t := range m.allTasks {
		if t.ParentID.Valid && t.ParentID.String == parentID {
			count++
		}
	}
	return count
}

func (m model) filterTasks() []db.Task {
	query := strings.ToLower(m.searchInput.Value())
	var res []db.Task

	for _, t := range m.allTasks {
		if t.ParentID.Valid && m.collapsedParents[t.ParentID.String] {
			continue
		}

		if m.filterStatus == "TODO" && t.Status == "DONE" {
			continue
		}
		if m.filterStatus == "DONE" && t.Status != "DONE" {
			continue
		}

		if query != "" && !strings.Contains(strings.ToLower(t.Title), query) {
			continue
		}

		res = append(res, t)
	}

	switch m.sortBy {
	case "PRIORITY":
		prioOrder := map[string]int{"HIGH": 1, "MEDIUM": 2, "LOW": 3}
		sort.SliceStable(res, func(i, j int) bool {
			return prioOrder[res[i].Priority] < prioOrder[res[j].Priority]
		})
	case "DUE_DATE":
		sort.SliceStable(res, func(i, j int) bool {
			if !res[i].DueDate.Valid {
				return false
			}
			if !res[j].DueDate.Valid {
				return true
			}
			return res[i].DueDate.Time.Before(res[j].DueDate.Time)
		})
	}

	return res
}

func renderProgressBar(doneCount, totalCount int) string {
	if totalCount == 0 {
		return ""
	}
	percent := (doneCount * 100) / totalCount
	width := 10
	filled := (doneCount * width) / totalCount
	empty := width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	return fmt.Sprintf("[%s] %d%% (%d/%d done)", bar, percent, doneCount, totalCount)
}

func formatDuration(sec int) string {
	m := sec / 60
	s := sec % 60
	return fmt.Sprintf("%02dm %02ds", m, s)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tickMsg:
		if m.timerActive {
			m.timerSeconds++
			return m, doTick()
		}

	case aiResultMsg:
		m.aiLoading = false
		if msg.err != nil {
			m.aiMessage = fmt.Sprintf("⚠️ AI failed: %v", msg.err)
		} else {
			_ = db.InsertSubtasksForParent(msg.parentID, msg.subtasks, msg.priority, m.gitCtx.RepoPath, m.gitCtx.BranchName)
			m.allTasks, _ = db.GetTasks(m.gitCtx.RepoPath, m.gitCtx.BranchName)
			m.tasks = m.filterTasks()
			m.aiMessage = fmt.Sprintf("✨ Added %d AI subtasks!", len(msg.subtasks))
		}
		return m, nil
	}

	if m.showHelp {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "?", "esc", "q":
				m.showHelp = false
				return m, nil
			}
		}
		return m, nil
	}

	if m.adding {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				val := m.textInput.Value()
				if val != "" {
					_, _ = db.InsertTask(val, "HIGH", nil, m.gitCtx.RepoPath, m.gitCtx.BranchName)
					m.allTasks, _ = db.GetTasks(m.gitCtx.RepoPath, m.gitCtx.BranchName)
					m.tasks = m.filterTasks()
					m.textInput.Reset()
				}
				m.adding = false
				return m, nil

			case "esc":
				m.adding = false
				m.textInput.Reset()
				return m, nil
			}
		}

		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	if m.editing {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				val := m.textInput.Value()
				if val != "" && m.editTaskID != "" {
					_ = db.UpdateTaskTitle(m.editTaskID, val)
					m.allTasks, _ = db.GetTasks(m.gitCtx.RepoPath, m.gitCtx.BranchName)
					m.tasks = m.filterTasks()
					m.textInput.Reset()
				}
				m.editing = false
				m.editTaskID = ""
				return m, nil

			case "esc":
				m.editing = false
				m.editTaskID = ""
				m.textInput.Reset()
				return m, nil
			}
		}

		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	if m.searching {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter", "esc":
				m.searching = false
				return m, nil
			}
		}

		m.searchInput, cmd = m.searchInput.Update(msg)
		m.tasks = m.filterTasks()
		m.cursor = 0
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.timerActive && m.activeTimerID != "" {
				_ = db.UpdateTimeSpent(m.activeTimerID, m.timerSeconds)
			}
			return m, tea.Quit

		case "o":
			if len(m.tasks) > 0 && m.cursor < len(m.tasks) {
				tID := m.tasks[m.cursor].ID
				subCount := m.getSubtaskCount(tID)
				if subCount > 0 {
					m.collapsedParents[tID] = !m.collapsedParents[tID]
					m.tasks = m.filterTasks()
				}
			}

		case "x":
			if len(m.tasks) > 0 && m.cursor < len(m.tasks) {
				tID := m.tasks[m.cursor].ID
				m.selectedIDs[tID] = !m.selectedIDs[tID]
			}

		case "f":
			if len(m.tasks) > 0 && m.cursor < len(m.tasks) {
				currentTask := m.tasks[m.cursor]
				if m.timerActive {
					_ = db.UpdateTimeSpent(m.activeTimerID, m.timerSeconds)
					m.allTasks, _ = db.GetTasks(m.gitCtx.RepoPath, m.gitCtx.BranchName)
					m.tasks = m.filterTasks()
					m.timerActive = false
					m.aiMessage = fmt.Sprintf("⏱️ Saved %s work time!", formatDuration(m.timerSeconds))
					m.timerSeconds = 0
					m.activeTimerID = ""
				} else {
					m.timerActive = true
					m.activeTimerID = currentTask.ID
					m.timerSeconds = 0
					m.aiMessage = fmt.Sprintf("⏱️ Focus Timer Started on %q", currentTask.Title)
					return m, doTick()
				}
			}

		case "?":
			m.showHelp = true
			return m, nil

		case "s":
			switch m.sortBy {
			case "DEFAULT":
				m.sortBy = "PRIORITY"
			case "PRIORITY":
				m.sortBy = "DUE_DATE"
			default:
				m.sortBy = "DEFAULT"
			}
			m.tasks = m.filterTasks()
			m.cursor = 0

		case "tab":
			switch m.filterStatus {
			case "ALL":
				m.filterStatus = "TODO"
			case "TODO":
				m.filterStatus = "DONE"
			default:
				m.filterStatus = "ALL"
			}
			m.tasks = m.filterTasks()
			m.cursor = 0

		case "c":
			purgedCount, _ := db.PurgeDoneTasks(m.gitCtx.RepoPath, m.gitCtx.BranchName)
			m.allTasks, _ = db.GetTasks(m.gitCtx.RepoPath, m.gitCtx.BranchName)
			m.tasks = m.filterTasks()
			m.cursor = 0
			m.selectedIDs = make(map[string]bool)
			if purgedCount > 0 {
				m.aiMessage = fmt.Sprintf("🧹 Purged %d completed tasks!", purgedCount)
			} else {
				m.aiMessage = "ℹ️ No completed tasks to clear."
			}

		case "p":
			selList := m.getSelectedIDsList()
			if len(selList) > 0 {
				_ = db.BatchUpdatePriority(selList, "HIGH")
				m.aiMessage = fmt.Sprintf("⚡ Updated priority to HIGH for %d tasks!", len(selList))
				m.selectedIDs = make(map[string]bool)
			} else if len(m.tasks) > 0 && m.cursor < len(m.tasks) {
				t := m.tasks[m.cursor]
				nextPrio := "MEDIUM"
				switch t.Priority {
				case "LOW":
					nextPrio = "MEDIUM"
				case "MEDIUM":
					nextPrio = "HIGH"
				case "HIGH":
					nextPrio = "LOW"
				}
				_ = db.UpdateTaskPriority(t.ID, nextPrio)
			}
			m.allTasks, _ = db.GetTasks(m.gitCtx.RepoPath, m.gitCtx.BranchName)
			m.tasks = m.filterTasks()

		case "A": // Trigger Async Non-blocking AI Decomposition
			if len(m.tasks) > 0 && m.cursor < len(m.tasks) && !m.aiLoading {
				selectedTask := m.tasks[m.cursor]
				m.aiLoading = true
				m.aiMessage = fmt.Sprintf("🤖 AI decomposing %q...", selectedTask.Title)
				return m, decomposeTaskCmd(selectedTask.ID, selectedTask.Title, selectedTask.Priority)
			}

		case "/":
			m.searching = true
			m.searchInput.Focus()
			return m, textinput.Blink

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}

		case "a":
			m.adding = true
			m.textInput.SetValue("")
			m.textInput.Focus()
			return m, textinput.Blink

		case "e":
			if len(m.tasks) > 0 && m.cursor < len(m.tasks) {
				m.editing = true
				m.editTaskID = m.tasks[m.cursor].ID
				m.textInput.SetValue(m.tasks[m.cursor].Title)
				m.textInput.Focus()
				return m, textinput.Blink
			}

		case "d":
			selList := m.getSelectedIDsList()
			if len(selList) > 0 {
				_ = db.BatchDeleteTasks(selList)
				m.aiMessage = fmt.Sprintf("🗑️ Bulk deleted %d tasks!", len(selList))
				m.selectedIDs = make(map[string]bool)
			} else if len(m.tasks) > 0 && m.cursor < len(m.tasks) {
				_ = db.DeleteTask(m.tasks[m.cursor].ID)
			}
			m.allTasks, _ = db.GetTasks(m.gitCtx.RepoPath, m.gitCtx.BranchName)
			m.tasks = m.filterTasks()
			if m.cursor >= len(m.tasks) && m.cursor > 0 {
				m.cursor--
			}

		case " ":
			selList := m.getSelectedIDsList()
			if len(selList) > 0 {
				_ = db.BatchCompleteTasks(selList)
				m.aiMessage = fmt.Sprintf("✅ Bulk completed %d tasks!", len(selList))
				m.selectedIDs = make(map[string]bool)
			} else if len(m.tasks) > 0 && m.cursor < len(m.tasks) {
				t := m.tasks[m.cursor]
				if t.Status != "DONE" {
					_ = db.CompleteTask(t.ID)
				}
			}
			m.allTasks, _ = db.GetTasks(m.gitCtx.RepoPath, m.gitCtx.BranchName)
			m.tasks = m.filterTasks()
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "\n"

	doneCount := 0
	for _, t := range m.allTasks {
		if t.Status == "DONE" {
			doneCount++
		}
	}
	pBar := renderProgressBar(doneCount, len(m.allTasks))

	if m.gitCtx.BranchName != "" {
		s += titleStyle.Render(fmt.Sprintf("⏳ Kairo Workspace [%s]", m.gitCtx.BranchName)) + " " + progressStyle.Render(pBar)
	} else {
		s += titleStyle.Render("⏳ Kairo Global Workspace") + " " + progressStyle.Render(pBar)
	}

	if m.timerActive {
		s += " " + timerActiveStyle.Render(fmt.Sprintf("⏱️ RUNNING: %s", formatDuration(m.timerSeconds)))
	}
	s += "\n\n"

	if m.showHelp {
		modalContent := "⌨️  Kairo Keyboard Reference\n" +
			"----------------------------------------\n" +
			" j / k or ↑/↓ : Navigate items\n" +
			" Space        : Toggle task DONE (or Bulk Complete)\n" +
			" o            : Toggle Expand/Collapse AI subtasks tree\n" +
			" x            : Toggle Multi-Select item ([•])\n" +
			" f            : Toggle Focus Timer (Stopwatch)\n" +
			" p            : Cycle Priority (or Bulk HIGH)\n" +
			" s            : Cycle Sort (Default/Priority/Due Date)\n" +
			" a            : Add new task inline\n" +
			" e            : Edit highlighted task title\n" +
			" d            : Delete highlighted task (or Bulk Delete)\n" +
			" c            : Clear/Purge all completed tasks\n" +
			" Shift+A      : AI Decompose highlighted task\n" +
			" Tab          : Switch filter tab (ALL/TODO/DONE)\n" +
			" /            : Search/Filter tasks in real-time\n" +
			" ?            : Toggle this Help Menu\n" +
			" q            : Quit dashboard\n\n" +
			"(Press '?' or Esc to close)"
		return s + helpBoxStyle.Render(modalContent) + "\n"
	}

	s += tabStyle.Render(fmt.Sprintf(" View Tab: [%s] | Sort: [%s]", m.filterStatus, m.sortBy))

	selCount := len(m.getSelectedIDsList())
	if selCount > 0 {
		s += "  " + multiSelectStyle.Render(fmt.Sprintf("[• %d SELECTED]", selCount))
	}

	if m.aiMessage != "" {
		s += "  " + aiStatusStyle.Render(m.aiMessage)
	}
	s += "\n\n"

	if m.adding {
		s += "➕ Add New Task:\n\n"
		s += m.textInput.View() + "\n\n"
		s += "(Press Enter to save, Esc to cancel)\n"
		return s
	}

	if m.editing {
		s += "✏️ Edit Selected Task:\n\n"
		s += m.textInput.View() + "\n\n"
		s += "(Press Enter to save changes, Esc to cancel)\n"
		return s
	}

	if m.searching {
		s += "🔍 Search: " + m.searchInput.View() + " (Press Enter or Esc when done)\n\n"
	} else if m.searchInput.Value() != "" {
		s += fmt.Sprintf("🔍 Search Query: %q\n\n", m.searchInput.Value())
	}

	if len(m.tasks) == 0 {
		s += statusHeaderStyle.Render("🎉 No matching tasks! Press 'a' to add, 'Tab' to switch view, '?' for help, 'q' to quit.") + "\n\n"
		return s
	}

	s += statusHeaderStyle.Render("j/k: Nav • Space: Done • o: Expand • x: Select • f: Timer • p: Prio • s: Sort • Shift+A: AI • ?: Help • q: Quit") + "\n\n"

	now := time.Now()

	for i, task := range m.tasks {
		cursor := " "
		if m.cursor == i {
			cursor = "❯"
		}

		selectTag := " "
		if m.selectedIDs[task.ID] {
			selectTag = "•"
		}

		indent := ""
		if task.ParentID.Valid {
			indent = "   └─ "
		}

		timeLabel := ""
		if task.TimeSpentSeconds > 0 {
			timeLabel = fmt.Sprintf(" ⏱️ %s", formatDuration(task.TimeSpentSeconds))
		}

		var line string

		if task.Status == "DONE" {
			plainText := fmt.Sprintf("%s%s %s[x] [%s] %s%s", indent, cursor, selectTag, task.Priority, task.Title, timeLabel)
			line = doneItemStyle.Render(plainText)
		} else {
			checkbox := "[ ]"

			var prioBadge string
			switch task.Priority {
			case "HIGH":
				prioBadge = highPriorityStyle.Render()
			case "LOW":
				prioBadge = lowPriorityStyle.Render()
			default:
				prioBadge = medPriorityStyle.Render()
			}

			dueLabel := ""
			if task.DueDate.Valid {
				dTime := task.DueDate.Time
				if dTime.Before(now) {
					dueLabel = overdueBadgeStyle.Render() + " "
				} else if dTime.Year() == now.Year() && dTime.YearDay() == now.YearDay() {
					dueLabel = dueTodayBadgeStyle.Render() + " "
				} else {
					dueLabel = dueBadgeStyle.Render(fmt.Sprintf("[%s]", dTime.Format("Jan 02"))) + " "
				}
			}

			subCount := m.getSubtaskCount(task.ID)
			treeTag := ""
			if subCount > 0 {
				if m.collapsedParents[task.ID] {
					treeTag = subtaskTagStyle.Render(fmt.Sprintf(" ▶ [%d subtasks]", subCount)) + " "
				} else {
					treeTag = subtaskTagStyle.Render(fmt.Sprintf(" ▼ [%d subtasks]", subCount)) + " "
				}
			}

			line = fmt.Sprintf("%s%s %s%s %s %s%s%s%s", indent, cursor, selectTag, checkbox, prioBadge, dueLabel, treeTag, task.Title, timeLabel)

			if m.cursor == i {
				line = selectedItemStyle.Render(line)
			} else {
				line = normalItemStyle.Render(line)
			}
		}

		s += line + "\n"
	}

	return s
}