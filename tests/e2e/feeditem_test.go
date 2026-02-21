package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestFeedItemPage_RendersItems(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/feed-item", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Feed items are rendered with data-signals attribute (Datastar).
	items := page.Locator("main [data-signals]")
	count, err := items.Count()
	if err != nil {
		t.Fatalf("count feed items: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 feed items with data-signals, got %d", count)
	}
}

func TestFeedItemPage_ClickExpandsContent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/feed-item", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The basic feed item has ID "demo-basic".
	item := page.Locator("#demo-basic")

	// Click the header text to expand.
	header := item.Locator("text=Click to expand")
	if err := header.Click(); err != nil {
		t.Fatalf("click header: %v", err)
	}

	// Content text should become visible.
	contentText := item.Locator("text=This is the expanded content area")
	if err := contentText.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("content did not become visible: %v", err)
	}
}

func TestFeedItemPage_ClickCollapses(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/feed-item", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	item := page.Locator("#demo-basic")
	header := item.Locator("text=Click to expand")

	// Open.
	if err := header.Click(); err != nil {
		t.Fatalf("click to open: %v", err)
	}
	contentText := item.Locator("text=This is the expanded content area")
	if err := contentText.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("content did not open: %v", err)
	}

	// Close.
	if err := header.Click(); err != nil {
		t.Fatalf("click to close: %v", err)
	}
	if err := contentText.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateHidden,
	}); err != nil {
		t.Fatalf("content did not hide after second click: %v", err)
	}
}

func TestFeedItemPage_UrgencyDotsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/feed-item", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	dots := []string{"bg-error", "bg-warning"}
	for _, d := range dots {
		loc := page.Locator("main .rounded-full." + d)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s dots: %v", d, err)
		}
		if count == 0 {
			t.Errorf("no %s urgency dot found", d)
		}
	}
}

func TestFeedItemPage_ThreadPillsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/feed-item", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	pills := page.Locator("main .bg-primary\\/10")
	count, err := pills.Count()
	if err != nil {
		t.Fatalf("count thread pills: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 thread pills, got %d", count)
	}
}
