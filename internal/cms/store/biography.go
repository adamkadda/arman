package store

import (
	"context"
	"fmt"

	"github.com/adamkadda/arman/internal/content"
)

type BiographyStore struct {
	db Executor
}

func NewBiographyStore(db Executor) *BiographyStore {
	return &BiographyStore{
		db: db,
	}
}

type biographyRow struct {
	content string `db:"content"`
	variant string `db:"variant"`
}

func (r *biographyRow) toBiography() content.Biography {
	return content.Biography{
		Content: r.content,
		Variant: content.BiographyVariant(r.variant),
	}
}

func (s *BiographyStore) Get(
	ctx context.Context,
	variant content.BiographyVariant,
) (*content.Biography, error) {
	query := `
	SELECT
		content
		variant
	FROM biographies
	WHERE variant = $1
	`

	pgxRows, err := s.db.Query(ctx, query, variant)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[biographyRow](pgxRows)
	if err != nil {
		return nil, err
	}

	biography := row.toBiography()

	return &biography, nil
}

func (s *BiographyStore) Update(
	ctx context.Context,
	b content.Biography,
) (*content.Biography, error) {
	query := `
	UPDATE biographies
	SET
		content = $1
	WHERE variant = $2
	RETURNING
		content
	`

	pgxRows, err := s.db.Query(ctx, query,
		b.Content,
		b.Variant,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[biographyRow](pgxRows)
	if err != nil {
		return nil, err
	}

	biography := row.toBiography()

	return &biography, nil
}
