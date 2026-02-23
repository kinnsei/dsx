package validators

import "testing"

func TestIBAN(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"empty", "", true},
		{"valid GB", "GB29 NWBK 6016 1331 9268 19", true},
		{"valid DE", "DE89370400440532013000", true},
		{"valid NL", "NL91ABNA0417164300", true},
		{"too short", "GB29", false},
		{"bad checksum", "GB00 NWBK 6016 1331 9268 19", false},
		{"wrong length for country", "GB29 NWBK 6016 13", false},
		{"invalid chars", "GB29!NWBK60161331926819", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := IBAN(tt.input)
			if r.Valid != tt.valid {
				t.Errorf("IBAN(%q) valid=%v, want %v (error: %s)", tt.input, r.Valid, tt.valid, r.Error)
			}
		})
	}
}
