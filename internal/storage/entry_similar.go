// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package storage // import "miniflux.app/v2/internal/storage"

import (
	"fmt"
	"miniflux.app/v2/internal/model"
)

// CreateSimilarEntry add a new entry similar.
func (s *Storage) CreateSimilarEntry(similar *model.EntrySimilar) error {
	query := `
		INSERT INTO entry_similar
			(entry_id, similar_entry_id, similarity)
		VALUES
			($1, $2, $3)
		ON CONFLICT DO NOTHING
	`
	_, err := s.db.Exec(
		query,
		similar.EntryID,
		similar.SimilarEntryID,
		similar.Similarity,
	)

	if err != nil {
		return fmt.Errorf(`store: unable to create similar_entry (%d -> %d): %s`, similar.EntryID, similar.SimilarEntryID, err)
	}
	return nil
}

// FindSimilarEntries returns a list of similar entries.
func (s *Storage) FindSimilarEntries(entryID int64) ([]*model.EntrySimilar, error) {
	query := `
		SELECT
			entry_id, similar_entry_id, similarity
		FROM
			entry_similar
		WHERE
			entry_id=$1
	`
	rows, err := s.db.Query(query, entryID)
	if err != nil {
		return nil, fmt.Errorf(`store: unable to fetch similar entries: %v`, err)
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			fmt.Printf(`store: unable to close similar entries rows: %v`, err)
		}
	}()
	similar := make([]*model.EntrySimilar, 0)
	err = rows.Scan(&similar)
	if err != nil {
		return nil, fmt.Errorf(`store: unable to scan similar entries: %v`, err)
	}
	return similar, nil
}
