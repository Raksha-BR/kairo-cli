// internal/db/tasks.go
package db

import (
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID               string
	Title            string
	Status           string
	Priority         string
	DueDate          sql.NullTime
	TimeSpentSeconds int
}

type TaskStats struct {
	Total      int
	Done       int
	Todo       int
	HighPrio   int
	MedPrio    int
	LowPrio    int
	Overdue    int
	TotalSpent int
}

func InsertTask(title string, priority string, dueDate *time.Time, repoPath string, branchName string) error {
	id := uuid.New().String()
	if priority == "" {
		priority = "MEDIUM"
	}

	var query string
	var err error

	if dueDate != nil {
		query = `INSERT INTO tasks (id, title, status, priority, due_date, repo_path, branch_name) VALUES (?, ?, 'TODO', ?, ?, ?, ?)`
		_, err = DB.Exec(query, id, title, priority, *dueDate, repoPath, branchName)
	} else {
		query = `INSERT INTO tasks (id, title, status, priority, repo_path, branch_name) VALUES (?, ?, 'TODO', ?, ?, ?)`
		_, err = DB.Exec(query, id, title, priority, repoPath, branchName)
	}

	return err
}

func GetTasks(repoPath string, branchName string) ([]Task, error) {
	query := `SELECT id, title, status, priority, due_date, time_spent_seconds FROM tasks WHERE repo_path = ? AND branch_name = ? ORDER BY created_at DESC`

	if repoPath == "" {
		query = `SELECT id, title, status, priority, due_date, time_spent_seconds FROM tasks WHERE repo_path = '' OR repo_path IS NULL ORDER BY created_at DESC`
	}

	rows, err := DB.Query(query, repoPath, branchName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Status, &t.Priority, &t.DueDate, &t.TimeSpentSeconds); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func GetStats(repoPath string, branchName string) (TaskStats, error) {
	tasks, err := GetTasks(repoPath, branchName)
	if err != nil {
		return TaskStats{}, err
	}

	now := time.Now()
	stats := TaskStats{Total: len(tasks)}

	for _, t := range tasks {
		stats.TotalSpent += t.TimeSpentSeconds
		if t.Status == "DONE" {
			stats.Done++
		} else {
			stats.Todo++
			if t.DueDate.Valid && t.DueDate.Time.Before(now) {
				stats.Overdue++
			}
		}

		switch t.Priority {
		case "HIGH":
			stats.HighPrio++
		case "LOW":
			stats.LowPrio++
		default:
			stats.MedPrio++
		}
	}
	return stats, nil
}

func UpdateTaskTitle(id string, newTitle string) error {
	query := `UPDATE tasks SET title = ? WHERE id = ?`
	_, err := DB.Exec(query, newTitle, id)
	return err
}

func UpdateTaskPriority(id string, newPriority string) error {
	query := `UPDATE tasks SET priority = ? WHERE id = ?`
	_, err := DB.Exec(query, newPriority, id)
	return err
}

func UpdateTimeSpent(id string, additionalSeconds int) error {
	query := `UPDATE tasks SET time_spent_seconds = time_spent_seconds + ? WHERE id = ?`
	_, err := DB.Exec(query, additionalSeconds, id)
	return err
}

func InsertSubtasksForParent(subtasks []string, priority string, repoPath string, branchName string) error {
	for _, sub := range subtasks {
		if err := InsertTask("└─ "+sub, priority, nil, repoPath, branchName); err != nil {
			return err
		}
	}
	return nil
}

func CompleteTask(id string) error {
	query := `UPDATE tasks SET status = 'DONE' WHERE id = ?`
	_, err := DB.Exec(query, id)
	return err
}

func DeleteTask(id string) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := DB.Exec(query, id)
	return err
}

// Batch Operations for Multi-Selected Tasks
func BatchCompleteTasks(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := `UPDATE tasks SET status = 'DONE' WHERE id IN (` + strings.Join(placeholders, ",") + `)`
	_, err := DB.Exec(query, args...)
	return err
}

func BatchDeleteTasks(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := `DELETE FROM tasks WHERE id IN (` + strings.Join(placeholders, ",") + `)`
	_, err := DB.Exec(query, args...)
	return err
}

func BatchUpdatePriority(ids []string, priority string) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids)+1)
	args[0] = priority
	for i, id := range ids {
		placeholders[i] = "?"
		args[i+1] = id
	}
	query := `UPDATE tasks SET priority = ? WHERE id IN (` + strings.Join(placeholders, ",") + `)`
	_, err := DB.Exec(query, args...)
	return err
}

func PurgeDoneTasks(repoPath string, branchName string) (int64, error) {
	query := `DELETE FROM tasks WHERE status = 'DONE' AND repo_path = ? AND branch_name = ?`
	if repoPath == "" {
		query = `DELETE FROM tasks WHERE status = 'DONE' AND (repo_path = '' OR repo_path IS NULL)`
	}

	res, err := DB.Exec(query, repoPath, branchName)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}