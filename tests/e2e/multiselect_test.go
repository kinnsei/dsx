package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestMultiSelectPage_Renders(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/multi-select", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	heading := page.Locator("h1", pw.PageLocatorOptions{HasText: "Multi-Select"})
	if err := heading.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("heading not visible: %v", err)
	}
}

func TestMultiSelectPage_SelectAndRemove(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/multi-select", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Click the search input to open dropdown.
	searchInput := page.Locator("#languages-search")
	if err := searchInput.Click(); err != nil {
		t.Fatalf("click search input: %v", err)
	}

	// Wait for results to load.
	frenchOption := page.Locator("#languages-results div[role='option']", pw.PageLocatorOptions{HasText: "French"})
	if err := frenchOption.WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateVisible,
		Timeout: pw.Float(5000),
	}); err != nil {
		t.Fatalf("French option not visible: %v", err)
	}

	// Click French.
	if err := frenchOption.Click(); err != nil {
		t.Fatalf("click French: %v", err)
	}

	// French tag should appear.
	frenchTag := page.Locator("#languages-tags .badge", pw.PageLocatorOptions{HasText: "French"})
	if err := frenchTag.WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateVisible,
		Timeout: pw.Float(5000),
	}); err != nil {
		t.Fatalf("French tag did not appear: %v", err)
	}

	// French should no longer be in the dropdown.
	if err := frenchOption.WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateHidden,
		Timeout: pw.Float(3000),
	}); err != nil {
		t.Fatalf("French option should be hidden after selection: %v", err)
	}

	// Remove French by clicking its ✕ button.
	removeBtn := frenchTag.Locator("button")
	if err := removeBtn.Click(); err != nil {
		t.Fatalf("click remove French: %v", err)
	}

	// French tag should disappear.
	if err := frenchTag.WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateHidden,
		Timeout: pw.Float(5000),
	}); err != nil {
		t.Fatalf("French tag did not disappear after removal: %v", err)
	}
}
