package v1

import (
	"miniflux.app/v2/internal/model"
	"testing"
)

func TestCalculateSimilarity_EmptyInput(t *testing.T) {
	comparator := NewComparator(0.5) // Assuming a threshold of 0.5 for similarity
	entries := []*model.Entry{}

	similars, err := comparator.CalculateSimilarity(entries)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(similars) != 0 {
		t.Errorf("Expected no similarities, got %d", len(similars))
	}
}

func TestCalculateSimilarity_SingleEntry(t *testing.T) {
	comparator := NewComparator(0.5)
	entries := []*model.Entry{createEntry("Single story content")}

	similars, err := comparator.CalculateSimilarity(entries)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(similars) != 0 {
		t.Errorf("Expected no similarities, got %d", len(similars))
	}
}

func TestCalculateSimilarity_MultipleEntriesNoSimilarities(t *testing.T) {
	comparator := NewComparator(0.5)
	entries := []*model.Entry{
		createEntry("First unique story"),
		createEntry("Second unique story"),
	}

	similars, err := comparator.CalculateSimilarity(entries)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(similars) != 0 {
		t.Errorf("Expected no similarities, got %d", len(similars))
	}
}

func TestCalculateSimilarity_MultipleEntriesWithSimilarities(t *testing.T) {
	comparator := NewComparator(0.5)
	entries := []*model.Entry{
		createEntry("Similar story one"),
		createEntry("Similar story two"),
		createEntry("Unique story"),
	}

	similars, err := comparator.CalculateSimilarity(entries)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(similars) == 0 {
		t.Errorf("Expected similarities, got none")
	}
}

func TestCalculateSimilarity_IdenticalEntries(t *testing.T) {
	comparator := NewComparator(0.5)
	entries := []*model.Entry{
		createEntry("Identical story"),
		createEntry("Identical story"),
	}

	similars, err := comparator.CalculateSimilarity(entries)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(similars) != 1 || similars[0].Similarity != 1.0 {
		t.Errorf("Expected identical entries to be fully similar")
	}
}

// Helper function to create model.Entry instances for testing
func createEntry(content string) *model.Entry {
	return &model.Entry{Content: content}
}
