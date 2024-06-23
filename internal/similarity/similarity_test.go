package similarity

import (
	"github.com/stretchr/testify/require"
	"miniflux.app/v2/internal/model"
	"testing"
)

func TestCalculateSimilarity(t *testing.T) {
	// Initialize a new similarity object
	s := NewSimilarity(0.2)

	// Define your test cases here
	testCases := []struct {
		name     string
		input    []*model.Entry
		expected []*model.EntrySimilar
	}{
		//{
		//	name:     "Test 1",
		//	input:    []*model.Entry{},      // Add your test entries here
		//	expected: model.EntrySimilars{}, // Add your expected result here
		//},
		{
			name: "Test 2",
			input: []*model.Entry{
				{
					ID:      1,
					Title:   "My Story 1",
					URL:     "https://example.com/story1",
					Content: "This is a fake story 1 that is about different things that should be similar.",
				},
				{
					ID:      2,
					Title:   "My Other Story 2",
					URL:     "https://example.com/story2",
					Content: "What will I do with this story 1 that is about different things that should be similar.",
				},
				{
					ID:      3,
					Title:   "This is a totally different story",
					URL:     "https://example.com/story3",
					Content: "This story is about something completely different.",
				},
			}, // Add your test entries here
			expected: []*model.EntrySimilar{
				{
					EntryID:        1,
					SimilarEntryID: 2,
					Similarity:     0.26042439905730275,
				},
			}, // Add your expected result here
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := s.CalculateSimilarity(tc.input)
			require.NoError(t, err)
			if !compareEntrySimilars(output, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, output)
			}
		})
	}
}

// compareEntrySimilars is a helper function to compare two slices of EntrySimilar
func compareEntrySimilars(a, b []*model.EntrySimilar) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].EntryID != b[i].EntryID || a[i].SimilarEntryID != b[i].SimilarEntryID || a[i].Similarity != b[i].Similarity {
			return false
		}
	}
	return true
}
