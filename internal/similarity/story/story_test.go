package story

import (
	"testing"
)

// Test cases for RemoveHTMLTags
func TestRemoveHTMLTags(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<p>This is a <strong>sample</strong> text with <a href=\"http://example.com\">HTML</a> tags.</p>", "This is a sample text with HTML tags."},
		{"<div><p>Nested <span>tags</span> example.</p></div>", "Nested tags example."},
		{"No HTML tags here!", "No HTML tags here!"},
		{"<a href='#'>Link</a>", "Link"},
		{"<img src='image.jpg' alt='image'/>", ""},
	}

	for _, test := range tests {
		result := RemoveHTMLTags(test.input)
		if result != test.expected {
			t.Errorf("RemoveHTMLTags(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

// Test cases for RemoveNonWordsAndPunctuation
func TestRemoveNonWordsAndPunctuation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"This is a sample text with punctuation, and CAPITAL letters!", "this is a sample text with punctuation and capital letters"},
		{"Multiple    spaces   should be  reduced.", "multiple spaces should be reduced"},
		{"Numbers 123 and symbols #$%^ should be removed.", "numbers 123 and symbols should be removed"},
		{"MixOfCAPS and lower.", "mixofcaps and lower"},
		{"Special characters: @#$%^&*()!", "special characters"},
	}

	for _, test := range tests {
		result := RemoveNonWordsAndPunctuation(test.input)
		if result != test.expected {
			t.Errorf("RemoveNonWordsAndPunctuation(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

// Test cases for combined functionality
func TestCombinedFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<p>This is a <strong>sample</strong> text with <a href=\"http://example.com\">HTML</a> tags, punctuation! and CAPITAL letters.</p>", "this is a sample text with html tags punctuation and capital letters"},
		{"<div>Text with <span>nested <strong>HTML</strong></span> tags and <em>special characters</em> like @#$%.</div>", "text with nested html tags and special characters like"},
		{"<a href='#'>Link</a> with words and 123 numbers.", "link with words and 123 numbers"},
		{"<img src='image.jpg' alt='image'/> Description with symbols #$%^&*", "description with symbols"},
		{"No HTML but special characters: @#$%^&*()!", "no html but special characters"},
	}

	for _, test := range tests {
		htmlCleaned := RemoveHTMLTags(test.input)
		finalResult := RemoveNonWordsAndPunctuation(htmlCleaned)
		if finalResult != test.expected {
			t.Errorf("CombinedFunctions(%q) = %q; want %q", test.input, finalResult, test.expected)
		}
	}
}
