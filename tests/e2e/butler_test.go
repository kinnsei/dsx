package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestButlerPage_Renders(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/examples/butler", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	heading := page.Locator("h1", pw.PageLocatorOptions{
		HasText: "Customer Care Hub",
	})
	count, err := heading.Count()
	if err != nil {
		t.Fatalf("count heading: %v", err)
	}
	if count == 0 {
		t.Error("no Customer Care Hub heading found")
	}
}

func TestButlerPage_BriefingPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/examples/butler", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	briefing := page.Locator("main .bg-neutral")
	count, err := briefing.Count()
	if err != nil {
		t.Fatalf("count briefing: %v", err)
	}
	if count == 0 {
		t.Error("no briefing card found on butler page")
	}
}

func TestButlerPage_CommandBarPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/examples/butler", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Command bar has data-signals and a placeholder.
	commandBar := page.Locator("main [data-signals]").First()
	count, err := commandBar.Count()
	if err != nil {
		t.Fatalf("count command bar: %v", err)
	}
	if count == 0 {
		t.Error("no command bar found on butler page")
	}
}

func TestButlerPage_ScrollStripPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/examples/butler", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	strip := page.Locator("main .overflow-x-auto")
	count, err := strip.Count()
	if err != nil {
		t.Fatalf("count scroll strip: %v", err)
	}
	if count == 0 {
		t.Error("no scroll strip found on butler page")
	}
}

func TestButlerPage_FeedItemsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/examples/butler", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Feed items have data-signals.
	items := page.Locator("main [data-signals]")
	count, err := items.Count()
	if err != nil {
		t.Fatalf("count feed items: %v", err)
	}
	// At least command bar + several feed items.
	if count < 3 {
		t.Errorf("expected at least 3 interactive items, got %d", count)
	}
}
