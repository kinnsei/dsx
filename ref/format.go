package ref

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// FormatMoney formats an amount with the currency symbol and locale-specific
// separators for the given country. Falls back to plain number + currency code
// if country or currency is not found.
func (r *Registry) FormatMoney(amount float64, currencyCode, countryCode string) string {
	cur, curOK := r.Currency(currencyCode)
	cty, ctyOK := r.Country(countryCode)

	decimals := 2
	symbol := currencyCode
	if curOK {
		decimals = cur.Decimals
		symbol = cur.Symbol
	}

	decSep := "."
	thousSep := ","
	symbolAfter := false
	if ctyOK {
		decSep = cty.DecimalSep
		thousSep = cty.ThousandsSep
		symbolAfter = cty.SymbolAfter
	}

	num := formatNumber(amount, decimals, decSep, thousSep)

	if symbolAfter {
		return num + " " + symbol
	}
	return symbol + num
}

// FormatNumber formats a number with locale-specific separators for the
// given country. Uses 2 decimal places by default.
func (r *Registry) FormatNumber(n float64, countryCode string) string {
	cty, ok := r.Country(countryCode)
	if !ok {
		return formatNumber(n, 2, ".", ",")
	}
	return formatNumber(n, 2, cty.DecimalSep, cty.ThousandsSep)
}

// FormatPercent formats a ratio (0.1925) as a percentage string with
// locale-specific decimal separator. Uses 2 decimal places.
func (r *Registry) FormatPercent(n float64, countryCode string) string {
	pct := n * 100
	cty, ok := r.Country(countryCode)
	if !ok {
		return formatNumber(pct, 2, ".", "") + "%"
	}
	return formatNumber(pct, 2, cty.DecimalSep, "") + "%"
}

// FormatDate formats an ISO 8601 date string (YYYY-MM-DD) using the
// country's locale conventions. Returns the input unchanged if parsing fails
// or the country is not found.
func (r *Registry) FormatDate(isoDate, countryCode string) string {
	t, err := time.Parse("2006-01-02", isoDate)
	if err != nil {
		return isoDate
	}
	cty, ok := r.Country(countryCode)
	if !ok {
		return isoDate
	}
	return formatDateParts(t.Day(), int(t.Month()), t.Year(), cty.DateOrder, cty.DateSep)
}

// FormatDateTime formats a time.Time using the country's date conventions
// with 24-hour time appended (e.g., "22/03/1985 14:30").
func (r *Registry) FormatDateTime(t time.Time, countryCode string) string {
	cty, ok := r.Country(countryCode)
	if !ok {
		return t.Format("2006-01-02 15:04")
	}
	date := formatDateParts(t.Day(), int(t.Month()), t.Year(), cty.DateOrder, cty.DateSep)
	return date + " " + fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
}

func formatDateParts(day, month, year int, order DateOrder, sep string) string {
	d := fmt.Sprintf("%02d", day)
	m := fmt.Sprintf("%02d", month)
	y := fmt.Sprintf("%04d", year)
	switch order {
	case DateOrderMDY:
		return m + sep + d + sep + y
	case DateOrderYMD:
		return y + sep + m + sep + d
	default: // DMY
		return d + sep + m + sep + y
	}
}

// FormatIBAN formats an IBAN into groups of 4 characters.
// Input is normalized to uppercase with spaces removed.
func FormatIBAN(iban string) string {
	s := strings.ToUpper(strings.ReplaceAll(iban, " ", ""))
	var b strings.Builder
	for i, ch := range s {
		if i > 0 && i%4 == 0 {
			b.WriteByte(' ')
		}
		b.WriteRune(ch)
	}
	return b.String()
}

// FormatPhone formats a phone number with the country's international prefix.
// If the number starts with "0", it is replaced with the country prefix.
// Digits are grouped with spaces for readability.
func (r *Registry) FormatPhone(number, countryCode string) string {
	digits := stripNonDigitsKeepPlus(number)
	cty, ok := r.Country(countryCode)

	prefixLen := 0
	if ok {
		// Replace leading 0 with country prefix.
		if strings.HasPrefix(digits, "0") {
			digits = cty.PhonePrefix + digits[1:]
		}
		// If no prefix yet, add country prefix.
		if !strings.HasPrefix(digits, "+") {
			digits = cty.PhonePrefix + digits
		}
		prefixLen = len(cty.PhonePrefix)
	}

	return groupPhoneDigits(digits, prefixLen)
}

// formatNumber renders a float with the given decimal places, decimal
// separator, and thousands separator.
func formatNumber(n float64, decimals int, decSep, thousSep string) string {
	neg := n < 0
	n = math.Abs(n)

	// Round to the requested decimal places.
	factor := math.Pow(10, float64(decimals))
	rounded := math.Round(n*factor) / factor

	intPart := int64(rounded)
	fracPart := rounded - float64(intPart)

	// Integer part with thousands separator.
	intStr := strconv.FormatInt(intPart, 10)
	if thousSep != "" {
		intStr = insertThousandsSep(intStr, thousSep)
	}

	var result string
	if decimals > 0 {
		fracStr := strconv.FormatFloat(fracPart, 'f', decimals, 64)
		// fracStr is "0.XX" — take everything after "0."
		result = intStr + decSep + fracStr[2:]
	} else {
		result = intStr
	}

	if neg {
		return "-" + result
	}
	return result
}

// insertThousandsSep inserts sep every 3 digits from the right.
func insertThousandsSep(s, sep string) string {
	if len(s) <= 3 {
		return s
	}
	var b strings.Builder
	start := len(s) % 3
	if start == 0 {
		start = 3
	}
	b.WriteString(s[:start])
	for i := start; i < len(s); i += 3 {
		b.WriteString(sep)
		b.WriteString(s[i : i+3])
	}
	return b.String()
}

// stripNonDigitsKeepPlus removes everything except digits and leading '+'.
func stripNonDigitsKeepPlus(s string) string {
	var b strings.Builder
	for i, ch := range s {
		if ch == '+' && i == 0 {
			b.WriteRune(ch)
		} else if ch >= '0' && ch <= '9' {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

// groupPhoneDigits groups a phone number string into readable segments.
// prefixLen is the length of the country prefix (e.g., 3 for "+32").
// "+32470123456" with prefixLen=3 → "+32 470 12 34 56"
func groupPhoneDigits(s string, prefixLen int) string {
	prefix := ""
	rest := s
	if prefixLen > 0 && len(s) >= prefixLen {
		prefix = s[:prefixLen]
		rest = s[prefixLen:]
	} else if strings.HasPrefix(s, "+") {
		// Fallback: scan for prefix end.
		i := 1
		for i < len(s) && s[i] >= '0' && s[i] <= '9' && i <= 4 {
			i++
		}
		prefix = s[:i]
		rest = s[i:]
	}

	// Group remaining digits: first group of 3, then groups of 2.
	var parts []string
	if prefix != "" {
		parts = append(parts, prefix)
	}
	if len(rest) > 3 {
		parts = append(parts, rest[:3])
		rest = rest[3:]
		for len(rest) >= 2 {
			parts = append(parts, rest[:2])
			rest = rest[2:]
		}
		if rest != "" {
			parts = append(parts, rest)
		}
	} else if rest != "" {
		parts = append(parts, rest)
	}

	return strings.Join(parts, " ")
}
