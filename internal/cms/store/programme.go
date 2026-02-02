package store

import (
	"context"
	"fmt"

	"github.com/adamkadda/arman/internal/cms/model"
	"github.com/adamkadda/arman/internal/content"
)

type ProgrammeStore struct {
	db Executor
}

func NewProgrammeStore(db Executor) *ProgrammeStore {
	return &ProgrammeStore{
		db: db,
	}
}

type programmeRow struct {
	programmeID    int    `db:"programme_id"`
	programmeTitle string `db:"programme_title"`
	eventCount     int    `db:"event_count"`
}

func (r *programmeRow) toProgramme() content.Programme {
	return content.Programme{
		ID:    r.programmeID,
		Title: r.programmeTitle,
	}
}

func (r *programmeRow) toProgrammeWithDetails() model.ProgrammeWithDetails {
	return model.ProgrammeWithDetails{
		Programme:  r.toProgramme(),
		EventCount: r.eventCount,
	}
}

func (s *ProgrammeStore) Get(
	ctx context.Context,
	id int,
) (*content.Programme, error) {
	query := `
	SELECT
		programme_id,
		programme_title
	FROM programmes
	WHERE programme_id = $1
	`

	pgxRows, err := s.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[programmeRow](pgxRows)
	if err != nil {
		return nil, err
	}

	programme := row.toProgramme()

	return &programme, nil
}

func (s *ProgrammeStore) GetWithDetails(
	ctx context.Context,
	id int,
) (*model.ProgrammeWithDetails, error) {
	query := `
	SELECT
		programme_id,
		programme_title
	COALESCE(e.event_count, 0) AS event_count
	FROM programmes p
	LEFT JOIN (
	SELECT programme_id, COUNT(*) AS event_count
	FROM events
	WHERE venue_id = $1 AND status = "published"
	GROUP BY programme_id
	) e ON e.programme_id = p.programme_id
	WHERE programme_id = $1
	`

	pgxRows, err := s.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[programmeRow](pgxRows)
	if err != nil {
		return nil, err
	}

	programme := row.toProgrammeWithDetails()

	return &programme, nil
}

// ListWithDetails returns a list of programmes with additional details. This method
// was the alternative to an N+1 query approach of making multiple queries for each
// Programme. It's a tradeoff between performance and adherence to a (strict) clean
// separation of concerns.
func (s *ProgrammeStore) ListWithDetails(
	ctx context.Context,
) ([]model.ProgrammeWithDetails, error) {
	query := `
	SELECT
		programme_id,
		programme_title
	COALESCE(e.event_count, 0) AS event_count
	FROM programmes p
	LEFT JOIN (
	SELECT programme_id, COUNT(*) AS event_count
	FROM events
	WHERE status = "published"
	GROUP BY programme_id
	) e ON e.programme_id = p.programme_id
	ORDER BY p.programme_id
	`

	pgxRows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	rows, err := collectRows[programmeRow](pgxRows)
	if err != nil {
		return nil, err
	}

	programmes := make([]model.ProgrammeWithDetails, len(rows))
	for i, row := range rows {
		programmes[i] = row.toProgrammeWithDetails()
	}

	return programmes, nil
}

func (s *ProgrammeStore) Create(
	ctx context.Context,
	p content.Programme,
) (*content.Programme, error) {
	query := `
	INSERT INTO programmes (
		programme_title
	)
	VALUES ($1)
	RETURNING
		programme_id
		programme_title
	`

	pgxRows, err := s.db.Query(ctx, query,
		p.Title,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[programmeRow](pgxRows)
	if err != nil {
		return nil, err
	}

	programme := row.toProgramme()

	return &programme, nil
}

func (s *ProgrammeStore) Update(
	ctx context.Context,
	p content.Programme,
) (*content.Programme, error) {
	query := `
	UPDATE programmes
	SET
		programme_title = $1
	WHERE programme_id = $2
	RETURNING
		programme_id,
		programme_title
	`

	pgxRows, err := s.db.Query(ctx, query,
		p.Title,
		p.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[programmeRow](pgxRows)
	if err != nil {
		return nil, err
	}

	programme := row.toProgramme()

	return &programme, nil
}

func (s *ProgrammeStore) Delete(
	ctx context.Context,
	id int,
) error {
	query := `
	DELETE
	FROM programmes
	WHERE programme_id = $1
	`

	cmdTag, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	return checkAffected(cmdTag)
}

type ProgrammePieceStore struct {
	db Executor
}

func NewProgrammePieceStore(db Executor) *ProgrammePieceStore {
	return &ProgrammePieceStore{
		db: db,
	}
}

type programmePieceRow struct {
	pieceID    int    `db:"piece_id"`
	pieceTitle string `db:"piece_title"`
	composerID int    `db:"composer_id"`
	fullName   string `db:"full_name"`
	shortName  string `db:"short_name"`
	sequence   int    `db:"sequence"`
}

func (r *programmePieceRow) toProgrammePiece() content.ProgrammePiece {
	return content.ProgrammePiece{
		Piece: content.Piece{
			ID:    r.pieceID,
			Title: r.pieceTitle,
		},
		Composer: content.Composer{
			ID:        r.composerID,
			FullName:  r.fullName,
			ShortName: r.shortName,
		},
		Sequence: r.sequence,
	}
}

func (s *ProgrammePieceStore) ListByProgrammeID(
	ctx context.Context,
	id int,
) ([]content.ProgrammePiece, error) {
	query := `
	SELECT
	p.piece_id,
	p.piece_title,
	c.composer_id,
	c.full_name,
	c.short_name,
	pp.sequence
	FROM programme_pieces pp
	JOIN pieces p ON p.piece_id = pp.piece_id
	JOIN composers c ON c.composer_id = p.composer_id
	WHERE pp.programme_id = $1
	ORDER BY pp.sequence ASC
	`

	pgxRows, err := s.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	rows, err := collectRows[programmePieceRow](pgxRows)
	if err != nil {
		return nil, err
	}

	programmePieces := make([]content.ProgrammePiece, len(rows))
	for i, row := range rows {
		programmePieces[i] = row.toProgrammePiece()
	}

	return programmePieces, nil
}

func (s *ProgrammePieceStore) Update(
	ctx context.Context,
	id int,
	ids []int,
) ([]content.ProgrammePiece, error) {
	deleteQuery := `
	DELETE
	FROM programme_pieces
	WHERE programme_id = $1
	`
	_, err := s.db.Exec(ctx, deleteQuery, id)
	if err != nil {
		return nil, fmt.Errorf("delete query failed: %w", err)
	}

	if len(ids) == 0 {
		return []content.ProgrammePiece{}, nil
	}

	sequences := make([]int, len(ids))
	for i := range ids {
		sequences[i] = i + 1
	}

	// UNNEST here helps me turn arrays (or columns) of data into rows.
	// It's especially helpful here because we don't know for sure how many
	// pieces we want to insert. I prefer this approach over dynamically building
	// a query. Lastly, `::int[]` is a type annotation to help UNNEST
	// typecast the variable we insert there.
	insertQuery := `
	INSERT INTO programme_pieces (
		programme_id,
		piece_id,
		sequence
	)
	SELECT
		$1,
		piece_id,
		sequence
	FROM UNNEST($2::int[], $3::int[]) AS t(piece_id, sequence)
	`

	_, err = s.db.Exec(ctx, insertQuery,
		id,
		ids,
		sequences,
	)
	if err != nil {
		return nil, fmt.Errorf("insert query failed: %w", err)
	}

	selectQuery := `
	SELECT
		p.piece_id,
		p.piece_title,
		c.composer_id,
		c.full_name,
		c.short_name,
		pp.sequence
	FROM programme_pieces pp
	JOIN pieces p ON p.piece_id = pp.piece_id
	JOIN composers c ON c.composer_id = p.composer_id
	WHERE pp.programme_id = $1
	ORDER BY pp.sequence ASC
	`

	pgxRows, err := s.db.Query(ctx, selectQuery, id)
	if err != nil {
		return nil, fmt.Errorf("insert query failed: %w", err)
	}

	rows, err := collectRows[programmePieceRow](pgxRows)
	if err != nil {
		return nil, err
	}

	programmePieces := make([]content.ProgrammePiece, len(rows))
	for i, row := range rows {
		programmePieces[i] = row.toProgrammePiece()
	}

	return programmePieces, nil
}
