package cli

import (
	"log/slog"
	"miniflux.app/v2/internal/model"
	"miniflux.app/v2/internal/similarity"
	"miniflux.app/v2/internal/storage"
)

func calcSimilarity(store *storage.Storage, similarityThreshold float64) error {
	slog.Debug("Calculating similarity",
		slog.Float64("threshold", similarityThreshold),
		slog.String("action", "calc_similarity.go:calcSimilarity()"))
	// Get all users
	users, err := store.Users()
	if err != nil {
		printErrorAndExit(err)
	}
	for _, user := range users {
		slog.Debug("Processing user", slog.Int64("userID", user.ID))
		// Calculate Similar
		builder := store.NewEntryQueryBuilder(user.ID)
		stories, err := similarity.NewSimilarity(similarityThreshold).Calculate(builder)
		if err != nil {
			printErrorAndExit(err)
		}

		// Create similar entries
		for _, story := range stories {
			for _, similar := range story.Similar {
				sims, err := store.FindSimilarEntries(story.ID)
				if err != nil {
					return err
				}
				found := false
				for _, sim := range sims {
					if sim.EntryID == story.ID {
						found = true
						slog.Debug("Found existing similar story", slog.Int64("storyID", story.ID))
					}
				}
				if !found {
					slog.Debug("Found new similar",
						slog.Int64("entryID", story.ID),
						slog.Int64("similarEntryID", similar.Source.ID),
						slog.Float64("similarity", similar.Similarity),
					)
					err := store.CreateSimilarEntry(&model.EntrySimilar{
						EntryID:        story.ID,
						SimilarEntryID: similar.Source.ID,
						Similarity:     similar.Similarity,
					})
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
