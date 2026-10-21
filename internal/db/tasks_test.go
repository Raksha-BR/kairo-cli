// internal/db/tasks_test.go
package db

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) {
	var err error
	DB, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	if err := createTables(); err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}
}

func TestInsertAndGetTasks(t *testing.T) {
	setupTestDB(t)
	defer DB.Close()

	repoPath := "/home/user/projects/kairo"
	branchName := "main"

	// Test Insert
	err := InsertTask("Write unit tests", repoPath, branchName)
	if err != nil {
		t.Fatalf("failed to insert task: %v", err)
	}

	// Test Retrieve
	tasks, err := GetTasks(repoPath, branchName)
	if err != nil {
		t.Fatalf("failed to retrieve tasks: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	if tasks[0].Title != "Write unit tests" {
		t.Errorf("expected title 'Write unit tests', got %q", tasks[0].Title)
	}
}

func TestCompleteAndDeleteTask(t *testing.T) {
	setupTestDB(t)
	defer DB.Close()

	repoPath := "/home/user/projects/kairo"
	branchName := "feature/auth"

	_ = InsertTask("Refactor JWT", repoPath, branchName)
	tasks, _ := GetTasks(repoPath, branchName)
	taskID := tasks[0].ID

	// Test Complete
	if err := CompleteTask(taskID); err != nil {
		t.Fatalf("failed to complete task: %v", err)
	}

	updatedTasks, _ := GetTasks(repoPath, branchName)
	if updatedTasks[0].Status != "DONE" {
		t.Errorf("expected status 'DONE', got %q", updatedTasks[0].Status)
	}

	// Test Delete
	if err := DeleteTask(taskID); err != nil {
		t.Fatalf("failed to delete task: %v", err)
	}

	finalTasks, _ := GetTasks(repoPath, branchName)
	if len(finalTasks) != 0 {
		t.Errorf("expected 0 tasks after deletion, got %d", len(finalTasks))
	}
}