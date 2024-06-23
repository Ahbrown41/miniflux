// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package model // import "miniflux.app/v2/internal/model"

// EntrySimilar represents a similar entry in the application.
type EntrySimilar struct {
	// EntryID is the unique identifier of the entry.
	EntryID int64 `json:"entry_id"`
	// SimilarEntryID is the unique identifier of the similar entry.
	SimilarEntryID int64 `json:"similar_entry_id"`
	// Similarity is the similarity score between the two entries.
	Similarity float64 `json:"similarity"`
}
