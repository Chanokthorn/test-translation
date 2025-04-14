package translation

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
)

type TokenCounter interface {
	CountTokens(ctx context.Context, texts []string) (int, error)
	GetAccCount() (int, error)
	ResetAccCount() error
}

type tokenCounter struct {
	tokenizerPath string
	accCount      int
	mu            sync.Mutex
}

func NewTokenCounter(tokenizerPath string) TokenCounter {
	return &tokenCounter{
		tokenizerPath: tokenizerPath,
	}
}

type CountTokensResponse struct {
	Data []int `json:"data"`
}

func (tc *tokenCounter) CountTokens(ctx context.Context, texts []string) (int, error) {
	cmd := exec.CommandContext(ctx, "../tokenizer/venv/bin/python3", tc.tokenizerPath, "-p", "openai", "-t")
	cmd.Args = append(cmd.Args, texts...)

	fmt.Printf("Executing command: %s\n", cmd.String())
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Output", string(output))
		return 0, fmt.Errorf("failed to execute command: %w", err)
	}

	var res CountTokensResponse
	err = json.Unmarshal(output, &res)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	totalTokens := 0
	for _, tokenCount := range res.Data {
		totalTokens += tokenCount
	}

	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.accCount += totalTokens

	return totalTokens, nil
}

func (tc *tokenCounter) GetAccCount() (int, error) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	return tc.accCount, nil
}

func (tc *tokenCounter) ResetAccCount() error {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.accCount = 0
	return nil
}
