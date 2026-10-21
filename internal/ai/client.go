// internal/ai/client.go
package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SubTaskResponse structured JSON expected from LLM
type SubTaskResponse struct {
	Subtasks []string `json:"subtasks"`
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

// DecomposeTask sends a task description to Ollama and receives a list of subtasks
func DecomposeTask(taskTitle string) ([]string, error) {
	// Prompt engineering forcing structured JSON output
	prompt := fmt.Sprintf(`You are a software project manager assistant. 
Decompose the following high-level software task into 3 to 5 concise, actionable, sequential sub-tasks.
Task: "%s"

Respond strictly with valid JSON in this exact format:
{"subtasks": ["subtask 1", "subtask 2", "subtask 3"]}`, taskTitle)

	reqBody := OllamaRequest{
		Model:  "llama3.2", // Default lightweight local model, change if using another (e.g. mistral)
		Prompt: prompt,
		Stream: false,
		Format: "json",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("could not connect to Ollama (is Ollama running locally on port 11434?): %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(bodyBytes, &ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to parse Ollama response: %v", err)
	}

	var parsedResp SubTaskResponse
	if err := json.Unmarshal([]byte(ollamaResp.Response), &parsedResp); err != nil {
		return nil, fmt.Errorf("failed to parse AI JSON subtasks: %v", err)
	}

	return parsedResp.Subtasks, nil
}