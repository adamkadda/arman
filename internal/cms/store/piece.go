package store

import (
	"context"
	"fmt"

	"github.com/adamkadda/arman/internal/cms/model"
	"github.com/adamkadda/arman/internal/content"
)

type PostgresPieceStore struct {
	db Executor
}

func NewPostgresPieceStore(db Executor) *PostgresPieceStore {
	return &PostgresPieceStore{
		db: db,
	}
}

type pieceRow struct {
	pieceID         int    `db:"piece_id"`
	pieceTitle      string `db:"piece_title"`
	composerID      int    `db:"composer_id"`
	programme_count int    `db:"programme_count"`
}

func (r *pieceRow) toPiece() content.Piece {
	return content.Piece{
		ID:         r.pieceID,
		Title:      r.pieceTitle,
		ComposerID: r.composerID,
	}
}

func (r *pieceRow) toPieceWithDetails() model.PieceWithDetails {
	return model.PieceWithDetails{
		Piece:          r.toPiece(),
		ProgrammeCount: r.programme_count,
	}
}

func (s *PostgresPieceStore) Get(
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

func (s *PostgresPieceStore) GetWithDetails(
	ctx context.Context,
	id int,
) (*model.PieceWithDetails, error) {
	query := `
	SELECT
		piece_id,
		piece_title,
		composer_id
	COALESCE(pp.programme_count, 0) AS programme_count
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

func (s *PostgresPieceStore) ListWithDetails(
	ctx context.Context,
) ([]model.PieceWithDetails, error) {
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

	pieces := make([]model.PieceWithDetails, len(rows))
	for i, row := range rows {
		pieces[i] = row.toPieceWithDetails()
	}

	return pieces, nil
}

func (s *PostgresPieceStore) Create(
	ctx context.Context,
	p content.Piece,
) (*content.Piece, error) {
	query := `
	INSERT INTO pieces (
		piece_title,
		composer_id
	)
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

func (s *PostgresPieceStore) Update(
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

func (s *PostgresPieceStore) Delete(
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
