package doctor

import "testing"

func TestIsVersionOlder(t *testing.T) {
	tests := []struct {
		current  string
		required string
		expected bool
	}{
		{"1.26.0", "1.26", false},
		{"1.26.1", "1.26", false},
		{"1.27.0", "1.26", false},
		{"2.0.0", "1.26", false},
		{"1.26", "1.26", false},
		{"1.25.9", "1.26", true},
		{"1.9.0", "1.26", true},
		{"1.20", "1.26", true},
		{"1.26-rc1", "1.26", false},
	}

	for _, tt := range tests {
		got := isVersionOlder(tt.current, tt.required)
		if got != tt.expected {
			t.Errorf("isVersionOlder(%q, %q) = %v; want %v", tt.current, tt.required, got, tt.expected)
		}
	}
}
