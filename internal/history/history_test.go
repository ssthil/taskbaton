package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestArchiveCreatesFile(t *testing.T) {
	dir := t.TempDir()
	content := "# Baton — feature/auth\n\nsome content"
	if err := os.WriteFile(filepath.Join(dir, "current.md"), []byte(content), 0600); err != nil {
		t.Fatalf("writing current.md: %v", err)
	}

	now := time.Date(2026, 6, 9, 0, 0, 0, 0, time.UTC)
	if err := Archive(dir, "feature/auth", now); err != nil {
		t.Fatalf("Archive returned error: %v", err)
	}

	dest := filepath.Join(dir, "history", "2026-06-09_feature-auth.md")
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("archived file not found at %s: %v", dest, err)
	}
	if string(got) != content {
		t.Errorf("archived content mismatch:\ngot:  %q\nwant: %q", string(got), content)
	}
}

func TestArchiveSlugifiesStage(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "current.md"), []byte("content"), 0600); err != nil {
		t.Fatalf("writing current.md: %v", err)
	}

	now := time.Date(2026, 6, 9, 0, 0, 0, 0, time.UTC)
	if err := Archive(dir, "feature/user auth", now); err != nil {
		t.Fatalf("Archive returned error: %v", err)
	}

	dest := filepath.Join(dir, "history", "2026-06-09_feature-user-auth.md")
	if _, err := os.Stat(dest); err != nil {
		t.Errorf("expected archived file at %s, got error: %v", dest, err)
	}
}

func TestArchiveMissingCurrentMd(t *testing.T) {
	dir := t.TempDir()
	// Do not create current.md — Archive should return nil silently.
	now := time.Date(2026, 6, 9, 0, 0, 0, 0, time.UTC)
	if err := Archive(dir, "any-stage", now); err != nil {
		t.Errorf("Archive with missing current.md should return nil, got: %v", err)
	}
}

func TestListEmpty(t *testing.T) {
	dir := t.TempDir()
	// No history/ subdirectory exists.
	entries, err := List(dir)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %v", entries)
	}
}

func TestListOrder(t *testing.T) {
	dir := t.TempDir()
	histDir := filepath.Join(dir, "history")
	if err := os.MkdirAll(histDir, 0700); err != nil {
		t.Fatalf("creating history dir: %v", err)
	}

	files := []string{
		"2026-06-11_c.md",
		"2026-06-09_a.md",
		"2026-06-10_b.md",
	}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(histDir, f), []byte("x"), 0600); err != nil {
			t.Fatalf("writing %s: %v", f, err)
		}
	}

	got, err := List(dir)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	want := []string{
		"2026-06-09_a.md",
		"2026-06-10_b.md",
		"2026-06-11_c.md",
	}
	if len(got) != len(want) {
		t.Fatalf("expected %d entries, got %d: %v", len(want), len(got), got)
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("entry[%d]: got %q, want %q", i, got[i], w)
		}
	}
}
