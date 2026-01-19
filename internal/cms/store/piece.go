package store

import (
	"context"
	"fmt"

	"github.com/adamkadda/arman/internal/cms/models"
	"github.com/adamkadda/arman/internal/content"
)

type PieceStore struct {
	db Executor
}

func NewPieceStore(db Executor) *PieceStore {
	return &PieceStore{
		db: db,
	}
}

type pieceRow struct {
	pieceID         int    `db:"piece_id"`
	piece_title     string `db:"piece_title"`
	composerID      int    `db:"composer_id"`
	programme_count int    `db:"programme_count"`
}

func (r *pieceRow) toPiece() content.Piece {
	return content.Piece{
		ID:         r.pieceID,
		Title:      r.piece_title,
		ComposerID: r.composerID,
	}
}

func (r *pieceRow) toPieceWithDetails() models.PieceWithDetails {
	return models.PieceWithDetails{
		Piece:          r.toPiece(),
		ProgrammeCount: r.programme_count,
	}
}

func (s *PieceStore) Get(
	ctx context.Context,
	id int,
) (*content.Piece, error) {
	query := `
	SELECT
		piece_id,
		piece_title,
		composer_id
	FROM pieces
	WHERE piece_id = $1
	`

	pgxRows, err := s.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[pieceRow](pgxRows)
	if err != nil {
		return nil, err
	}

	piece := row.toPiece()

	return &piece, nil
}

func (s *PieceStore) GetWithDetails(
	ctx context.Context,
	id int,
) (*models.PieceWithDetails, error) {
	query := `
	SELECT
		piece_id,
		piece_title,
		composer_id
	COALESCE() AS programme_count
	FROM pieces p
	LEFT JOIN (
		SELECT piece_id, COUNT(*) AS programme_count
		FROM programme_pieces
		WHERE piece_id = $1
		GROUP BY composer_id
	) pp on pp.piece_id = p.piece_id
	WHERE piece_id = $1
	`

	pgxRows, err := s.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[pieceRow](pgxRows)
	if err != nil {
		return nil, err
	}

	piece := row.toPieceWithDetails()

	return &piece, nil
}

func (s *PieceStore) ListWithDetails(
	ctx context.Context,
) ([]models.PieceWithDetails, error) {
	query := `
	SELECT
		piece_id,
		piece_title,
		composer_id
	COALESCE() AS programme_count
	FROM pieces p
	LEFT JOIN (
		SELECT piece_id, COUNT(*) AS programme_count
		FROM programme_pieces
		GROUP BY composer_id
	) pp on pp.piece_id = p.piece_id
	ORDER BY p.piece_id
	`

	pgxRows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	rows, err := collectRows[pieceRow](pgxRows)
	if err != nil {
		return nil, err
	}

	pieces := make([]models.PieceWithDetails, 0, len(rows))
	for _, row := range rows {
		pieces = append(pieces, row.toPieceWithDetails())
	}

	return pieces, nil
}

func (s *PieceStore) Create(
	ctx context.Context,
	p content.Piece,
) (*content.Piece, error) {
	query := `
	INSERT INTO pieces (piece_title, composer_id)
	VALUES ($1, $2)
	RETURNING
		piece_id,
		piece_title,
		composer_id
	`

	pgxRows, err := s.db.Query(ctx, query,
		p.Title,
		p.ComposerID,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[pieceRow](pgxRows)
	if err != nil {
		return nil, err
	}

	piece := row.toPiece()

	return &piece, nil
}

func (s *PieceStore) Update(
	ctx context.Context,
	p content.Piece,
) (*content.Piece, error) {
	query := `
	UPDATE pieces
	SET
		piece_title = $1
		composer_id = $2
	WHERE piece_id = $3
	RETURNING
		piece_id,
		piece_title,
		composer_id
	`

	pgxRows, err := s.db.Query(ctx, query,
		p.Title,
		p.ComposerID,
		p.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[pieceRow](pgxRows)
	if err != nil {
		return nil, err
	}

	piece := row.toPiece()

	return &piece, nil
}

func (s *PieceStore) Delete(
	ctx context.Context,
	id int,
) error {
	query := `
	DELETE
	FROM pieces 
	WHERE piece_id = $1
	`

	cmdTag, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	return checkAffected(cmdTag)
}
