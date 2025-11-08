package tui

import "testing"

func TestFormatTokenCount(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		expected string
	}{
		{"small number", 123, "123"},
		{"exact thousand", 1000, "1.0K"},
		{"thousands", 1234, "1.2K"},
		{"tens of thousands", 12345, "12.3K"},
		{"hundreds of thousands", 123456, "123.5K"},
		{"exact million", 1000000, "1.0M"},
		{"millions", 1234567, "1.2M"},
		{"tens of millions", 12345678, "12.3M"},
		{"zero", 0, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTokenCount(tt.count)
			if result != tt.expected {
				t.Errorf("formatTokenCount(%d) = %s, want %s", tt.count, result, tt.expected)
			}
		})
	}
}
