package yamltree

import (
	"fmt"
	"sort"
)

// BuildTree converts unmarshalled YAML data (map[string]any, []any, or scalar)
// into a slice of Node for rendering. Map keys are sorted for consistent output.
func BuildTree(data any) []Node {
	return buildNodes(data, "")
}

func buildNodes(data any, parentPath string) []Node {
	switch v := data.(type) {
	case map[string]any:
		return buildMap(v, parentPath)
	case []any:
		return buildSequence(v, parentPath)
	default:
		return nil
	}
}

func buildMap(m map[string]any, parentPath string) []Node {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	nodes := make([]Node, 0, len(keys))
	for _, k := range keys {
		path := joinPath(parentPath, k)
		nodes = append(nodes, buildNode(k, path, m[k]))
	}
	return nodes
}

func buildSequence(s []any, parentPath string) []Node {
	nodes := make([]Node, 0, len(s))
	for i, item := range s {
		key := fmt.Sprintf("%d", i)
		path := joinPath(parentPath, key)
		nodes = append(nodes, buildNode(key, path, item))
	}
	return nodes
}

func buildNode(key, path string, value any) Node {
	switch v := value.(type) {
	case map[string]any:
		return Node{
			Key:      key,
			Path:     path,
			Type:     NodeMap,
			Children: buildMap(v, path),
		}
	case []any:
		return Node{
			Key:      key,
			Path:     path,
			Type:     NodeSequence,
			Children: buildSequence(v, path),
		}
	default:
		return Node{
			Key:   key,
			Path:  path,
			Type:  NodeScalar,
			Value: formatScalar(v),
		}
	}
}

func formatScalar(v any) string {
	if v == nil {
		return "null"
	}
	return fmt.Sprintf("%v", v)
}

func joinPath(parent, key string) string {
	if parent == "" {
		return key
	}
	return parent + "." + key
}
