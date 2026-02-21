package aichat

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/plaenen/webx/ds"
	"github.com/plaenen/webx/utils"
	"github.com/starfederation/datastar-go/datastar"
)

// quickReplyExpr builds a Datastar expression that sets the input signal
// to the reply value, then POSTs to the submit URL.
func quickReplyExpr(chatID, value, submitURL string) string {
	signals := utils.Signals(chatID, AIChatSignals{})
	return fmt.Sprintf("%s; %s",
		signals.SetString("input", value),
		ds.PostOnce(submitURL),
	)
}

// ReadSignals reads the AI chat's namespaced signals from the request.
func ReadSignals(chatID string, r *http.Request, dest *AIChatSignals) error {
	sanitizedID := strings.ReplaceAll(chatID, "-", "_")
	wrapper := map[string]any{sanitizedID: dest}
	if err := datastar.ReadSignals(r, &wrapper); err != nil {
		return fmt.Errorf("read aichat signals: %w", err)
	}
	return nil
}
