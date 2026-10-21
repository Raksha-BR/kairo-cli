// internal/db/tasks.go
package db

import (
	"github.com/google/uuid"
)

// Task represents a row in the tasks table
type Task struct {
	ID     string
	Title  string
	Status string
}

// internal/db/tasks.go (Update these two functions)

// InsertTask saves a new task, now including Git context
func InsertTask(title string, repoPath string, branchName string) error {
	id := uuid.New().String()
	query := `INSERT INTO tasks (id, title, status, repo_path, branch_name) VALUES (?, ?, 'TODO', ?, ?)`
	
	_, err := DB.Exec(query, id, title, repoPath, branchName)
	return err
}

// GetTasks retrieves tasks filtered by the current Git context
func GetTasks(repoPath string, branchName string) ([]Task, error) {
	query := `SELECT id, title, status FROM tasks WHERE repo_path = ? AND branch_name = ?`
	
	// If not in a git repo, fetch global tasks where repo_path is empty
	if repoPath == "" {
		query = `SELECT id, title, status FROM tasks WHERE repo_path = '' OR repo_path IS NULL`
	}

	rows, err := DB.Query(query, repoPath, branchName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
// CompleteTask marks a task as DONE by its ID
func CompleteTask(id string) error {
	query := `UPDATE tasks SET status = 'DONE' WHERE id = ?`
	_, err := DB.Exec(query, id)
	return err
}

// DeleteTask removes a task from the database by its ID
func DeleteTask(id string) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := DB.Exec(query, id)
	return err
}

// InsertTaskWithSubtasks saves a parent task and a list of AI generated subtasks
func InsertTaskWithSubtasks(parentTitle string, repoPath string, branchName string, subtasks []string) error {
	parentID := uuid.New().String()

	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert Parent Task
	parentQuery := `INSERT INTO tasks (id, title, status, repo_path, branch_name) VALUES (?, ?, 'TODO', ?, ?)`
	_, err = tx.Exec(parentQuery, parentID, parentTitle, repoPath, branchName)
	if err != nil {
		return err
	}

	// Insert Subtasks linked by parent_id
	subQuery := `INSERT INTO tasks (id, title, status, repo_path, branch_name, parent_id) VALUES (?, ?, 'TODO', ?, ?, ?)`
	for _, sub := range subtasks {
		subID := uuid.New().String()
		_, err = tx.Exec(subQuery, subID, "└─ "+sub, repoPath, branchName, parentID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}