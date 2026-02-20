package ds

import (
	"github.com/starfederation/datastar-go/datastar"
)

// Redirect sends an SSE event that navigates the browser to the given URL.
// Wraps the Datastar SDK's sse.Redirect() which handles Firefox quirks automatically.
func (s *Sender) Redirect(sse *datastar.ServerSentEventGenerator, url string) error {
	return sse.Redirect(url)
}
