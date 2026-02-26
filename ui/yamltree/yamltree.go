package yamltree

import "github.com/a-h/templ"

// NodeType classifies a YAML node.
type NodeType int

const (
	NodeScalar   NodeType = iota // leaf value (string, number, bool, null)
	NodeMap                      // mapping (key-value pairs)
	NodeSequence                 // sequence (ordered list)
)

// Node represents a single element in the YAML tree.
type Node struct {
	Key      string   // display key (or array index as string)
	Path     string   // dot-delimited path from root, e.g. "server.ports.0"
	Type     NodeType // scalar, map, or sequence
	Value    string   // string representation for scalars; empty for collections
	Children []Node   // non-nil for map/sequence nodes
}

// TreeSignals holds the reactive state for the YAML tree editor.
// Only 5 signals regardless of tree depth.
type TreeSignals struct {
	EditingPath  string `json:"editing_path"`  // path of node being edited, "" = none
	EditingValue string `json:"editing_value"` // current input buffer
	AddParent    string `json:"add_parent"`    // parent path for add operation, "" = none
	AddKey       string `json:"add_key"`       // key name for new entry
	AddValue     string `json:"add_value"`     // value for new entry
}

// Props configures the YAML tree editor component.
type Props struct {
	ID         string           // required — namespaces signals and DOM IDs
	Class      string           // additional CSS classes on root container
	Attributes templ.Attributes // passthrough HTML attributes
	Data       any              // unmarshalled YAML data (map[string]any or []any)
	ActionURL  string           // base URL for edit/add/remove handlers
	ReadOnly   bool             // when true, hides edit/add/remove controls
}
