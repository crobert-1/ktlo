package config

import (
	"testing"
	"time"
)

func TestParseLookbackDuration(t *testing.T) {
	tests := []struct {
		name      string
		lookback  string
		expected  time.Duration
		expectErr bool
	}{
		{
			name:     "minutes",
			lookback: "30m",
			expected: 30 * time.Minute,
		},
		{
			name:     "hours",
			lookback: "12h",
			expected: 12 * time.Hour,
		},
		{
			name:     "days",
			lookback: "7d",
			expected: 7 * 24 * time.Hour,
		},
		{
			name:      "empty",
			expectErr: true,
		},
		{
			name:      "zero duration",
			lookback:  "0s",
			expectErr: true,
		},
		{
			name:      "negative duration",
			lookback:  "-1h",
			expectErr: true,
		},
		{
			name:      "invalid duration",
			lookback:  "tomorrow",
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := parseLookbackDuration(test.lookback)
			if (err != nil) != test.expectErr {
				t.Fatalf("parseLookbackDuration(%q) error = %v, want error %t", test.lookback, err, test.expectErr)
			}
			if err == nil && actual != test.expected {
				t.Errorf("parseLookbackDuration(%q) = %s, want %s", test.lookback, actual, test.expected)
			}
		})
	}
}
