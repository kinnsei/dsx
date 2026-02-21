package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestScrollStripPage_RendersCards(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/scroll-strip", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Scroll strip cards have min-w-36 class.
	cards := page.Locator("main .min-w-36")
	count, err := cards.Count()
	if err != nil {
		t.Fatalf("count cards: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 scroll strip cards, got %d", count)
	}
}

func TestScrollStripPage_TitlePresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/scroll-strip", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	title := page.Locator("main .uppercase", pw.PageLocatorOptions{
		HasText: "Pinned cases",
	})
	count, err := title.Count()
	if err != nil {
		t.Fatalf("count title: %v", err)
	}
	if count == 0 {
		t.Error("no scroll strip title found")
	}
}

func TestScrollStripPage_ActiveCard(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/scroll-strip", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Active cards use bg-neutral and border-neutral.
	active := page.Locator("main .bg-neutral.border-neutral.min-w-36")
	count, err := active.Count()
	if err != nil {
		t.Fatalf("count active cards: %v", err)
	}
	if count == 0 {
		t.Error("no active scroll strip card found")
	}
}

func TestScrollStripPage_HorizontalScroll(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/scroll-strip", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The scroll container has overflow-x-auto.
	scrollContainer := page.Locator("main .overflow-x-auto")
	count, err := scrollContainer.Count()
	if err != nil {
		t.Fatalf("count scroll containers: %v", err)
	}
	if count == 0 {
		t.Error("no horizontal scroll container found")
	}
}
