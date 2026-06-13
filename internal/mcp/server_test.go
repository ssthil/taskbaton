package mcp

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ssthil/taskbaton/internal/baton"
)

func runRPC(t *testing.T, s *Server, requestJSON string) map[string]any {
	t.Helper()
	in := strings.NewReader(requestJSON + "\n")
	var buf bytes.Buffer
	if err := s.Serve(in, &buf); err != nil {
		t.Fatalf("Serve error: %v", err)
	}
	line, _, _ := strings.Cut(buf.String(), "\n")
	var out map[string]any
	if err := json.Unmarshal([]byte(line), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v\nraw: %s", err, buf.String())
	}
	return out
}

func TestInitialize(t *testing.T) {
	s := New(t.TempDir())
	resp := runRPC(t, s, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`)

	result, ok := resp["result"].(map[string]any)
	if !ok {
		t.Fatalf("result is not an object: %#v", resp["result"])
	}
	if got := result["protocolVersion"]; got != "2024-11-05" {
		t.Errorf("protocolVersion: got %q, want %q", got, "2024-11-05")
	}
	serverInfo, ok := result["serverInfo"].(map[string]any)
	if !ok {
		t.Fatalf("serverInfo is not an object: %#v", result["serverInfo"])
	}
	if got := serverInfo["name"]; got != "taskbaton" {
		t.Errorf("serverInfo.name: got %q, want %q", got, "taskbaton")
	}
}

func TestPing(t *testing.T) {
	s := New(t.TempDir())
	resp := runRPC(t, s, `{"jsonrpc":"2.0","id":2,"method":"ping","params":{}}`)
	result, ok := resp["result"].(map[string]any)
	if !ok {
		t.Fatalf("result is not an object: %#v", resp["result"])
	}
	if len(result) != 0 {
		t.Errorf("ping result should be {}, got %v", result)
	}
}

func TestResourcesList(t *testing.T) {
	s := New(t.TempDir())
	resp := runRPC(t, s, `{"jsonrpc":"2.0","id":3,"method":"resources/list","params":{}}`)
	result, ok := resp["result"].(map[string]any)
	if !ok {
		t.Fatalf("result is not an object: %#v", resp["result"])
	}
	raw, ok := result["resources"].([]any)
	if !ok {
		t.Fatalf("resources is not an array: %#v", result["resources"])
	}
	if len(raw) != 3 {
		t.Errorf("resources length: got %d, want 3", len(raw))
	}
	var found bool
	for _, item := range raw {
		if entry, ok := item.(map[string]any); ok && entry["uri"] == "baton://current" {
			found = true
		}
	}
	if !found {
		t.Error(`resources list does not contain uri "baton://current"`)
	}
}

func TestResourcesReadNoBaton(t *testing.T) {
	s := New(t.TempDir())
	resp := runRPC(t, s,
		`{"jsonrpc":"2.0","id":4,"method":"resources/read","params":{"uri":"baton://current"}}`)
	rpcErr, ok := resp["error"].(map[string]any)
	if !ok {
		t.Fatalf("expected error, got: %#v", resp)
	}
	if got := rpcErr["code"]; got != float64(-32602) {
		t.Errorf("error.code: got %v, want -32602", got)
	}
}

func TestResourcesReadCurrent(t *testing.T) {
	dir := t.TempDir()
	const content = "# Baton — test-stage\n\nSome content here."
	if err := os.WriteFile(filepath.Join(dir, "current.md"), []byte(content), 0600); err != nil {
		t.Fatalf("setup: %v", err)
	}
	s := New(dir)
	resp := runRPC(t, s,
		`{"jsonrpc":"2.0","id":5,"method":"resources/read","params":{"uri":"baton://current"}}`)
	result, ok := resp["result"].(map[string]any)
	if !ok {
		t.Fatalf("expected result, got: %#v", resp)
	}
	contents, ok := result["contents"].([]any)
	if !ok || len(contents) == 0 {
		t.Fatalf("contents is missing or empty: %#v", result["contents"])
	}
	entry, ok := contents[0].(map[string]any)
	if !ok {
		t.Fatalf("contents[0] is not an object: %#v", contents[0])
	}
	text, _ := entry["text"].(string)
	if !strings.Contains(text, "test-stage") {
		t.Errorf("contents[0].text does not contain %q: %s", "test-stage", text)
	}
}

func TestToolsList(t *testing.T) {
	s := New(t.TempDir())
	resp := runRPC(t, s, `{"jsonrpc":"2.0","id":6,"method":"tools/list","params":{}}`)
	result, ok := resp["result"].(map[string]any)
	if !ok {
		t.Fatalf("result is not an object: %#v", resp["result"])
	}
	raw, ok := result["tools"].([]any)
	if !ok {
		t.Fatalf("tools is not an array: %#v", result["tools"])
	}
	if len(raw) != 3 {
		t.Errorf("tools length: got %d, want 3", len(raw))
	}
	required := map[string]bool{"get_status": false, "get_next_tasks": false, "get_constraints": false}
	for _, item := range raw {
		if entry, ok := item.(map[string]any); ok {
			if name, ok := entry["name"].(string); ok {
				required[name] = true
			}
		}
	}
	for name, found := range required {
		if !found {
			t.Errorf("tools list is missing %q", name)
		}
	}
}

func TestToolsCallGetStatus(t *testing.T) {
	dir := t.TempDir()
	if err := baton.Write(dir, baton.New("my-stage")); err != nil {
		t.Fatalf("setup: %v", err)
	}
	s := New(dir)
	resp := runRPC(t, s,
		`{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"get_status"}}`)
	result, ok := resp["result"].(map[string]any)
	if !ok {
		t.Fatalf("expected result, got: %#v", resp)
	}
	contents, ok := result["content"].([]any)
	if !ok || len(contents) == 0 {
		t.Fatalf("content is missing or empty: %#v", result["content"])
	}
	text, _ := contents[0].(map[string]any)["text"].(string)
	if !strings.Contains(text, "my-stage") {
		t.Errorf("get_status text does not contain %q: %s", "my-stage", text)
	}
}

func TestToolsCallGetNextTasks(t *testing.T) {
	dir := t.TempDir()
	b := baton.New("planning")
	b.NextTasks = []string{"do A", "do B"}
	if err := baton.Write(dir, b); err != nil {
		t.Fatalf("setup: %v", err)
	}
	s := New(dir)
	resp := runRPC(t, s,
		`{"jsonrpc":"2.0","id":8,"method":"tools/call","params":{"name":"get_next_tasks"}}`)
	result, ok := resp["result"].(map[string]any)
	if !ok {
		t.Fatalf("expected result, got: %#v", resp)
	}
	contents, ok := result["content"].([]any)
	if !ok || len(contents) == 0 {
		t.Fatalf("content is missing or empty: %#v", result["content"])
	}
	text, _ := contents[0].(map[string]any)["text"].(string)
	if !strings.Contains(text, "do A") {
		t.Errorf("get_next_tasks text does not contain %q: %s", "do A", text)
	}
}

func TestNotificationNoResponse(t *testing.T) {
	s := New(t.TempDir())
	in := strings.NewReader(`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n")
	var buf bytes.Buffer
	if err := s.Serve(in, &buf); err != nil {
		t.Fatalf("Serve error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for notification, got: %s", buf.String())
	}
}

func TestUnknownMethod(t *testing.T) {
	s := New(t.TempDir())
	resp := runRPC(t, s, `{"jsonrpc":"2.0","id":9,"method":"no/such/method","params":{}}`)
	rpcErr, ok := resp["error"].(map[string]any)
	if !ok {
		t.Fatalf("expected error, got: %#v", resp)
	}
	if got := rpcErr["code"]; got != float64(-32601) {
		t.Errorf("error.code: got %v, want -32601", got)
	}
}
