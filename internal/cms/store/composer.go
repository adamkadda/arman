package store

import (
	"context"
	"fmt"

	"github.com/adamkadda/arman/internal/cms/models"
	"github.com/adamkadda/arman/internal/content"
)

type ComposerStore struct {
	db Executor
}

func NewComposerStore(db Executor) *ComposerStore {
	return &ComposerStore{
		db: db,
	}
}

type composerRow struct {
	composerID int    `db:"composer_id"`
	fullName   string `db:"full_name"`
	shortName  string `db:"short_name"`
	pieceCount int    `db:"piece_count"`
}

func (r *composerRow) toComposer() content.Composer {
	return content.Composer{
		ID:        r.composerID,
		FullName:  r.fullName,
		ShortName: r.shortName,
	}
}

func (r *composerRow) toComposerWithDetails() models.ComposerWithDetails {
	return models.ComposerWithDetails{
		Composer:   r.toComposer(),
		PieceCount: r.pieceCount,
	}
}

func (s *ComposerStore) Get(
	ctx context.Context,
	id int,
) (*content.Composer, error) {
	query := `
	SELECT
		composer_id,
		full_name,
		short_name
	FROM composers
	WHERE composer_id = $1
	`

	pgxRows, err := s.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[composerRow](pgxRows)
	if err != nil {
		return nil, err
	}

	composer := row.toComposer()

	return &composer, nil
}

func (s *ComposerStore) GetWithDetails(
	ctx context.Context,
	id int,
) (*models.ComposerWithDetails, error) {
	query := `
	SELECT
		composer_id,
		full_name,
		short_name,
	COALESCE(p.piece_count, 0) AS piece_count
	FROM composers c
	LEFT JOIN (
		SELECT composer_id, COUNT(*) AS piece_count
		FROM pieces
		WHERE composer_id = $1
		GROUP BY composer_id
	) p ON p.composer_id = c.composer_id
	WHERE composer_id = $1
	`

	pgxRows, err := s.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[composerRow](pgxRows)
	if err != nil {
		return nil, err
	}

	composer := row.toComposerWithDetails()

	return &composer, nil
}

func (s *ComposerStore) ListWithDetails(
	ctx context.Context,
) ([]models.ComposerWithDetails, error) {
	query := `
	SELECT
		composer_id,
		full_name,
		short_name,
	COALESCE(p.piece_count, 0) AS piece_count
	FROM composers c
	LEFT JOIN (
		SELECT composer_id, COUNT(*) AS piece_count
		FROM pieces
		GROUP BY composer_id
	) p ON p.composer_id = c.composer_id
	ORDER BY c.composer_id
	`

	pgxRows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	rows, err := collectRows[composerRow](pgxRows)
	if err != nil {
		return nil, err
	}

	composers := make([]models.ComposerWithDetails, 0, len(rows))
	for _, row := range rows {
		composers = append(composers, row.toComposerWithDetails())
	}

	return composers, nil
}

func (s *ComposerStore) Create(
	ctx context.Context,
	c content.Composer,
) (*content.Composer, error) {
	query := `
	INSERT INTO composers (full_name, short_name)
	VALUES ($1, $2)
	RETURNING
		composer_id,
		full_name,
		short_name
	`

	pgxRows, err := s.db.Query(ctx, query,
		c.FullName,
		c.ShortName,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[composerRow](pgxRows)
	if err != nil {
		return nil, err
	}

	composer := row.toComposer()

	return &composer, nil
}

func (s *ComposerStore) Update(
	ctx context.Context,
	c content.Composer,
) (*content.Composer, error) {
	query := `
	UPDATE composers
	SET
		full_name = $1
		short_name = $2
	WHERE composer_id = $3
	RETURNING
		composer_id,
		full_name,
		short_name
	`

	pgxRows, err := s.db.Query(ctx, query,
		c.FullName,
		c.ShortName,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[composerRow](pgxRows)
	if err != nil {
		return nil, err
	}

	composer := row.toComposer()

	return &composer, nil
}

func (s *ComposerStore) Delete(
	ctx context.Context,
	id int,
) error {
	query := `
	DELETE
	FROM composers
	WHERE composer_id = $1
	`

	cmdTag, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	return checkAffected(cmdTag)
}
