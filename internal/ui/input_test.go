package ui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSanitizeCSVPath(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{name: "plain path", in: "./usage.example.csv", out: "./usage.example.csv"},
		{name: "double quoted path", in: "\"./usage.example.csv\"", out: "./usage.example.csv"},
		{name: "single quoted path", in: "'./usage.example.csv'", out: "./usage.example.csv"},
		{name: "quoted path with spaces", in: "  \"./my usage.csv\"  ", out: "./my usage.csv"},
		{name: "nested quotes", in: "\"'./usage.example.csv'\"", out: "./usage.example.csv"},
		{name: "mismatched quote kept", in: "\"./usage.example.csv", out: "\"./usage.example.csv"},
		{name: "empty quoted path", in: "\"\"", out: ""},
		{name: "wrapped path newlines removed", in: "/Users/me/Desktop/team-usage-2026-03-\n13.csv", out: "/Users/me/Desktop/team-usage-2026-03-13.csv"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := sanitizeCSVPath(tc.in)
			if got != tc.out {
				t.Errorf("sanitizeCSVPath(%q) = %q, want %q", tc.in, got, tc.out)
			}
		})
	}
}

func TestSanitizeCSVPathExpandsHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir() error = %v", err)
	}

	got := sanitizeCSVPath("~/Desktop/usage.csv")
	want := filepath.Join(home, "Desktop/usage.csv")
	if got != want {
		t.Errorf("sanitizeCSVPath(%q) = %q, want %q", "~/Desktop/usage.csv", got, want)
	}
}
