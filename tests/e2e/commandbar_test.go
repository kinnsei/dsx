package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestCommandBarPage_RendersCollapsed(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/command-bar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Command bars have data-signals attribute (Datastar).
	bars := page.Locator("main [data-signals]")
	count, err := bars.Count()
	if err != nil {
		t.Fatalf("count command bars: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 command bars, got %d", count)
	}
}

func TestCommandBarPage_ClickOpensTextMode(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/command-bar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	bar := page.Locator("#demo-text-only")

	// Click the collapsed placeholder to open.
	placeholder := bar.Locator("text=Type a message...")
	if err := placeholder.Click(); err != nil {
		t.Fatalf("click placeholder: %v", err)
	}

	// Textarea should become visible.
	textarea := bar.Locator("textarea")
	if err := textarea.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("textarea did not become visible: %v", err)
	}
}

func TestCommandBarPage_CloseButton(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/command-bar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	bar := page.Locator("#demo-text-only")

	// Open.
	placeholder := bar.Locator("text=Type a message...")
	if err := placeholder.Click(); err != nil {
		t.Fatalf("click placeholder: %v", err)
	}

	textarea := bar.Locator("textarea")
	if err := textarea.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("textarea did not become visible: %v", err)
	}

	// Click close button (X icon).
	closeBtn := bar.Locator("button.btn-square").First()
	if err := closeBtn.Click(); err != nil {
		t.Fatalf("click close: %v", err)
	}

	// Textarea should be hidden again.
	if err := textarea.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateHidden,
	}); err != nil {
		t.Fatalf("textarea did not hide after close: %v", err)
	}
}

func TestCommandBarPage_AllModesHasTabs(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/command-bar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	bar := page.Locator("#demo-all-modes")

	// Open the command bar.
	placeholder := bar.Locator("text=Type, upload, or record...")
	if err := placeholder.Click(); err != nil {
		t.Fatalf("click placeholder: %v", err)
	}

	// Should show Type, File, and Voice tabs.
	typeTab := bar.Locator("button", pw.LocatorLocatorOptions{HasText: "Type"})
	if err := typeTab.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("Type tab not visible: %v", err)
	}

	fileTab := bar.Locator("button", pw.LocatorLocatorOptions{HasText: "File"})
	if err := fileTab.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("File tab not visible: %v", err)
	}

	voiceTab := bar.Locator("button", pw.LocatorLocatorOptions{HasText: "Voice"})
	if err := voiceTab.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("Voice tab not visible: %v", err)
	}
}

func TestCommandBarPage_SwitchToFileMode(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/command-bar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	bar := page.Locator("#demo-all-modes")

	// Open.
	placeholder := bar.Locator("text=Type, upload, or record...")
	if err := placeholder.Click(); err != nil {
		t.Fatalf("click placeholder: %v", err)
	}

	// Click File tab.
	fileTab := bar.Locator("button", pw.LocatorLocatorOptions{HasText: "File"})
	if err := fileTab.Click(); err != nil {
		t.Fatalf("click File tab: %v", err)
	}

	// File drop zone should appear (has file input).
	fileInput := bar.Locator("input[type='file']")
	count, err := fileInput.Count()
	if err != nil {
		t.Fatalf("count file inputs: %v", err)
	}
	if count == 0 {
		t.Error("no file input found in file mode")
	}
}

func TestCommandBarPage_SwitchToVoiceMode(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/command-bar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	bar := page.Locator("#demo-all-modes")

	// Open.
	placeholder := bar.Locator("text=Type, upload, or record...")
	if err := placeholder.Click(); err != nil {
		t.Fatalf("click placeholder: %v", err)
	}

	// Click Voice tab.
	voiceTab := bar.Locator("button", pw.LocatorLocatorOptions{HasText: "Voice"})
	if err := voiceTab.Click(); err != nil {
		t.Fatalf("click Voice tab: %v", err)
	}

	// "Tap to start recording" text should appear.
	recordText := bar.Locator("text=Tap to start recording")
	if err := recordText.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("voice idle text not visible: %v", err)
	}
}

func TestCommandBarPage_Suggestions(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/command-bar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	bar := page.Locator("#demo-suggestions")

	// Open.
	placeholder := bar.Locator("text=How can I help?")
	if err := placeholder.Click(); err != nil {
		t.Fatalf("click placeholder: %v", err)
	}

	// Suggestion chips should be visible.
	suggestions := []string{"Check balance", "Transfer funds", "Report issue"}
	for _, s := range suggestions {
		chip := bar.Locator("button", pw.LocatorLocatorOptions{HasText: s})
		if err := chip.WaitFor(pw.LocatorWaitForOptions{
			State: pw.WaitForSelectorStateVisible,
		}); err != nil {
			t.Errorf("suggestion chip %q not visible: %v", s, err)
		}
	}
}
