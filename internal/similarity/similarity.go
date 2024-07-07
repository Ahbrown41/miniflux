package similarity

import (
	"fmt"
	"log/slog"
	"miniflux.app/v2/internal/model"
	"miniflux.app/v2/internal/similarity/pool"
	"miniflux.app/v2/internal/similarity/story"
	"miniflux.app/v2/internal/similarity/v2"
	"miniflux.app/v2/internal/storage"
	"time"
)

type calculateSimilarity struct {
	threshold float64
}

type Similarity interface {
	Calculate(builder *storage.EntryQueryBuilder) ([]*story.Story, error)
}

// NewSimilarity returns a new similarity.
func NewSimilarity(threshold float64) Similarity {
	return &calculateSimilarity{threshold: threshold}
}

// processor - Process a given loop of all stories minus the offset
func (s *calculateSimilarity) processor(builder *storage.EntryQueryBuilder, offset int) func(story1 *story.Story) (*story.Story, error) {
	return func(story1 *story.Story) (*story.Story, error) {
		subStart := time.Now()
		iterations := 0
		builder.WithOffset(offset)
		comp := v2.NewComparator()
		err := builder.EntryProcessor(func(entry2 *model.Entry) error {
			story2 := story.FromEntry(entry2)
			sim, err := comp.Compare(story1.Content, story2.Content)
			if err != nil {
				return err
			}
			if sim >= s.threshold {
				story1.Similar = append(story1.Similar, &story.Similar{
					Source:     story2,
					Similarity: sim,
				})
			}
			iterations++
			return nil
		})
		if err != nil {
			return nil, err
		}
		slog.Debug("Processed stories",
			slog.Int("offset", offset),
			slog.Int("iterations", iterations),
			slog.String("Duration", fmt.Sprintf("%s", time.Since(subStart))))
		return story1, nil
	}
}

// Calculate calculates the similarity between entries.
func (s *calculateSimilarity) Calculate(builder *storage.EntryQueryBuilder) ([]*story.Story, error) {
	threadPool := pool.NewThreadPool[*story.Story, *story.Story](5, 200)
	threadPool.Start()
	builder.WithSorting("id", "asc")
	offset := 1
	count, err := builder.CountEntries()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}
	slog.Debug("Total entries", slog.Int("count", count))

	err = builder.EntryProcessor(func(entry1 *model.Entry) error {
		thisBuilder := builder
		slog.Debug("Processing entry",
			slog.Int64("EntryID", entry1.ID),
			slog.Int("Offset", offset))
		threadPool.TaskQueue <- pool.Task[*story.Story, *story.Story]{
			ID:       offset,
			Data:     story.FromEntry(entry1),
			Function: s.processor(thisBuilder, offset),
		}
		offset++
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Collect results.
	stories := []*story.Story{}
	go func() {
		for result := range threadPool.ResultChan {
			if len(result.Value.Similar) > 0 {
				stories = append(stories, result.Value)
			}
		}
	}()

	// Collect errors.
	go func() {
		for err := range threadPool.ErrChan {
			fmt.Printf("Error: %v\n", err)
		}
	}()

	// Stop the thread pool after processing.
	threadPool.Stop()
	return stories, nil
}
