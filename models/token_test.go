package models

import (
	"strings"
	"testing"
)

func TestGenerateNumericCode(t *testing.T) {
	for range 100 {
		code, err := GenerateNumericCode(6)
		if err != nil {
			t.Fatalf("GenerateNumericCode() error = %v", err)
		}
		if len(code) != 6 {
			t.Fatalf("GenerateNumericCode() length = %d, want 6", len(code))
		}
		if strings.IndexFunc(code, func(character rune) bool {
			return character < '0' || character > '9'
		}) >= 0 {
			t.Fatalf("GenerateNumericCode() = %q, want only digits", code)
		}
	}
}
