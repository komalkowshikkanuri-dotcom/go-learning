package service

import (
	"strings"
	"testing"
)

func TestGenerateShortCode(t *testing.T) {
	code, err := GenerateShortCode(6)
	if err != nil {
		t.Fatalf("GenerateShortCode returned an error: %v", err)
	}

	if len(code) != 6 {
		t.Fatalf("expected length 6, got %d", len(code))
	}

	for _, char := range code {
		if !strings.ContainsRune(alphabet, char) {
			t.Fatalf("unexpected character %q in code %q", char, code)
		}
	}
}

func TestGenerateShortCodeRejectsInvalidLength(t *testing.T) {
	_, err := GenerateShortCode(0)
	if err == nil {
		t.Fatal("expected an error for length 0")
	}
}