package story

import (
	"miniflux.app/v2/internal/model"
	"regexp"
	"strings"
)

type Story struct {
	ID          int64
	Title       string
	Link        string
	Description string
	Content     string
	Similar     []*Similar
}

type Similar struct {
	Source     *Story
	Similarity float64
}

// FromEntry - Extract story
func FromEntry(entry *model.Entry) *Story {
	s := Story{}
	content := RemoveNonWordsAndPunctuation(RemoveHTMLTags(entry.Title + " " + entry.Content))
	s.ID = entry.ID
	s.Title = entry.Title
	s.Link = entry.URL
	s.Description = entry.Content
	s.Content = content
	return &s
}

// Generate string to compare
func (s *Story) toString() string {
	return s.Content
}

// RemoveHTMLTags removes all HTML tags from the given text.
func RemoveHTMLTags(text string) string {
	// Define a regular expression to match HTML tags.
	re := regexp.MustCompile("<.*?>")
	// Replace all HTML tags with an empty string.
	cleanText := re.ReplaceAllString(text, "")
	return cleanText
}

// RemoveNonWordsAndPunctuation removes non-word characters and punctuation marks, and lowercases the text.
func RemoveNonWordsAndPunctuation(text string) string {
	// Define a regular expression to match non-word characters and punctuation marks.
	re := regexp.MustCompile("[^\\w\\s]+")
	// Replace all non-word characters and punctuation marks with an empty string.
	cleanText := re.ReplaceAllString(text, "")
	// Convert the text to lowercase.
	cleanText = strings.ToLower(cleanText)
	// Replace multiple spaces with a single space.
	re = regexp.MustCompile("\\s+")
	cleanText = re.ReplaceAllString(cleanText, " ")
	return strings.Trim(cleanText, " ")
}
