package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestBriefingPage_Renders(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/briefing", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Briefing cards use bg-neutral by default.
	cards := page.Locator("main .bg-neutral")
	count, err := cards.Count()
	if err != nil {
		t.Fatalf("count briefing cards: %v", err)
	}
	if count < 2 {
		t.Errorf("expected at least 2 briefing cards, got %d", count)
	}
}

func TestBriefingPage_TitlePresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/briefing", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	title := page.Locator("main .bg-neutral .uppercase", pw.PageLocatorOptions{
		HasText: "Morning briefing",
	})
	count, err := title.Count()
	if err != nil {
		t.Fatalf("count title: %v", err)
	}
	if count == 0 {
		t.Error("no briefing title found")
	}
}

func TestBriefingPage_CustomColors(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/briefing", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	variants := []string{"bg-error", "bg-success", "bg-info"}
	for _, v := range variants {
		loc := page.Locator("main ." + v)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", v, err)
		}
		if count == 0 {
			t.Errorf("no %s briefing card found", v)
		}
	}
}
