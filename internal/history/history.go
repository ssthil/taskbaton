package history

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var nonAlphaHyphen = regexp.MustCompile(`[^a-z0-9-]`)

func slug(stage string) string {
	s := strings.ToLower(stage)
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, " ", "-")
	s = nonAlphaHyphen.ReplaceAllString(s, "")
	return s
}

func Archive(batonDir string, stage string, now time.Time) error {
	src := filepath.Join(batonDir, "current.md")
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil
	}

	histDir := filepath.Join(batonDir, "history")
	if err := os.MkdirAll(histDir, 0700); err != nil {
		return fmt.Errorf("history archive: %w", err)
	}

	filename := now.Format("2006-01-02") + "_" + slug(stage) + ".md"
	dst := filepath.Join(histDir, filename)

	if err := copyFile(src, dst); err != nil {
		return fmt.Errorf("history archive: %w", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Close()
}

func List(batonDir string) ([]string, error) {
	histDir := filepath.Join(batonDir, "history")

	entries, err := os.ReadDir(histDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("history archive: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}

	sort.Strings(names)
	return names, nil
}
