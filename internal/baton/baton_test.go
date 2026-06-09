package baton_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ssthil/taskbaton/internal/baton"
)

func TestRoundTrip(t *testing.T) {
	dir := t.TempDir()

	want := baton.Baton{
		Stage:         "planning",
		Status:        "open",
		From:          "claude",
		SealedAt:      "",
		Next:          "codegen",
		Completed:     []string{"task-a", "task-b"},
		Decisions:     []string{"use postgres"},
		NextTasks:     []string{"write migrations"},
		Constraints:   []string{"no breaking changes"},
		OpenQuestions: []string{"which port?"},
	}

	if err := baton.Write(dir, want); err != nil {
		t.Fatalf("Write: %v", err)
	}

	got, err := baton.Read(dir)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	if got.Stage != want.Stage {
		t.Errorf("Stage: got %q, want %q", got.Stage, want.Stage)
	}
	if got.Status != want.Status {
		t.Errorf("Status: got %q, want %q", got.Status, want.Status)
	}
	if got.From != want.From {
		t.Errorf("From: got %q, want %q", got.From, want.From)
	}
	if got.Next != want.Next {
		t.Errorf("Next: got %q, want %q", got.Next, want.Next)
	}

	checkSlice := func(name string, got, want []string) {
		t.Helper()
		if len(got) != len(want) {
			t.Errorf("%s: len got %d, want %d", name, len(got), len(want))
			return
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("%s[%d]: got %q, want %q", name, i, got[i], want[i])
			}
		}
	}

	checkSlice("Completed", got.Completed, want.Completed)
	checkSlice("Decisions", got.Decisions, want.Decisions)
	checkSlice("NextTasks", got.NextTasks, want.NextTasks)
	checkSlice("Constraints", got.Constraints, want.Constraints)
	checkSlice("OpenQuestions", got.OpenQuestions, want.OpenQuestions)
}

func TestNewSlicesNotNil(t *testing.T) {
	b := baton.New("x")

	if b.Completed == nil {
		t.Error("Completed must be non-nil")
	}
	if b.Decisions == nil {
		t.Error("Decisions must be non-nil")
	}
	if b.NextTasks == nil {
		t.Error("NextTasks must be non-nil")
	}
	if b.Constraints == nil {
		t.Error("Constraints must be non-nil")
	}
	if b.OpenQuestions == nil {
		t.Error("OpenQuestions must be non-nil")
	}

	// Slices must be empty, not nil.
	if len(b.Completed) != 0 {
		t.Errorf("Completed: expected empty, got %v", b.Completed)
	}
	if len(b.Decisions) != 0 {
		t.Errorf("Decisions: expected empty, got %v", b.Decisions)
	}
	if len(b.NextTasks) != 0 {
		t.Errorf("NextTasks: expected empty, got %v", b.NextTasks)
	}
	if len(b.Constraints) != 0 {
		t.Errorf("Constraints: expected empty, got %v", b.Constraints)
	}
	if len(b.OpenQuestions) != 0 {
		t.Errorf("OpenQuestions: expected empty, got %v", b.OpenQuestions)
	}
}

func TestWriteCreatesDualFiles(t *testing.T) {
	dir := t.TempDir()

	b := baton.New("setup")
	if err := baton.Write(dir, b); err != nil {
		t.Fatalf("Write: %v", err)
	}

	jsonPath := filepath.Join(dir, "current.json")
	mdPath := filepath.Join(dir, "current.md")

	if _, err := os.Stat(jsonPath); err != nil {
		t.Errorf("current.json not found: %v", err)
	}
	if _, err := os.Stat(mdPath); err != nil {
		t.Errorf("current.md not found: %v", err)
	}
}
