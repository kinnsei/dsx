package ref_test

import (
	"testing"
	"time"

	"github.com/laenen-partners/domains/ref"
)

func TestCountryLookup(t *testing.T) {
	r := ref.New()

	c, ok := r.Country("BE")
	if !ok {
		t.Fatal("expected to find BE")
	}
	if c.Name != "Belgium" {
		t.Errorf("got %q, want Belgium", c.Name)
	}
	if c.PhonePrefix != "+32" {
		t.Errorf("got %q, want +32", c.PhonePrefix)
	}

	// Case insensitive
	c2, ok := r.Country("be")
	if !ok || c2.Code != "BE" {
		t.Error("case-insensitive lookup failed")
	}

	// Not found
	_, ok = r.Country("XX")
	if ok {
		t.Error("expected XX not found")
	}
}

func TestCountryLanguages(t *testing.T) {
	r := ref.New()

	c, ok := r.Country("BE")
	if !ok {
		t.Fatal("expected to find BE")
	}
	if len(c.Languages) != 3 {
		t.Fatalf("expected 3 languages for BE, got %d", len(c.Languages))
	}
	want := map[string]bool{"nl": true, "fr": true, "de": true}
	for _, l := range c.Languages {
		if !want[l] {
			t.Errorf("unexpected language %q for BE", l)
		}
	}
}

func TestCountryTimezone(t *testing.T) {
	r := ref.New()

	c, ok := r.Country("BE")
	if !ok {
		t.Fatal("expected to find BE")
	}
	if c.DefaultTimezone != "Europe/Brussels" {
		t.Errorf("got %q, want Europe/Brussels", c.DefaultTimezone)
	}
}

func TestCurrencyLookup(t *testing.T) {
	r := ref.New()

	c, ok := r.Currency("EUR")
	if !ok {
		t.Fatal("expected to find EUR")
	}
	if c.Name != "Euro" {
		t.Errorf("got %q, want Euro", c.Name)
	}
	if c.Symbol != "€" {
		t.Errorf("got %q, want €", c.Symbol)
	}
	if c.Decimals != 2 {
		t.Errorf("got %d, want 2", c.Decimals)
	}

	// Zero-decimal currency
	jpy, ok := r.Currency("JPY")
	if !ok {
		t.Fatal("expected to find JPY")
	}
	if jpy.Decimals != 0 {
		t.Errorf("JPY decimals: got %d, want 0", jpy.Decimals)
	}

	// 3-decimal currency
	bhd, ok := r.Currency("BHD")
	if !ok {
		t.Fatal("expected to find BHD")
	}
	if bhd.Decimals != 3 {
		t.Errorf("BHD decimals: got %d, want 3", bhd.Decimals)
	}
}

func TestCurrencyForCountry(t *testing.T) {
	r := ref.New()

	c, ok := r.CurrencyForCountry("BE")
	if !ok {
		t.Fatal("expected currency for BE")
	}
	if c.Code != "EUR" {
		t.Errorf("got %q, want EUR", c.Code)
	}

	c2, ok := r.CurrencyForCountry("US")
	if !ok {
		t.Fatal("expected currency for US")
	}
	if c2.Code != "USD" {
		t.Errorf("got %q, want USD", c2.Code)
	}
}

func TestPhonePrefix(t *testing.T) {
	r := ref.New()

	prefix, ok := r.PhonePrefix("BE")
	if !ok || prefix != "+32" {
		t.Errorf("got %q, want +32", prefix)
	}

	// Reverse lookup
	c, ok := r.CountryForPhonePrefix("+32")
	if !ok || c.Code != "BE" {
		t.Error("reverse phone prefix lookup failed for +32")
	}

	// Shared prefix (+1: US and CA)
	cs, ok := r.CountriesForPhonePrefix("+1")
	if !ok || len(cs) < 2 {
		t.Errorf("expected at least 2 countries for +1, got %d", len(cs))
	}
}

func TestCountriesForCurrency(t *testing.T) {
	r := ref.New()

	cs, ok := r.CountriesForCurrency("EUR")
	if !ok {
		t.Fatal("expected countries for EUR")
	}
	if len(cs) < 10 {
		t.Errorf("expected at least 10 EUR countries, got %d", len(cs))
	}

	// CHF: CH and LI
	cs2, ok := r.CountriesForCurrency("CHF")
	if !ok || len(cs2) != 2 {
		t.Errorf("expected 2 CHF countries, got %d", len(cs2))
	}
}

// ---------------------------------------------------------------------------
// Languages
// ---------------------------------------------------------------------------

func TestLanguageLookup(t *testing.T) {
	r := ref.New()

	l, ok := r.Language("fr")
	if !ok {
		t.Fatal("expected to find fr")
	}
	if l.Name != "French" {
		t.Errorf("got %q, want French", l.Name)
	}
	if l.NativeName != "Français" {
		t.Errorf("got %q, want Français", l.NativeName)
	}

	// Case insensitive
	l2, ok := r.Language("FR")
	if !ok || l2.Code != "fr" {
		t.Error("case-insensitive language lookup failed")
	}

	_, ok = r.Language("xx")
	if ok {
		t.Error("expected xx not found")
	}
}

func TestLanguagesForCountry(t *testing.T) {
	r := ref.New()

	langs, ok := r.LanguagesForCountry("BE")
	if !ok {
		t.Fatal("expected languages for BE")
	}
	if len(langs) != 3 {
		t.Errorf("expected 3 languages, got %d", len(langs))
	}

	codes := map[string]bool{}
	for _, l := range langs {
		codes[l.Code] = true
	}
	for _, want := range []string{"nl", "fr", "de"} {
		if !codes[want] {
			t.Errorf("missing language %q for BE", want)
		}
	}
}

func TestCountriesForLanguage(t *testing.T) {
	r := ref.New()

	// French is spoken in BE, FR, LU, CA, CH, MA
	cs, ok := r.CountriesForLanguage("fr")
	if !ok {
		t.Fatal("expected countries for fr")
	}
	if len(cs) < 4 {
		t.Errorf("expected at least 4 countries for fr, got %d", len(cs))
	}
}

func TestSearchLanguages(t *testing.T) {
	r := ref.New()

	results := r.SearchLanguages("fren")
	if len(results) == 0 {
		t.Fatal("expected results for 'fren'")
	}
	if results[0].Code != "fr" {
		t.Errorf("first result: got %q, want fr", results[0].Code)
	}

	// Exact code match
	results2 := r.SearchLanguages("fr")
	if len(results2) == 0 || results2[0].Code != "fr" {
		t.Error("exact code match should be first")
	}

	if len(r.SearchLanguages("")) != 0 {
		t.Error("empty query should return nil")
	}
}

// ---------------------------------------------------------------------------
// Timezones
// ---------------------------------------------------------------------------

func TestTimezoneLookup(t *testing.T) {
	r := ref.New()

	tz, ok := r.Timezone("Europe/Brussels")
	if !ok {
		t.Fatal("expected to find Europe/Brussels")
	}
	if tz.UTCOffset != "+01:00" {
		t.Errorf("got %q, want +01:00", tz.UTCOffset)
	}
	if tz.CountryCode != "BE" {
		t.Errorf("got %q, want BE", tz.CountryCode)
	}

	_, ok = r.Timezone("Mars/Olympus_Mons")
	if ok {
		t.Error("expected Mars/Olympus_Mons not found")
	}
}

func TestTimezonesForCountry(t *testing.T) {
	r := ref.New()

	// US has multiple timezones
	tzs, ok := r.TimezonesForCountry("US")
	if !ok {
		t.Fatal("expected timezones for US")
	}
	if len(tzs) < 4 {
		t.Errorf("expected at least 4 US timezones, got %d", len(tzs))
	}

	// BE has one
	tzs2, ok := r.TimezonesForCountry("BE")
	if !ok || len(tzs2) != 1 {
		t.Errorf("expected 1 BE timezone, got %d", len(tzs2))
	}
}

func TestDefaultTimezone(t *testing.T) {
	r := ref.New()

	tz, ok := r.DefaultTimezone("AE")
	if !ok {
		t.Fatal("expected default timezone for AE")
	}
	if tz.Name != "Asia/Dubai" {
		t.Errorf("got %q, want Asia/Dubai", tz.Name)
	}
}

func TestSearchTimezones(t *testing.T) {
	r := ref.New()

	results := r.SearchTimezones("brussels")
	if len(results) == 0 {
		t.Fatal("expected results for 'brussels'")
	}
	if results[0].Name != "Europe/Brussels" {
		t.Errorf("got %q, want Europe/Brussels", results[0].Name)
	}

	// Region search
	results2 := r.SearchTimezones("europe/b")
	if len(results2) == 0 {
		t.Fatal("expected results for 'europe/b'")
	}
}

// ---------------------------------------------------------------------------
// Search (countries + currencies)
// ---------------------------------------------------------------------------

func TestSearchCountries(t *testing.T) {
	r := ref.New()

	results := r.SearchCountries("belg")
	if len(results) == 0 {
		t.Fatal("expected results for 'belg'")
	}
	if results[0].Code != "BE" {
		t.Errorf("first result: got %q, want BE", results[0].Code)
	}

	// Exact code match comes first
	results2 := r.SearchCountries("be")
	if len(results2) == 0 || results2[0].Code != "BE" {
		t.Error("exact code match should be first")
	}

	// Empty query
	if len(r.SearchCountries("")) != 0 {
		t.Error("empty query should return nil")
	}
}

func TestSearchCurrencies(t *testing.T) {
	r := ref.New()

	results := r.SearchCurrencies("euro")
	if len(results) == 0 {
		t.Fatal("expected results for 'euro'")
	}
	if results[0].Code != "EUR" {
		t.Errorf("first result: got %q, want EUR", results[0].Code)
	}

	// Symbol search
	results2 := r.SearchCurrencies("€")
	if len(results2) == 0 || results2[0].Code != "EUR" {
		t.Error("symbol search should find EUR")
	}
}

// ---------------------------------------------------------------------------
// Formatting
// ---------------------------------------------------------------------------

func TestFormatMoney(t *testing.T) {
	r := ref.New()

	tests := []struct {
		amount   float64
		currency string
		country  string
		want     string
	}{
		{1234.50, "EUR", "BE", "1.234,50 €"},
		{1234.50, "USD", "US", "$1,234.50"},
		{1234.50, "GBP", "GB", "£1,234.50"},
		{1000, "JPY", "JP", "¥1,000"},
		{1234.567, "BHD", "BH", "1,234.567 BD"},
		{0.99, "EUR", "BE", "0,99 €"},
		{1000000, "EUR", "FR", "1 000 000,00 €"},
	}

	for _, tt := range tests {
		got := r.FormatMoney(tt.amount, tt.currency, tt.country)
		if got != tt.want {
			t.Errorf("FormatMoney(%v, %q, %q) = %q, want %q",
				tt.amount, tt.currency, tt.country, got, tt.want)
		}
	}
}

func TestFormatNumber(t *testing.T) {
	r := ref.New()

	tests := []struct {
		n       float64
		country string
		want    string
	}{
		{1234567.89, "BE", "1.234.567,89"},
		{1234567.89, "US", "1,234,567.89"},
		{0.5, "BE", "0,50"},
	}

	for _, tt := range tests {
		got := r.FormatNumber(tt.n, tt.country)
		if got != tt.want {
			t.Errorf("FormatNumber(%v, %q) = %q, want %q",
				tt.n, tt.country, got, tt.want)
		}
	}
}

func TestFormatPercent(t *testing.T) {
	r := ref.New()

	tests := []struct {
		n       float64
		country string
		want    string
	}{
		{0.1925, "BE", "19,25%"},
		{0.1925, "US", "19.25%"},
		{1.0, "US", "100.00%"},
	}

	for _, tt := range tests {
		got := r.FormatPercent(tt.n, tt.country)
		if got != tt.want {
			t.Errorf("FormatPercent(%v, %q) = %q, want %q",
				tt.n, tt.country, got, tt.want)
		}
	}
}

func TestFormatDate(t *testing.T) {
	r := ref.New()

	tests := []struct {
		date    string
		country string
		want    string
	}{
		{"1985-03-22", "BE", "22/03/1985"},
		{"1985-03-22", "US", "03/22/1985"},
		{"1985-03-22", "JP", "1985/03/22"},
		{"1985-03-22", "DE", "22.03.1985"},
		{"1985-03-22", "SE", "1985-03-22"},
		{"1985-03-22", "GB", "22/03/1985"},
		{"1985-03-22", "CA", "1985-03-22"},
		{"2024-01-05", "BE", "05/01/2024"},
	}

	for _, tt := range tests {
		got := r.FormatDate(tt.date, tt.country)
		if got != tt.want {
			t.Errorf("FormatDate(%q, %q) = %q, want %q",
				tt.date, tt.country, got, tt.want)
		}
	}

	// Invalid date returns input unchanged
	got := r.FormatDate("not-a-date", "BE")
	if got != "not-a-date" {
		t.Errorf("invalid date: got %q, want %q", got, "not-a-date")
	}

	// Unknown country returns input unchanged
	got = r.FormatDate("1985-03-22", "XX")
	if got != "1985-03-22" {
		t.Errorf("unknown country: got %q, want %q", got, "1985-03-22")
	}
}

func TestFormatDateTime(t *testing.T) {
	r := ref.New()

	ts := time.Date(1985, 3, 22, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		country string
		want    string
	}{
		{"BE", "22/03/1985 14:30"},
		{"US", "03/22/1985 14:30"},
		{"JP", "1985/03/22 14:30"},
		{"DE", "22.03.1985 14:30"},
	}

	for _, tt := range tests {
		got := r.FormatDateTime(ts, tt.country)
		if got != tt.want {
			t.Errorf("FormatDateTime(_, %q) = %q, want %q",
				tt.country, got, tt.want)
		}
	}

	// Unknown country falls back to ISO format
	got := r.FormatDateTime(ts, "XX")
	if got != "1985-03-22 14:30" {
		t.Errorf("unknown country: got %q, want %q", got, "1985-03-22 14:30")
	}
}

func TestFormatIBAN(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"BE68539007547034", "BE68 5390 0754 7034"},
		{"be68539007547034", "BE68 5390 0754 7034"},
		{"DE89 3704 0044 0532 0130 00", "DE89 3704 0044 0532 0130 00"},
	}

	for _, tt := range tests {
		got := ref.FormatIBAN(tt.input)
		if got != tt.want {
			t.Errorf("FormatIBAN(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatPhone(t *testing.T) {
	r := ref.New()

	tests := []struct {
		number  string
		country string
		want    string
	}{
		{"0470123456", "BE", "+32 470 12 34 56"},
		{"+32470123456", "BE", "+32 470 12 34 56"},
		{"+12025551234", "US", "+1 202 55 51 23 4"},
		{"+442071234567", "GB", "+44 207 12 34 56 7"},
	}

	for _, tt := range tests {
		got := r.FormatPhone(tt.number, tt.country)
		if got != tt.want {
			t.Errorf("FormatPhone(%q, %q) = %q, want %q",
				tt.number, tt.country, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Integrity
// ---------------------------------------------------------------------------

func TestCountriesAndCurrencies(t *testing.T) {
	r := ref.New()

	countries := r.Countries()
	if len(countries) == 0 {
		t.Error("expected countries")
	}

	currencies := r.Currencies()
	if len(currencies) == 0 {
		t.Error("expected currencies")
	}

	// Every country's currency should exist in the registry.
	for _, c := range countries {
		if _, ok := r.Currency(c.CurrencyCode); !ok {
			t.Errorf("country %s references unknown currency %s", c.Code, c.CurrencyCode)
		}
	}

	// Every country's languages should exist in the registry.
	for _, c := range countries {
		for _, lang := range c.Languages {
			if _, ok := r.Language(lang); !ok {
				t.Errorf("country %s references unknown language %s", c.Code, lang)
			}
		}
	}

	// Every country's default timezone should exist in the registry.
	for _, c := range countries {
		if _, ok := r.Timezone(c.DefaultTimezone); !ok {
			t.Errorf("country %s references unknown timezone %s", c.Code, c.DefaultTimezone)
		}
	}
}

func TestLanguagesAndTimezones(t *testing.T) {
	r := ref.New()

	languages := r.Languages()
	if len(languages) == 0 {
		t.Error("expected languages")
	}

	timezones := r.Timezones()
	if len(timezones) == 0 {
		t.Error("expected timezones")
	}

	// Every timezone's country should exist.
	for _, tz := range timezones {
		if _, ok := r.Country(tz.CountryCode); !ok {
			t.Errorf("timezone %s references unknown country %s", tz.Name, tz.CountryCode)
		}
	}
}
