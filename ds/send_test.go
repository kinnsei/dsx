package ds_test

import (
	"context"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/kinnsei/dsx/ds"
	"github.com/starfederation/datastar-go/datastar"
)

// testComponent is a minimal templ component for testing.
func testComponent(id, text string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<div id="`+id+`">`+text+`</div>`)
		return err
	})
}

func sse(t *testing.T) *datastar.ServerSentEventGenerator {
	t.Helper()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	sse, err := datastar.NewSSE(w, r)
	if err != nil {
		t.Fatalf("NewSSE: %v", err)
	}
	return sse
}

func TestSender_Patch(t *testing.T) {
	sse := sse(t)
	comp := testComponent("my-el", "hello")

	err := ds.Send.Patch(sse, comp)
	if err != nil {
		t.Fatalf("Patch: %v", err)
	}
}

func TestSender_Drawer(t *testing.T) {
	sse := sse(t)
	comp := testComponent("drawer-content", "Drawer content")

	err := ds.Send.Drawer(context.Background(), sse, comp)
	if err != nil {
		t.Fatalf("Drawer: %v", err)
	}
}

func TestSender_Drawer_Expandable(t *testing.T) {
	sse := sse(t)
	comp := testComponent("drawer-content", "Drawer content")

	err := ds.Send.Drawer(context.Background(), sse, comp, ds.WithDrawerExpandable())
	if err != nil {
		t.Fatalf("Drawer(expandable): %v", err)
	}
}

func TestSender_Drawer_CustomMaxWidth(t *testing.T) {
	sse := sse(t)
	comp := testComponent("drawer-content", "Drawer content")

	err := ds.Send.Drawer(context.Background(), sse, comp,
		ds.WithDrawerMaxWidth("max-w-2xl"))
	if err != nil {
		t.Fatalf("Drawer(maxWidth): %v", err)
	}
}

func TestSender_Drawer_AllOptions(t *testing.T) {
	sse := sse(t)
	comp := testComponent("drawer-content", "Drawer content")

	err := ds.Send.Drawer(context.Background(), sse, comp,
		ds.WithDrawerMaxWidth("max-w-4xl"),
		ds.WithDrawerExpandable())
	if err != nil {
		t.Fatalf("Drawer(all options): %v", err)
	}
}

func TestSender_HideDrawer(t *testing.T) {
	sse := sse(t)
	err := ds.Send.HideDrawer(sse)
	if err != nil {
		t.Fatalf("HideDrawer: %v", err)
	}
}

func TestSender_Modal(t *testing.T) {
	sse := sse(t)
	comp := testComponent("modal-content", "Modal content")

	err := ds.Send.Modal(context.Background(), sse, comp)
	if err != nil {
		t.Fatalf("Modal: %v", err)
	}
}

func TestSender_Modal_CustomMaxWidth(t *testing.T) {
	sse := sse(t)
	comp := testComponent("modal-content", "Modal content")

	err := ds.Send.Modal(context.Background(), sse, comp,
		ds.WithModalMaxWidth("max-w-3xl"))
	if err != nil {
		t.Fatalf("Modal(maxWidth): %v", err)
	}
}

func TestSender_HideModal(t *testing.T) {
	sse := sse(t)
	err := ds.Send.HideModal(sse)
	if err != nil {
		t.Fatalf("HideModal: %v", err)
	}
}

func TestSender_Toast(t *testing.T) {
	sse := sse(t)
	err := ds.Send.Toast(sse, ds.ToastSuccess, "Operation completed")
	if err != nil {
		t.Fatalf("Toast(success): %v", err)
	}
}

func TestSender_Toast_Levels(t *testing.T) {
	tests := []struct {
		level   ds.ToastLevel
		message string
	}{
		{ds.ToastInfo, "Info message"},
		{ds.ToastSuccess, "Success message"},
		{ds.ToastWarning, "Warning message"},
		{ds.ToastError, "Error message"},
	}

	for _, tt := range tests {
		t.Run(string(tt.level), func(t *testing.T) {
			sse := sse(t)
			err := ds.Send.Toast(sse, tt.level, tt.message)
			if err != nil {
				t.Fatalf("Toast(%s): %v", tt.level, err)
			}
		})
	}
}

func TestSender_Toast_Persistent(t *testing.T) {
	sse := sse(t)
	err := ds.Send.Toast(sse, ds.ToastInfo, "Persistent toast",
		ds.WithToastPersistent())
	if err != nil {
		t.Fatalf("Toast(persistent): %v", err)
	}
}

func TestSender_Toast_WithDuration(t *testing.T) {
	sse := sse(t)
	err := ds.Send.Toast(sse, ds.ToastInfo, "Short toast",
		ds.WithToastDuration(1000))
	if err != nil {
		t.Fatalf("Toast(duration): %v", err)
	}
}

func TestSender_Toast_WithAction(t *testing.T) {
	sse := sse(t)
	err := ds.Send.Toast(sse, ds.ToastWarning, "Action required",
		ds.WithToastAction("Retry", "/api/retry"))
	if err != nil {
		t.Fatalf("Toast(action): %v", err)
	}
}

func TestSender_Toast_WithLink(t *testing.T) {
	sse := sse(t)
	err := ds.Send.Toast(sse, ds.ToastInfo, "See details",
		ds.WithToastLink("here", "/details"))
	if err != nil {
		t.Fatalf("Toast(link): %v", err)
	}
}

func TestSender_Toast_AllOptions(t *testing.T) {
	sse := sse(t)
	err := ds.Send.Toast(sse, ds.ToastError, "Complex toast",
		ds.WithToastDuration(5000),
		ds.WithToastAction("Fix", "/api/fix"),
		ds.WithToastLink("more info", "/help"))
	if err != nil {
		t.Fatalf("Toast(all options): %v", err)
	}
}

func TestSender_ToastComponent(t *testing.T) {
	sse := sse(t)
	comp := testComponent("toast-content", "Custom toast")

	err := ds.Send.ToastComponent(context.Background(), sse, comp)
	if err != nil {
		t.Fatalf("ToastComponent: %v", err)
	}
}

func TestSender_Toast_PatchesToContainer(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	sse, err := datastar.NewSSE(w, r)
	if err != nil {
		t.Fatalf("NewSSE: %v", err)
	}

	_ = ds.Send.Toast(sse, ds.ToastSuccess, "Test")

	body := w.Body.String()
	if !strings.Contains(body, "#toast-container") {
		t.Errorf("toast patch should use #toast-container selector, body: %s", body[:min(len(body), 200)])
	}
}

func TestSender_Drawer_PatchesToContainer(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	sse, err := datastar.NewSSE(w, r)
	if err != nil {
		t.Fatalf("NewSSE: %v", err)
	}

	comp := testComponent("drawer-content", "Content")
	_ = ds.Send.Drawer(context.Background(), sse, comp)

	body := w.Body.String()
	if !strings.Contains(body, "drawer-panel") {
		t.Errorf("drawer patch should reference drawer-panel, body: %s", body[:min(len(body), 200)])
	}
}

func TestSender_Modal_PatchesToContainer(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	sse, err := datastar.NewSSE(w, r)
	if err != nil {
		t.Fatalf("NewSSE: %v", err)
	}

	comp := testComponent("modal-content", "Content")
	_ = ds.Send.Modal(context.Background(), sse, comp)

	body := w.Body.String()
	if !strings.Contains(body, "modal-panel") {
		t.Errorf("modal patch should reference modal-panel, body: %s", body[:min(len(body), 200)])
	}
}
