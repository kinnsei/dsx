package aichat

import (
	"fmt"
	"net/http"

	"github.com/plaenen/webx/ds"
)

// quickReplyExpr builds a Datastar expression that sets the input signal
// to the reply value, then POSTs to the submit URL.
func quickReplyExpr(chatID, value, submitURL string) string {
	signals := ds.NewSignals(chatID, AIChatSignals{})
	return fmt.Sprintf("%s; %s",
		signals.SetString("input", value),
		ds.PostOnce(submitURL),
	)
}

// ReadSignals reads the AI chat's namespaced signals from the request.
func ReadSignals(chatID string, r *http.Request, dest *AIChatSignals) error {
	return ds.ReadSignals(chatID, r, dest)
}
