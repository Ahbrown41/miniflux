package similarity

import (
	"github.com/jdkato/prose/v2"
	"gonum.org/v1/gonum/floats"
	"log"
	"log/slog"
	"math"
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

type similarity struct {
	threshold float64
}

type Similarity interface {
	CalculateSimilarity(items model.Entries) ([]*model.EntrySimilar, error)
}

// NewSimilarity returns a new similarity.
func NewSimilarity(threshold float64) Similarity {
	return &similarity{threshold: threshold}
}

// CalculateSimilarity calculates the similarity between a list of entries.
func (s *similarity) CalculateSimilarity(items model.Entries) ([]*model.EntrySimilar, error) {
	stories := s.extractStories(items)
	tfidf, uniqueTokens := s.computeTFIDF(stories)
	groups := s.groupSimilarStories(stories, tfidf, uniqueTokens, s.threshold) // Increased threshold for better accuracy
	entrySimilars := make([]*model.EntrySimilar, 0)

	for i, stories := range groups {
		slog.Debug("Group", slog.Int("group", i+1), slog.Int("stories", len(stories)))
		for _, story := range stories {
			for _, similar := range story.Similar {
				if story.ID == similar.Source.ID {
					continue
				}
				entrySimilars = append(entrySimilars, &model.EntrySimilar{
					EntryID:        story.ID,
					SimilarEntryID: similar.Source.ID,
					Similarity:     similar.Similarity,
				})
				slog.Debug("Creating Simularity", slog.String("link", story.Link), slog.String("link", similar.Source.Link), slog.Float64("similarity", similar.Similarity))
			}
		}
	}
	return entrySimilars, nil
}

// removeHTMLTags removes HTML tags from a string
func (s *similarity) removeHTMLTags(text string) string {
	re := regexp.MustCompile("<[^>]*>")
	return re.ReplaceAllString(text, "")
}

// PreprocessText removes non-word characters, HTML tags, and lowercases the text
func (s *similarity) preprocessText(text string) string {
	text = s.removeHTMLTags(text)
	re := regexp.MustCompile(`[\W_]+`)
	text = re.ReplaceAllString(text, " ")
	return strings.ToLower(text)
}

func (s *similarity) extractStories(items []*model.Entry) []*Story {
	stories := []*Story{}
	for _, item := range items {
		content := s.preprocessText(item.Title + " " + item.Content)
		stories = append(stories, &Story{
			ID:          item.ID,
			Title:       item.Title,
			Link:        item.URL,
			Description: item.Content,
			Content:     content,
		})
	}
	return stories
}

func (s *similarity) tokenize(text string) []string {
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

func (s *similarity) computeTFIDF(stories []*Story) (map[string]map[string]float64, []string) {
	tf := make(map[string]map[string]float64)
	df := make(map[string]int)
	tokensPerStory := make(map[string][]string)
	for _, story := range stories {
		tokens := s.tokenize(story.Content)
		tokensPerStory[story.Link] = tokens
		tf[story.Link] = make(map[string]float64)
		for _, token := range tokens {
			tf[story.Link][token]++
		}
		for token := range tf[story.Link] {
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

func (s *similarity) cosineSimilarity(vecA, vecB []float64) float64 {
	return floats.Dot(vecA, vecB) / (math.Sqrt(floats.Dot(vecA, vecA)) * math.Sqrt(floats.Dot(vecB, vecB)))
}

func (s *similarity) compareStories(storyA, storyB *Story, tfidf map[string]map[string]float64, uniqueTokens []string) float64 {
	vecA := make([]float64, len(uniqueTokens))
	vecB := make([]float64, len(uniqueTokens))
	for i, token := range uniqueTokens {
		vecA[i] = tfidf[storyA.Link][token]
		vecB[i] = tfidf[storyB.Link][token]
	}
	return s.cosineSimilarity(vecA, vecB)
}

func (s *similarity) groupSimilarStories(stories []*Story, tfidf map[string]map[string]float64, uniqueTokens []string, threshold float64) [][]*Story {
	groups := [][]*Story{}
	visited := make(map[string]bool)

	for _, story := range stories {
		if visited[story.Link] {
			continue
		}
		group := []*Story{story}
		visited[story.Link] = true
		for _, otherStory := range stories {
			if visited[otherStory.Link] {
				continue
			}
			similarity := s.compareStories(story, otherStory, tfidf, uniqueTokens)
			//slog.Debug("Comparing stories vs threshold",
			//	slog.String("story1", story.Title),
			//	slog.String("story2", otherStory.Title),
			//	slog.Float64("similarity", similarity),
			//)
			if similarity >= threshold {
				story.Similar = append(story.Similar, &Similar{
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
