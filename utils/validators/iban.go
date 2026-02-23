package validators

import (
	"math/big"
	"strings"
	"unicode"
)

// IBANResult holds the result of IBAN validation.
type IBANResult struct {
	Valid       bool
	Error       string
	CountryCode string // two-letter country code (if valid format)
}

// ibanLengths maps country codes to expected IBAN lengths.
var ibanLengths = map[string]int{
	"AL": 28, "AD": 24, "AT": 20, "AZ": 28, "BH": 22,
	"BY": 28, "BE": 16, "BA": 20, "BR": 29, "BG": 22,
	"CR": 22, "HR": 21, "CY": 28, "CZ": 24, "DK": 18,
	"DO": 28, "TL": 23, "EG": 29, "SV": 28, "EE": 20,
	"FO": 18, "FI": 18, "FR": 27, "GE": 22, "DE": 22,
	"GI": 23, "GR": 27, "GL": 18, "GT": 28, "HU": 28,
	"IS": 26, "IQ": 23, "IE": 22, "IL": 23, "IT": 27,
	"JO": 30, "KZ": 20, "XK": 20, "KW": 30, "LV": 21,
	"LB": 28, "LI": 21, "LT": 20, "LU": 20, "MT": 31,
	"MR": 27, "MU": 30, "MC": 27, "MD": 24, "ME": 22,
	"NL": 18, "MK": 19, "NO": 15, "PK": 24, "PS": 29,
	"PL": 28, "PT": 25, "QA": 29, "RO": 24, "LC": 32,
	"SM": 27, "ST": 25, "SA": 24, "RS": 22, "SC": 31,
	"SK": 24, "SI": 19, "ES": 24, "SD": 18, "SE": 24,
	"CH": 21, "TN": 24, "TR": 26, "UA": 29, "AE": 23,
	"GB": 22, "VA": 22, "VG": 24,
}

// IBAN validates an IBAN (International Bank Account Number).
// Checks format, country-specific length, and MOD-97 checksum.
func IBAN(value string) IBANResult {
	if value == "" {
		return IBANResult{Valid: true}
	}

	// Remove spaces and convert to uppercase.
	iban := strings.ToUpper(strings.ReplaceAll(value, " ", ""))

	// Must be at least 5 characters (2 country + 2 check + 1 BBAN).
	if len(iban) < 5 {
		return IBANResult{Valid: false, Error: "IBAN is too short"}
	}

	// Must contain only letters and digits.
	for _, c := range iban {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return IBANResult{Valid: false, Error: "IBAN contains invalid characters"}
		}
	}

	// Country code must be two letters.
	cc := iban[:2]
	if !unicode.IsLetter(rune(cc[0])) || !unicode.IsLetter(rune(cc[1])) {
		return IBANResult{Valid: false, Error: "Invalid country code"}
	}

	// Check country-specific length.
	if expected, ok := ibanLengths[cc]; ok {
		if len(iban) != expected {
			return IBANResult{
				Valid:       false,
				Error:       "Invalid IBAN length for " + cc,
				CountryCode: cc,
			}
		}
	}

	// MOD-97 checksum: move first 4 chars to end, convert letters to digits, check mod 97 == 1.
	rearranged := iban[4:] + iban[:4]
	var digits strings.Builder
	for _, c := range rearranged {
		if unicode.IsLetter(c) {
			// A=10, B=11, ..., Z=35
			digits.WriteString(big.NewInt(int64(c - 'A' + 10)).String())
		} else {
			digits.WriteByte(byte(c))
		}
	}

	num := new(big.Int)
	num.SetString(digits.String(), 10)
	mod := new(big.Int).Mod(num, big.NewInt(97))
	if mod.Int64() != 1 {
		return IBANResult{
			Valid:       false,
			Error:       "Invalid IBAN checksum",
			CountryCode: cc,
		}
	}

	return IBANResult{Valid: true, CountryCode: cc}
}
