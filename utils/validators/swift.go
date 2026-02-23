package validators

import (
	"regexp"
	"strings"
)

// SWIFTResult holds the result of SWIFT/BIC validation.
type SWIFTResult struct {
	Valid       bool
	Error       string
	BankCode    string // 4-letter bank code
	CountryCode string // 2-letter country code
}

// swiftRegex matches SWIFT/BIC codes: 4 letters (bank) + 2 letters (country) + 2 alnum (location) + optional 3 alnum (branch).
var swiftRegex = regexp.MustCompile(`^[A-Z]{4}[A-Z]{2}[A-Z0-9]{2}([A-Z0-9]{3})?$`)

// SWIFT validates a SWIFT/BIC code.
// Accepts 8-character (head office) or 11-character (branch) codes.
func SWIFT(value string) SWIFTResult {
	if value == "" {
		return SWIFTResult{Valid: true}
	}

	code := strings.ToUpper(strings.TrimSpace(value))

	if len(code) != 8 && len(code) != 11 {
		return SWIFTResult{
			Valid: false,
			Error: "SWIFT code must be 8 or 11 characters",
		}
	}

	if !swiftRegex.MatchString(code) {
		return SWIFTResult{
			Valid: false,
			Error: "Invalid SWIFT code format",
		}
	}

	return SWIFTResult{
		Valid:       true,
		BankCode:    code[:4],
		CountryCode: code[4:6],
	}
}
