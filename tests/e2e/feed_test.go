package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestFeedPage_Renders(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/feed", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Feed containers use max-w-2xl mx-auto.
	feeds := page.Locator("main .max-w-2xl.mx-auto")
	count, err := feeds.Count()
	if err != nil {
		t.Fatalf("count feeds: %v", err)
	}
	if count < 1 {
		t.Errorf("expected at least 1 feed container, got %d", count)
	}
}

func TestFeedPage_CustomMaxWidth(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/feed", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Custom width example uses max-w-xl.
	narrow := page.Locator("main .max-w-xl.mx-auto")
	count, err := narrow.Count()
	if err != nil {
		t.Fatalf("count narrow feeds: %v", err)
	}
	if count == 0 {
		t.Error("no max-w-xl feed container found")
	}

	// Custom width example uses max-w-4xl.
	wide := page.Locator("main .max-w-4xl.mx-auto")
	count, err = wide.Count()
	if err != nil {
		t.Fatalf("count wide feeds: %v", err)
	}
	if count == 0 {
		t.Error("no max-w-4xl feed container found")
	}
}

func TestFeedPage_ContainsChildren(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/feed", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Feed items are rounded-box bordered cards.
	items := page.Locator("main .max-w-2xl .rounded-box.border")
	count, err := items.Count()
	if err != nil {
		t.Fatalf("count feed items: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 feed item cards, got %d", count)
	}
}
