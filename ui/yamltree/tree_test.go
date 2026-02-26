package yamltree

import (
	"testing"
)

func TestBuildTree_FlatMap(t *testing.T) {
	data := map[string]any{
		"host": "localhost",
		"port": 8080,
	}
	nodes := BuildTree(data)
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}
	// Sorted alphabetically
	if nodes[0].Key != "host" || nodes[0].Path != "host" || nodes[0].Value != "localhost" {
		t.Errorf("node 0: got key=%q path=%q value=%q", nodes[0].Key, nodes[0].Path, nodes[0].Value)
	}
	if nodes[1].Key != "port" || nodes[1].Path != "port" || nodes[1].Value != "8080" {
		t.Errorf("node 1: got key=%q path=%q value=%q", nodes[1].Key, nodes[1].Path, nodes[1].Value)
	}
}

func TestBuildTree_NestedMap(t *testing.T) {
	data := map[string]any{
		"server": map[string]any{
			"host": "0.0.0.0",
			"port": 3000,
		},
	}
	nodes := BuildTree(data)
	if len(nodes) != 1 {
		t.Fatalf("expected 1 root node, got %d", len(nodes))
	}
	server := nodes[0]
	if server.Type != NodeMap || server.Key != "server" || server.Path != "server" {
		t.Fatalf("server node: type=%d key=%q path=%q", server.Type, server.Key, server.Path)
	}
	if len(server.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(server.Children))
	}
	if server.Children[0].Path != "server.host" {
		t.Errorf("expected path server.host, got %q", server.Children[0].Path)
	}
}

func TestBuildTree_Sequence(t *testing.T) {
	data := map[string]any{
		"ports": []any{8080, 8443, 9090},
	}
	nodes := BuildTree(data)
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}
	ports := nodes[0]
	if ports.Type != NodeSequence {
		t.Fatalf("expected NodeSequence, got %d", ports.Type)
	}
	if len(ports.Children) != 3 {
		t.Fatalf("expected 3 children, got %d", len(ports.Children))
	}
	if ports.Children[0].Path != "ports.0" || ports.Children[0].Value != "8080" {
		t.Errorf("child 0: path=%q value=%q", ports.Children[0].Path, ports.Children[0].Value)
	}
	if ports.Children[2].Path != "ports.2" || ports.Children[2].Value != "9090" {
		t.Errorf("child 2: path=%q value=%q", ports.Children[2].Path, ports.Children[2].Value)
	}
}

func TestBuildTree_NullAndBool(t *testing.T) {
	data := map[string]any{
		"debug":   true,
		"verbose": false,
		"extra":   nil,
	}
	nodes := BuildTree(data)
	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(nodes))
	}
	// Sorted: debug, extra, verbose
	if nodes[0].Value != "true" {
		t.Errorf("debug: got %q", nodes[0].Value)
	}
	if nodes[1].Value != "null" {
		t.Errorf("extra: got %q", nodes[1].Value)
	}
	if nodes[2].Value != "false" {
		t.Errorf("verbose: got %q", nodes[2].Value)
	}
}

func TestBuildTree_ScalarRoot(t *testing.T) {
	nodes := BuildTree("hello")
	if len(nodes) != 0 {
		t.Fatalf("expected 0 nodes for scalar root, got %d", len(nodes))
	}
}

func TestBuildTree_Nil(t *testing.T) {
	nodes := BuildTree(nil)
	if len(nodes) != 0 {
		t.Fatalf("expected 0 nodes for nil, got %d", len(nodes))
	}
}
