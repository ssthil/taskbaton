package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ssthil/taskbaton/internal/baton"
	"github.com/ssthil/taskbaton/internal/history"
)

type Server struct {
	batonDir string
}

func New(batonDir string) *Server {
	return &Server{batonDir: batonDir}
}

type request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type response struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      any       `json:"id"`
	Result  any       `json:"result,omitempty"`
	Error   *rpcError `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func send(w io.Writer, resp response) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = w.Write(data)
	return err
}

func (s *Server) Serve(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var req request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			_ = send(w, response{
				JSONRPC: "2.0",
				ID:      nil,
				Error:   &rpcError{Code: -32700, Message: "parse error"},
			})
			continue
		}

		// Notifications have no id field or id is nil — no response.
		if req.ID == nil {
			s.dispatch(w, req)
			continue
		}

		result, rpcErr := s.dispatch(w, req)
		if rpcErr != nil {
			_ = send(w, response{JSONRPC: "2.0", ID: req.ID, Error: rpcErr})
		} else {
			_ = send(w, response{JSONRPC: "2.0", ID: req.ID, Result: result})
		}
	}
	return scanner.Err()
}

// dispatch handles a request and returns (result, error). For notifications it
// returns (nil, nil) without writing anything.
func (s *Server) dispatch(w io.Writer, req request) (any, *rpcError) {
	switch req.Method {
	case "initialize":
		return map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]any{
				"resources": map[string]any{"subscribe": false, "listChanged": false},
				"tools":     map[string]any{},
			},
			"serverInfo": map[string]any{"name": "taskbaton", "version": "0.1.0"},
		}, nil

	case "notifications/initialized":
		return nil, nil

	case "ping":
		return map[string]any{}, nil

	case "resources/list":
		return map[string]any{
			"resources": []map[string]any{
				{
					"uri":         "baton://current",
					"name":        "Current Baton",
					"description": "Active baton stage — decisions, constraints, and next tasks",
					"mimeType":    "text/markdown",
				},
				{
					"uri":         "baton://current.json",
					"name":        "Current Baton (JSON)",
					"description": "Machine-readable baton state",
					"mimeType":    "application/json",
				},
				{
					"uri":         "baton://history",
					"name":        "Baton History",
					"description": "List of archived baton stages",
					"mimeType":    "text/plain",
				},
			},
		}, nil

	case "resources/read":
		return s.handleResourceRead(req.Params)

	case "tools/list":
		return map[string]any{
			"tools": []map[string]any{
				{
					"name":        "get_status",
					"description": "Returns the current baton stage name and seal state",
					"inputSchema": map[string]any{"type": "object", "properties": map[string]any{}},
				},
				{
					"name":        "get_next_tasks",
					"description": "Returns the Next Tasks list from the current baton",
					"inputSchema": map[string]any{"type": "object", "properties": map[string]any{}},
				},
				{
					"name":        "get_constraints",
					"description": "Returns the Constraints list from the current baton",
					"inputSchema": map[string]any{"type": "object", "properties": map[string]any{}},
				},
			},
		}, nil

	case "tools/call":
		return s.handleToolCall(req.Params)

	default:
		return nil, &rpcError{Code: -32601, Message: "method not found: " + req.Method}
	}
}

func (s *Server) handleResourceRead(params json.RawMessage) (any, *rpcError) {
	var p struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, &rpcError{Code: -32602, Message: "invalid params"}
	}

	switch p.URI {
	case "baton://current":
		data, err := os.ReadFile(filepath.Join(s.batonDir, "current.md"))
		if err != nil {
			return nil, &rpcError{Code: -32602, Message: "no baton found — run: taskbaton init"}
		}
		return map[string]any{
			"contents": []map[string]any{
				{"uri": "baton://current", "mimeType": "text/markdown", "text": string(data)},
			},
		}, nil

	case "baton://current.json":
		data, err := os.ReadFile(filepath.Join(s.batonDir, "current.json"))
		if err != nil {
			return nil, &rpcError{Code: -32602, Message: "no baton found — run: taskbaton init"}
		}
		return map[string]any{
			"contents": []map[string]any{
				{"uri": "baton://current.json", "mimeType": "application/json", "text": string(data)},
			},
		}, nil

	case "baton://history":
		entries, err := history.List(s.batonDir)
		if err != nil {
			return nil, &rpcError{Code: -32602, Message: err.Error()}
		}
		var text string
		if len(entries) == 0 {
			text = "(no history yet)"
		} else {
			stripped := make([]string, len(entries))
			for i, e := range entries {
				stripped[i] = strings.TrimSuffix(e, ".md")
			}
			text = strings.Join(stripped, "\n")
		}
		return map[string]any{
			"contents": []map[string]any{
				{"uri": "baton://history", "mimeType": "text/plain", "text": text},
			},
		}, nil

	default:
		return nil, &rpcError{Code: -32602, Message: "unknown resource: " + p.URI}
	}
}

func (s *Server) handleToolCall(params json.RawMessage) (any, *rpcError) {
	var p struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, &rpcError{Code: -32602, Message: "invalid params"}
	}

	toolResult := func(text string, isError bool) any {
		return map[string]any{
			"content": []map[string]any{{"type": "text", "text": text}},
			"isError": isError,
		}
	}

	switch p.Name {
	case "get_status", "get_next_tasks", "get_constraints":
		b, err := baton.Read(s.batonDir)
		if err != nil {
			return toolResult(err.Error(), true), nil
		}

		switch p.Name {
		case "get_status":
			var sb strings.Builder
			fmt.Fprintf(&sb, "Stage:   %s\n", b.Stage)
			fmt.Fprintf(&sb, "Status:  %s", b.Status)
			if b.From != "" {
				fmt.Fprintf(&sb, "\nFrom:    %s", b.From)
			}
			if b.Next != "" {
				fmt.Fprintf(&sb, "\nNext:    %s", b.Next)
			}
			return toolResult(sb.String(), false), nil

		case "get_next_tasks":
			if len(b.NextTasks) == 0 {
				return toolResult("(no next tasks recorded)", false), nil
			}
			return toolResult("Next Tasks:\n- "+strings.Join(b.NextTasks, "\n- "), false), nil

		case "get_constraints":
			if len(b.Constraints) == 0 {
				return toolResult("(no constraints recorded)", false), nil
			}
			return toolResult("Constraints — Do Not Change:\n- "+strings.Join(b.Constraints, "\n- "), false), nil
		}
	}

	return toolResult("unknown tool: "+p.Name, true), nil
}
