package yamltree

import (
	"testing"
)

func TestSetAtPath_Simple(t *testing.T) {
	data := map[string]any{"host": "localhost", "port": 3000}
	if err := SetAtPath(data, "port", 8080); err != nil {
		t.Fatal(err)
	}
	if data["port"] != 8080 {
		t.Errorf("expected 8080, got %v", data["port"])
	}
}

func TestSetAtPath_Nested(t *testing.T) {
	data := map[string]any{
		"server": map[string]any{
			"host": "localhost",
			"port": 3000,
		},
	}
	if err := SetAtPath(data, "server.port", 9090); err != nil {
		t.Fatal(err)
	}
	server := data["server"].(map[string]any)
	if server["port"] != 9090 {
		t.Errorf("expected 9090, got %v", server["port"])
	}
}

func TestSetAtPath_ArrayIndex(t *testing.T) {
	data := map[string]any{
		"ports": []any{8080, 8443},
	}
	if err := SetAtPath(data, "ports.0", 9090); err != nil {
		t.Fatal(err)
	}
	ports := data["ports"].([]any)
	if ports[0] != 9090 {
		t.Errorf("expected 9090, got %v", ports[0])
	}
}

func TestSetAtPath_InvalidPath(t *testing.T) {
	data := map[string]any{"host": "localhost"}
	if err := SetAtPath(data, "missing.nested", "value"); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestDeleteAtPath_RootKey(t *testing.T) {
	data := map[string]any{"host": "localhost", "port": 3000}
	result, err := DeleteAtPath(data, "port")
	if err != nil {
		t.Fatal(err)
	}
	m := result.(map[string]any)
	if _, ok := m["port"]; ok {
		t.Error("expected port to be deleted")
	}
	if m["host"] != "localhost" {
		t.Error("expected host to remain")
	}
}

func TestDeleteAtPath_NestedKey(t *testing.T) {
	data := map[string]any{
		"server": map[string]any{
			"host": "localhost",
			"port": 3000,
		},
	}
	_, err := DeleteAtPath(data, "server.port")
	if err != nil {
		t.Fatal(err)
	}
	server := data["server"].(map[string]any)
	if _, ok := server["port"]; ok {
		t.Error("expected port to be deleted")
	}
}

func TestDeleteAtPath_MissingKey(t *testing.T) {
	data := map[string]any{"host": "localhost"}
	_, err := DeleteAtPath(data, "missing")
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestAddAtPath_Root(t *testing.T) {
	data := map[string]any{"host": "localhost"}
	if err := AddAtPath(data, "", "port", 3000); err != nil {
		t.Fatal(err)
	}
	if data["port"] != 3000 {
		t.Errorf("expected 3000, got %v", data["port"])
	}
}

func TestAddAtPath_Nested(t *testing.T) {
	data := map[string]any{
		"server": map[string]any{
			"host": "localhost",
		},
	}
	if err := AddAtPath(data, "server", "port", 3000); err != nil {
		t.Fatal(err)
	}
	server := data["server"].(map[string]any)
	if server["port"] != 3000 {
		t.Errorf("expected 3000, got %v", server["port"])
	}
}

func TestAddAtPath_DuplicateKey(t *testing.T) {
	data := map[string]any{"host": "localhost"}
	if err := AddAtPath(data, "", "host", "other"); err == nil {
		t.Error("expected error for duplicate key")
	}
}

func TestAddAtPath_NotAMap(t *testing.T) {
	data := map[string]any{"name": "test"}
	if err := AddAtPath(data, "name", "key", "value"); err == nil {
		t.Error("expected error when parent is not a map")
	}
}
