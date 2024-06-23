package similarity

import (
	"math"
	"strings"
)

// TextComparer is the struct that holds the methods to compare texts
type TextComparer struct{}

// NewTextComparer creates a new instance of TextComparer
func NewTextComparer() *TextComparer {
	return &TextComparer{}
}

// Tokenize splits text into lowercase words
func (tc *TextComparer) Tokenize(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

// TermFrequency computes the term frequency of words in a text
func (tc *TextComparer) TermFrequency(text string) map[string]float64 {
	tf := make(map[string]float64)
	words := tc.Tokenize(text)
	for _, word := range words {
		tf[word]++
	}
	for word := range tf {
		tf[word] /= float64(len(words))
	}
	return tf
}

// InverseDocumentFrequency computes the inverse document frequency for words across documents
func (tc *TextComparer) InverseDocumentFrequency(texts []string) map[string]float64 {
	idf := make(map[string]float64)
	totalDocuments := float64(len(texts))
	wordDocCount := make(map[string]float64)

	for _, text := range texts {
		words := tc.Tokenize(text)
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

// TFIDFVector computes the TF-IDF vector for a text
func (tc *TextComparer) TFIDFVector(text string, idf map[string]float64) map[string]float64 {
	tf := tc.TermFrequency(text)
	tfidf := make(map[string]float64)
	for word, tfVal := range tf {
		tfidf[word] = tfVal * idf[word]
	}
	return tfidf
}

// CosineSimilarity computes the cosine similarity between two vectors
func (tc *TextComparer) CosineSimilarity(vec1, vec2 map[string]float64) float64 {
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

// Compare compares two articles and returns their similarity score
func (tc *TextComparer) Compare(article1, article2 string) float64 {
	texts := []string{article1, article2}
	idf := tc.InverseDocumentFrequency(texts)
	tfidf1 := tc.TFIDFVector(article1, idf)
	tfidf2 := tc.TFIDFVector(article2, idf)
	return tc.CosineSimilarity(tfidf1, tfidf2)
}
