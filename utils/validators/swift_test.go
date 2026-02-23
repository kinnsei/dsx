package validators

import "testing"

func TestSWIFT(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"empty", "", true},
		{"valid 8 char", "DEUTDEFF", true},
		{"valid 11 char", "DEUTDEFF500", true},
		{"valid lowercase", "deutdeff", true},
		{"too short", "DEUT", false},
		{"wrong length 9", "DEUTDEFF5", false},
		{"invalid chars", "DEUT12FF", false},
		{"numbers in bank code", "1234DEFF", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := SWIFT(tt.input)
			if r.Valid != tt.valid {
				t.Errorf("SWIFT(%q) valid=%v, want %v (error: %s)", tt.input, r.Valid, tt.valid, r.Error)
			}
		})
	}
}
