// internal/ai/client.go
package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Raksha-BR/kairo-cli/internal/config"
)

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

func DecomposeTask(taskTitle string) ([]string, error) {
	prompt := fmt.Sprintf(`You are a software project manager assistant. 
Decompose the following high-level software task into 3 to 5 concise, actionable, sequential sub-tasks.
Task: "%s"

Respond strictly with valid JSON in this exact format:
{"subtasks": ["subtask 1", "subtask 2", "subtask 3"]}`, taskTitle)

	// Fetch dynamic settings from ~/.kairo.yaml
	model := config.AppConfig.AI.Model
	ollamaURL := config.AppConfig.AI.OllamaURL

	reqBody := OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
		Format: "json",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/api/generate", ollamaURL)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("could not connect to AI endpoint %s (is Ollama running?): %v", endpoint, err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(bodyBytes, &ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %v", err)
	}

	var parsedResp SubTaskResponse
	if err := json.Unmarshal([]byte(ollamaResp.Response), &parsedResp); err != nil {
		return nil, fmt.Errorf("failed to parse subtask JSON: %v", err)
	}

	return parsedResp.Subtasks, nil
}