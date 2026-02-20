package ds

// Send provides backend SSE operations (drawer, toast, etc.).
// Frontend attribute helpers remain as top-level ds.XXX functions.
var Send = &Sender{}

// Sender groups backend SSE operations that send events to the browser.
type Sender struct{}
