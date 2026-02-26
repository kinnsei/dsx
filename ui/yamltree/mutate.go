package yamltree

import (
	"fmt"
	"strconv"
	"strings"
)

// SetAtPath sets a value at a dot-delimited path in a nested map/slice structure.
func SetAtPath(data any, path string, value any) error {
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return fmt.Errorf("empty path")
	}

	parent, lastKey, err := walk(data, parts)
	if err != nil {
		return fmt.Errorf("set %q: %w", path, err)
	}

	switch p := parent.(type) {
	case map[string]any:
		p[lastKey] = value
		return nil
	case []any:
		idx, err := strconv.Atoi(lastKey)
		if err != nil || idx < 0 || idx >= len(p) {
			return fmt.Errorf("set %q: invalid array index %q", path, lastKey)
		}
		p[idx] = value
		return nil
	default:
		return fmt.Errorf("set %q: parent is not a map or slice", path)
	}
}

// DeleteAtPath removes a key at a dot-delimited path.
// For maps, the key is deleted. For slices, the element is removed.
func DeleteAtPath(data any, path string) (any, error) {
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return data, fmt.Errorf("empty path")
	}

	// Special case: deleting a root key
	if len(parts) == 1 {
		switch d := data.(type) {
		case map[string]any:
			if _, ok := d[parts[0]]; !ok {
				return data, fmt.Errorf("delete %q: key not found", path)
			}
			delete(d, parts[0])
			return data, nil
		case []any:
			idx, err := strconv.Atoi(parts[0])
			if err != nil || idx < 0 || idx >= len(d) {
				return data, fmt.Errorf("delete %q: invalid array index", path)
			}
			return append(d[:idx], d[idx+1:]...), nil
		default:
			return data, fmt.Errorf("delete %q: data is not a map or slice", path)
		}
	}

	parent, lastKey, err := walk(data, parts)
	if err != nil {
		return data, fmt.Errorf("delete %q: %w", path, err)
	}

	switch p := parent.(type) {
	case map[string]any:
		if _, ok := p[lastKey]; !ok {
			return data, fmt.Errorf("delete %q: key not found", path)
		}
		delete(p, lastKey)
		return data, nil
	case []any:
		idx, err := strconv.Atoi(lastKey)
		if err != nil || idx < 0 || idx >= len(p) {
			return data, fmt.Errorf("delete %q: invalid array index %q", path, lastKey)
		}
		// We need to update the parent's reference to this slice.
		// Since we can't do that from here, we need the grandparent.
		// For simplicity, rebuild at the grandparent level.
		newSlice := append(p[:idx], p[idx+1:]...)
		grandParentParts := parts[:len(parts)-1]
		return data, SetAtPath(data, strings.Join(grandParentParts, "."), newSlice)
	default:
		return data, fmt.Errorf("delete %q: parent is not a map or slice", path)
	}
}

// AddAtPath adds a new key/value at the given parent path.
// The parent must be a map.
func AddAtPath(data any, parentPath, key string, value any) error {
	var parent any
	if parentPath == "" {
		parent = data
	} else {
		parts := strings.Split(parentPath, ".")
		var err error
		parent, err = resolve(data, parts)
		if err != nil {
			return fmt.Errorf("add at %q: %w", parentPath, err)
		}
	}

	m, ok := parent.(map[string]any)
	if !ok {
		return fmt.Errorf("add at %q: parent is not a map", parentPath)
	}
	if _, exists := m[key]; exists {
		return fmt.Errorf("add at %q: key %q already exists", parentPath, key)
	}
	m[key] = value
	return nil
}

// walk traverses all path segments except the last, returning the parent container
// and the final key.
func walk(data any, parts []string) (parent any, lastKey string, err error) {
	if len(parts) == 1 {
		return data, parts[0], nil
	}
	parent, err = resolve(data, parts[:len(parts)-1])
	if err != nil {
		return nil, "", err
	}
	return parent, parts[len(parts)-1], nil
}

// resolve traverses the full path and returns the value at that path.
func resolve(data any, parts []string) (any, error) {
	current := data
	for _, key := range parts {
		switch c := current.(type) {
		case map[string]any:
			val, ok := c[key]
			if !ok {
				return nil, fmt.Errorf("key %q not found", key)
			}
			current = val
		case []any:
			idx, err := strconv.Atoi(key)
			if err != nil || idx < 0 || idx >= len(c) {
				return nil, fmt.Errorf("invalid array index %q", key)
			}
			current = c[idx]
		default:
			return nil, fmt.Errorf("cannot traverse into %T at key %q", current, key)
		}
	}
	return current, nil
}
