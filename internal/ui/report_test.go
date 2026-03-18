package ui

import (
	"fmt"
	"testing"
)

func TestFormatIntWithCommas(t *testing.T) {
	tests := []struct {
		input int64
		want  string
	}{
		{input: 0, want: "0"},
		{input: 12, want: "12"},
		{input: 12123, want: "12,123"},
		{input: 12123550, want: "12,123,550"},
		{input: -500, want: "-500"},
		{input: -1234567, want: "-1,234,567"},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d", tc.input), func(t *testing.T) {
			got := formatIntWithCommas(tc.input)
			if got != tc.want {
				t.Errorf("formatIntWithCommas(%d) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
