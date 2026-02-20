package ds

import (
	"github.com/a-h/templ"
	"github.com/starfederation/datastar-go/datastar"
)

// Patch renders a templ component and patches it into the DOM via SSE.
// The component's root element must have an id attribute for Datastar to find the target.
// Additional options (selector, mode) can be passed through.
func (s *Sender) Patch(sse *datastar.ServerSentEventGenerator, component templ.Component, opts ...datastar.PatchElementOption) error {
	return sse.PatchElementTempl(component, opts...)
}
