package v1

import (
	"github.com/jdkato/prose/v2"
	"gonum.org/v1/gonum/floats"
	"log"
	"log/slog"
	"math"
	"miniflux.app/v2/internal/model"
	"miniflux.app/v2/internal/similarity/story"
	"strings"
)

type Comparator interface {
	CalculateSimilarity(items model.Entries) ([]*model.EntrySimilar, error)
}

type compare struct {
	threshold float64
}

// NewComparator returns a new Comparator.
func NewComparator(threshold float64) Comparator {
	return &compare{threshold: threshold}
}

// Compare compares two articles and returns a compare score.
func (s *compare) Compare(article1, article2 string) (float64, error) {
	panic("implement me")
}

// CalculateSimilarity calculates the compare between a list of entries.
func (s *compare) CalculateSimilarity(items model.Entries) ([]*model.EntrySimilar, error) {
	stories := s.extractStories(items)
	tfidf, uniqueTokens := s.computeTFIDF(stories)
	groups := s.groupSimilarStories(stories, tfidf, uniqueTokens, s.threshold) // Increased threshold for better accuracy
	entrySimilars := make([]*model.EntrySimilar, 0)

	for i, stories := range groups {
		slog.Debug("Group", slog.Int("group", i+1), slog.Int("stories", len(stories)))
		for _, stor := range stories {
			for _, similar := range stor.Similar {
				if stor.ID == similar.Source.ID {
					continue
				}
				entrySimilars = append(entrySimilars, &model.EntrySimilar{
					EntryID:        stor.ID,
					SimilarEntryID: similar.Source.ID,
					Similarity:     similar.Similarity,
				})
				slog.Debug("Creating Simularity", slog.String("link", stor.Link), slog.String("link", similar.Source.Link), slog.Float64("compare", similar.Similarity))
			}
		}
	}
	return entrySimilars, nil
}

func (s *compare) extractStories(items []*model.Entry) []*story.Story {
	stories := []*story.Story{}
	for _, item := range items {
		stories = append(stories, story.FromEntry(item))
	}
	return stories
}

func (s *compare) tokenize(text string) []string {
	doc, err := prose.NewDocument(text)
	if err != nil {
		log.Fatalf("Failed to tokenize text: %v", err)
	}
	tokens := []string{}
	for _, tok := range doc.Tokens() {
		normalized := strings.ToLower(tok.Text)
		if _, found := stopwords[normalized]; !found {
			tokens = append(tokens, normalized)
		}
	}
	return tokens
}

func (s *compare) computeTFIDF(stories []*story.Story) (map[string]map[string]float64, []string) {
	tf := make(map[string]map[string]float64)
	df := make(map[string]int)
	tokensPerStory := make(map[string][]string)
	for _, stor := range stories {
		tokens := s.tokenize(stor.Content)
		tokensPerStory[stor.Link] = tokens
		tf[stor.Link] = make(map[string]float64)
		for _, token := range tokens {
			tf[stor.Link][token]++
		}
		for token := range tf[stor.Link] {
			df[token]++
		}
	}
	idf := make(map[string]float64)
	for token, count := range df {
		idf[token] = math.Log(float64(len(stories)) / float64(count))
	}
	tfidf := make(map[string]map[string]float64)
	for link, tfs := range tf {
		tfidf[link] = make(map[string]float64)
		for token, freq := range tfs {
			tfidf[link][token] = freq * idf[token]
		}
	}
	uniqueTokens := make([]string, 0, len(df))
	for token := range df {
		uniqueTokens = append(uniqueTokens, token)
	}
	return tfidf, uniqueTokens
}

func (s *compare) cosineSimilarity(vecA, vecB []float64) float64 {
	return floats.Dot(vecA, vecB) / (math.Sqrt(floats.Dot(vecA, vecA)) * math.Sqrt(floats.Dot(vecB, vecB)))
}

func (s *compare) compareStories(storyA, storyB *story.Story, tfidf map[string]map[string]float64, uniqueTokens []string) float64 {
	vecA := make([]float64, len(uniqueTokens))
	vecB := make([]float64, len(uniqueTokens))
	for i, token := range uniqueTokens {
		vecA[i] = tfidf[storyA.Link][token]
		vecB[i] = tfidf[storyB.Link][token]
	}
	return s.cosineSimilarity(vecA, vecB)
}

func (s *compare) groupSimilarStories(stories []*story.Story, tfidf map[string]map[string]float64, uniqueTokens []string, threshold float64) [][]*story.Story {
	groups := [][]*story.Story{}
	visited := make(map[string]bool)

	for _, store := range stories {
		if visited[store.Link] {
			continue
		}
		group := []*story.Story{store}
		visited[store.Link] = true
		for _, otherStory := range stories {
			if visited[otherStory.Link] {
				continue
			}
			similarity := s.compareStories(store, otherStory, tfidf, uniqueTokens)
			//slog.Debug("Comparing stories vs threshold",
			//	slog.String("story1", story.Title),
			//	slog.String("story2", otherStory.Title),
			//	slog.Float64("compare", compare),
			//)
			if similarity >= threshold {
				store.Similar = append(store.Similar, &story.Similar{
					Source:     otherStory,
					Similarity: similarity,
				})
				group = append(group, otherStory)
				visited[otherStory.Link] = true
			}
		}
		if len(group) > 1 { // Only add groups with multiple stories
			groups = append(groups, group)
		}
	}
	return groups
}
