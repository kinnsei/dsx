package ds

import (
	"fmt"

	"github.com/starfederation/datastar-go/datastar"
)

// Download triggers a file download in the browser via SSE without navigating away.
// It creates a temporary anchor element, triggers a click, and removes it.
func (s *Sender) Download(sse *datastar.ServerSentEventGenerator, url string, filename string) error {
	script := fmt.Sprintf(
		`const a=document.createElement('a');a.href='%s';a.download='%s';document.body.appendChild(a);a.click();a.remove()`,
		url, filename,
	)
	return sse.ExecuteScript(script)
}
