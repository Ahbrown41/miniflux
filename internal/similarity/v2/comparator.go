package v2

import (
	"math"
	"strings"
)

type Comparator interface {
	Compare(article1, article2 string) (float64, error)
}

// TextComparer is the struct that holds the methods to compare texts
type compare struct{}

// NewComparator creates a new instance of Comparator
func NewComparator() Comparator {
	return &compare{}
}

// Compare compares two articles and returns their similarity score
func (c *compare) Compare(article1, article2 string) (float64, error) {
	// Tokenize and remove stop words
	tokens1 := c.tokenizeAndRemoveStopWords(article1)
	tokens2 := c.tokenizeAndRemoveStopWords(article2)

	// Calculate TF-IDF vectors
	tfidf1 := c.calculateTFIDF(tokens1)
	tfidf2 := c.calculateTFIDF(tokens2)

	// Calculate and return cosine similarity
	similarity := c.cosineSimilarity(tfidf1, tfidf2)
	return similarity, nil
}

// tokenizeAndRemoveStopWords tokenizes a text and removes stop words
func (c *compare) tokenizeAndRemoveStopWords(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	filteredWords := make([]string, 0, len(words))
	for _, word := range words {
		if _, found := stopWords[word]; !found {
			filteredWords = append(filteredWords, word)
		}
	}
	return filteredWords
}

// tokenize splits the text into words
func (c *compare) tokenize(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

// termFrequency computes the term frequency of words in a text
func (tc *compare) termFrequency(text string) map[string]float64 {
	tf := make(map[string]float64)
	words := tc.tokenize(text)
	for _, word := range words {
		tf[word]++
	}
	for word := range tf {
		tf[word] /= float64(len(words))
	}
	return tf
}

// inverseDocumentFrequency computes the inverse document frequency for words across documents
func (tc *compare) inverseDocumentFrequency(texts []string) map[string]float64 {
	idf := make(map[string]float64)
	totalDocuments := float64(len(texts))
	wordDocCount := make(map[string]float64)

	for _, text := range texts {
		words := tc.tokenize(text)
		wordSet := make(map[string]struct{})
		for _, word := range words {
			wordSet[word] = struct{}{}
		}
		for word := range wordSet {
			wordDocCount[word]++
		}
	}

	for word, count := range wordDocCount {
		idf[word] = math.Log(totalDocuments / (1 + count))
	}

	return idf
}

// calculateTFIDF calculates the TF-IDF vector for a given text
func (c *compare) calculateTFIDF(tokens []string) map[string]float64 {
	idf := c.inverseDocumentFrequency([]string{strings.Join(tokens, " ")}) // Simplified; in practice, use a corpus-wide IDF
	tfidf := make(map[string]float64)
	for _, token := range tokens {
		tfVal, ok := c.termFrequency(strings.Join(tokens, " "))[token]
		if ok {
			tfidf[token] = tfVal * idf[token]
		}
	}
	return tfidf
}

// tfidFVector computes the TF-IDF vector for a text
func (tc *compare) tfidFVector(text string, idf map[string]float64) map[string]float64 {
	tf := tc.termFrequency(text)
	tfidf := make(map[string]float64)
	for word, tfVal := range tf {
		tfidf[word] = tfVal * idf[word]
	}
	return tfidf
}

// CosineSimilarity computes the cosine similarity between two vectors
func (tc *compare) cosineSimilarity(vec1, vec2 map[string]float64) float64 {
	var dotProduct, mag1, mag2 float64
	for word, val1 := range vec1 {
		val2 := vec2[word]
		dotProduct += val1 * val2
		mag1 += val1 * val1
	}
	for _, val2 := range vec2 {
		mag2 += val2 * val2
	}
	if mag1 == 0 || mag2 == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(mag1) * math.Sqrt(mag2))
}
