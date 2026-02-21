package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestAIChatPage_Renders(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/ai-chat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Page title should be present.
	h1 := page.Locator("h1", pw.PageLocatorOptions{HasText: "AI Chat"})
	if err := h1.WaitFor(); err != nil {
		t.Fatalf("h1 not found: %v", err)
	}
}

func TestAIChatPage_CollapsedBarVisible(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/ai-chat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The collapsed bar (Live Demo) should show the placeholder text.
	bar := page.Locator("#demo-aichat")
	if err := bar.WaitFor(); err != nil {
		t.Fatalf("collapsed bar not found: %v", err)
	}
}

func TestAIChatPage_ExpandsOnClick(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/ai-chat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Click the collapsed bar to expand.
	bar := page.Locator("#demo-aichat")
	if err := bar.Click(); err != nil {
		t.Fatalf("click bar: %v", err)
	}

	// The messages container should become visible.
	messages := page.Locator("#demo-aichat-messages")
	if err := messages.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("messages container not visible after expand: %v", err)
	}
}

func TestAIChatPage_UserMessageRendered(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/ai-chat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The standalone User Message example should be visible.
	msg := page.Locator("text=max needs football boots size 38")
	count, err := msg.Count()
	if err != nil {
		t.Fatalf("count user message: %v", err)
	}
	if count == 0 {
		t.Fatal("standalone user message not found")
	}
}

func TestAIChatPage_AssistantMessageRendered(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/ai-chat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The standalone Assistant Message should be visible.
	msg := page.Locator("text=Love it. When are you thinking?")
	count, err := msg.Count()
	if err != nil {
		t.Fatalf("count assistant message: %v", err)
	}
	if count == 0 {
		t.Fatal("standalone assistant message not found")
	}
}

func TestAIChatPage_TypingIndicatorRendered(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/ai-chat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The typing indicator example should contain the loading dots.
	dots := page.Locator(".loading-dots")
	count, err := dots.Count()
	if err != nil {
		t.Fatalf("count typing dots: %v", err)
	}
	if count == 0 {
		t.Fatal("typing indicator dots not found")
	}
}

func TestAIChatPage_QuickRepliesRendered(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/ai-chat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Quick reply buttons should be present.
	for _, label := range []string{"Netflix", "Spotify", "YouTube"} {
		btn := page.Locator("button", pw.PageLocatorOptions{HasText: label})
		count, err := btn.Count()
		if err != nil {
			t.Fatalf("count %s: %v", label, err)
		}
		if count == 0 {
			t.Errorf("quick reply button %q not found", label)
		}
	}
}

func TestAIChatPage_ResponseCardRendered(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/ai-chat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Response card title should be present.
	title := page.Locator("text=Buy football boots for Max (size 38)")
	count, err := title.Count()
	if err != nil {
		t.Fatalf("count card title: %v", err)
	}
	if count == 0 {
		t.Fatal("response card title not found")
	}

	// Tags should be present.
	shopping := page.Locator(".badge", pw.PageLocatorOptions{HasText: "Shopping"})
	count, err = shopping.Count()
	if err != nil {
		t.Fatalf("count shopping badge: %v", err)
	}
	if count == 0 {
		t.Fatal("Shopping tag not found")
	}
}

func TestAIChatPage_AssignRowRendered(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/ai-chat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Assign row should show member initials.
	me := page.Locator("text=ME").First()
	if err := me.WaitFor(); err != nil {
		t.Fatalf("assign row member ME not found: %v", err)
	}
}

func TestAIChatPage_ActionBarRendered(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/ai-chat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Action bar primary button.
	btn := page.Locator("button", pw.PageLocatorOptions{HasText: "Add to inbox"})
	count, err := btn.Count()
	if err != nil {
		t.Fatalf("count action button: %v", err)
	}
	if count == 0 {
		t.Fatal("action bar primary button not found")
	}
}
