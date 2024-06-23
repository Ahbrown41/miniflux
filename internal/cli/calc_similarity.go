package cli

import (
	"log/slog"
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
		// Get all feeds for the user
		feeds, err := store.Feeds(user.ID)
		if err != nil {
			printErrorAndExit(err)
		}
		for _, feed := range feeds {
			slog.Debug("Processing feed", slog.Int64("feedID", feed.ID))
			// Get all entries for the feed
			builder := store.NewEntryQueryBuilder(user.ID)
			builder.WithFeedID(feed.ID)
			entries, err := builder.GetEntries()
			if err != nil {
				printErrorAndExit(err)
			}
			sim := similarity.NewSimilarity(similarityThreshold)
			similars, err := sim.CalculateSimilarity(entries)
			if err != nil {
				printErrorAndExit(err)
			}
			// Create similar entries
			for _, similar := range similars {
				slog.Debug("Found Similar",
					slog.Int64("entryID", similar.EntryID),
					slog.Int64("similarEntryID", similar.SimilarEntryID),
					slog.Float64("similarity", similar.Similarity),
				)
				err := store.CreateSimilarEntry(similar)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
