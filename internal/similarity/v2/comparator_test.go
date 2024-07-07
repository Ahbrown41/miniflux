package v2

import (
	"math"
	"testing"
)

func TestArticleComparator(t *testing.T) {
	comparator := NewComparator()

	tests := []struct {
		name      string
		article1  string
		article2  string
		expected  float64
		threshold float64
	}{
		{
			name:      "Identical articles",
			article1:  "The quick brown fox jumps over the lazy dog.",
			article2:  "The quick brown fox jumps over the lazy dog.",
			expected:  1.0,
			threshold: 0.001,
		},
		{
			name:      "Completely different articles",
			article1:  "The quick brown fox jumps over the lazy dog.",
			article2:  "A completely different sentence with no common words.",
			expected:  0.0,
			threshold: 0.03,
		},
		{
			name:      "Partially similar articles",
			article1:  "The quick brown fox jumps over the lazy dog.",
			article2:  "The quick brown fox leaps over the lazy dog.",
			expected:  0.8, // Expecting a high similarity but not exactly 1
			threshold: 0.1,
		},
		{
			name:      "Empty articles",
			article1:  "",
			article2:  "",
			expected:  0.0,
			threshold: 0.001,
		},
		{
			name:      "One empty article",
			article1:  "The quick brown fox jumps over the lazy dog.",
			article2:  "",
			expected:  0.0,
			threshold: 0.001,
		},
		{
			name:      "Different lengths with some common content",
			article1:  "The quick brown fox jumps over the lazy dog.",
			article2:  "The quick brown fox jumps over the lazy dog. And then it ran away quickly.",
			expected:  0.8, // Expecting a moderate to high similarity
			threshold: 0.1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			similarity, err := comparator.Compare(test.article1, test.article2)
			if err != nil {
				t.Fatalf("Error comparing articles: %v", err)
			}
			if math.Abs(similarity-test.expected) > test.threshold {
				t.Errorf("Expected similarity ~%v, but got %v", test.expected, similarity)
			}
		})
	}
}
