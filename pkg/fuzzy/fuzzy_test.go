// fuzzy/fuzzy_test.go
package fuzzy

import (
	"testing"
)

func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		a        string
		b        string
		expected int
	}{
		// Identical strings
		{"", "", 0},
		{"test", "test", 0},

		// One string empty
		{"", "hello", 5},
		{"world", "", 5},

		// Single edit operations
		{"kitten", "sitten", 1},  // substitution
		{"kitten", "kitte", 1},   // deletion
		{"kitten", "kittens", 1}, // insertion

		// Multiple edits
		{"flaw", "lawn", 2},
		{"intention", "execution", 5},
	}

	for _, tt := range tests {
		result := LevenshteinDistance(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("LevenshteinDistance(%q, %q) = %d; want %d", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestSimilarityPercentage(t *testing.T) {
	tests := []struct {
		a        string
		b        string
		expected float64
	}{
		{"", "", 100.0},
		{"test", "test", 100.0},

		{"", "hello", 0.0},
		{"world", "", 0.0},

		{"kitten", "sitten", 83},
		{"kitten", "kitte", 83.0},
		{"kitten", "kittens", 85},

		{"flaw", "lawn", 50.0},
		{"intention", "execution", 44.44},
	}

	for _, tt := range tests {
		result := SimilarityPercentage(tt.a, tt.b)
		if !approxEqual(result, tt.expected, 1) {
			t.Errorf("SimilarityPercentage(%q, %q) = %.2f; want %.2f", tt.a, tt.b, result, tt.expected)
		}
	}
}

func approxEqual(a, b, epsilon float64) bool {
	if a > b {
		return a-b < epsilon
	}
	return b-a < epsilon
}
