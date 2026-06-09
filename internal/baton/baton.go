// Package baton provides core read/write logic for the .baton state directory.
package baton

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ssthil/taskbaton/internal/template"
)

// Baton is the in-memory representation of .baton/current.md + current.json.
type Baton struct {
	Stage         string   `json:"stage"`
	Status        string   `json:"status"`     // "open" | "sealed"
	From          string   `json:"from_tool"`
	SealedAt      string   `json:"sealed_at"`  // RFC3339 or ""
	Next          string   `json:"next_tool"`
	Completed     []string `json:"completed"`
	Decisions     []string `json:"decisions"`
	NextTasks     []string `json:"next_tasks"`
	Constraints   []string `json:"constraints"`
	OpenQuestions []string `json:"open_questions"`
}

// New returns a Baton with Status="open" and all slice fields initialized to
// empty (non-nil) slices.
func New(stage string) Baton {
	return Baton{
		Stage:         stage,
		Status:        "open",
		Completed:     []string{},
		Decisions:     []string{},
		NextTasks:     []string{},
		Constraints:   []string{},
		OpenQuestions: []string{},
	}
}

// Read reads and returns the Baton stored in <batonDir>/current.json.
func Read(batonDir string) (Baton, error) {
	path := filepath.Join(batonDir, "current.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return Baton{}, fmt.Errorf("baton read: %w", err)
	}
	var b Baton
	if err := json.Unmarshal(data, &b); err != nil {
		return Baton{}, fmt.Errorf("baton read: %w", err)
	}
	return b, nil
}

// Write persists b to <batonDir>/current.json and <batonDir>/current.md,
// creating batonDir (perm 0700) if it does not exist.
func Write(batonDir string, b Baton) error {
	if err := os.MkdirAll(batonDir, 0700); err != nil {
		return fmt.Errorf("baton write: %w", err)
	}

	// Write JSON.
	jsonData, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baton write: %w", err)
	}
	if err := os.WriteFile(filepath.Join(batonDir, "current.json"), jsonData, 0600); err != nil {
		return fmt.Errorf("baton write: %w", err)
	}

	// Write Markdown.
	rd := template.RenderData{
		Stage:         b.Stage,
		Status:        b.Status,
		From:          b.From,
		SealedAt:      b.SealedAt,
		Next:          b.Next,
		Completed:     b.Completed,
		Decisions:     b.Decisions,
		NextTasks:     b.NextTasks,
		Constraints:   b.Constraints,
		OpenQuestions: b.OpenQuestions,
	}
	md, err := template.Render(rd)
	if err != nil {
		return fmt.Errorf("baton write: %w", err)
	}
	if err := os.WriteFile(filepath.Join(batonDir, "current.md"), []byte(md), 0600); err != nil {
		return fmt.Errorf("baton write: %w", err)
	}

	return nil
}
