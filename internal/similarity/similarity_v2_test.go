package similarity

import (
	"fmt"
	"testing"
)

func TestNewTextComparer(t *testing.T) {
	tc := NewTextComparer()

	article1 := `
	What Grown Fish Your first articles different..
	`

	article2 := `
	Four second article text goes here, it is similar but also quite different
	`

	similarity := tc.Compare(article1, article2)
	fmt.Printf("Similarity score: %.4f\n", similarity)
}
