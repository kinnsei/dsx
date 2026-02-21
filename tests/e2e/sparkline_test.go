package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestSparklinePage_RendersCharts(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/sparkline", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	svgs := page.Locator("main svg polyline")
	count, err := svgs.Count()
	if err != nil {
		t.Fatalf("count polylines: %v", err)
	}
	if count == 0 {
		t.Error("no sparkline SVG polylines found")
	}
}

func TestSparklinePage_EndDotPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/sparkline", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	dots := page.Locator("main svg circle")
	count, err := dots.Count()
	if err != nil {
		t.Fatalf("count dots: %v", err)
	}
	if count == 0 {
		t.Error("no sparkline end dots (circles) found")
	}
}

func TestSparklinePage_CustomSizes(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/sparkline", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Check that SVGs with different widths exist.
	svg40 := page.Locator(`main svg[width="40"]`)
	count, err := svg40.Count()
	if err != nil {
		t.Fatalf("count 40px svg: %v", err)
	}
	if count == 0 {
		t.Error("no 40px wide sparkline found")
	}

	svg200 := page.Locator(`main svg[width="200"]`)
	count, err = svg200.Count()
	if err != nil {
		t.Fatalf("count 200px svg: %v", err)
	}
	if count == 0 {
		t.Error("no 200px wide sparkline found")
	}
}
