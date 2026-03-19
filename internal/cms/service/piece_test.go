package service

import (
	"cmp"
	"context"
	"slices"
	"testing"

	"github.com/adamkadda/arman/internal/cms/model"
	"github.com/adamkadda/arman/internal/cms/store"
	"github.com/adamkadda/arman/internal/content"
	"github.com/stretchr/testify/require"
)

func TestPieceService_Get(t *testing.T) {
	tests := []struct {
		name        string
		store       mockPieceStore
		expected    *content.Piece
		expectedErr error
	}{
		{
			name: "success",
			store: mockPieceStore{
				pieces: map[int]*content.Piece{
					1: {ID: 1, Title: "Foo Sonata", ComposerID: 1},
				},
			},
			expected:    &content.Piece{ID: 1, Title: "Foo Sonata", ComposerID: 1},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := PieceService{
				newPieceStore: func(db store.Executor) PieceStore {
					return &tt.store
				},
			}

			piece, err := svc.Get(testContext(), 1)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, piece)
			}
		})
	}
}

func TestPieceService_List(t *testing.T) {
	tests := []struct {
		name        string
		store       mockPieceStore
		expected    []model.PieceWithDetails
		expectedErr error
	}{
		{
			name: "success",
			store: mockPieceStore{
				detailedPieces: map[int]*model.PieceWithDetails{
					1: {Piece: content.Piece{ID: 1, Title: "Foo Sonata", ComposerID: 1}, ProgrammeCount: 0},
					2: {Piece: content.Piece{ID: 2, Title: "Bar Toccata", ComposerID: 2}, ProgrammeCount: 1},
					3: {Piece: content.Piece{ID: 3, Title: "Baz Prelude", ComposerID: 3}, ProgrammeCount: 2},
				},
			},
			expected: []model.PieceWithDetails{
				{Piece: content.Piece{ID: 1, Title: "Foo Sonata", ComposerID: 1}, ProgrammeCount: 0},
				{Piece: content.Piece{ID: 2, Title: "Bar Toccata", ComposerID: 2}, ProgrammeCount: 1},
				{Piece: content.Piece{ID: 3, Title: "Baz Prelude", ComposerID: 3}, ProgrammeCount: 2},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := PieceService{
				newPieceStore: func(db store.Executor) PieceStore {
					return &tt.store
				},
			}

			pieces, err := svc.List(testContext())

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, pieces)
			}
		})
	}
}

func TestPieceService_Create(t *testing.T) {
	tests := []struct {
		name          string
		cmd           model.PieceCommand
		pieceStore    mockPieceStore
		composerStore mockComposerStore
		expected      *content.Piece
		expectedErr   error
	}{
		{
			name: "operation mismatch",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationUpdate,
					Data:      content.Piece{Title: "Foo Sonata", ComposerID: 1},
				},
				Composer: model.ComposerIntent{
					Operation: model.OperationSelect,
					Data:      content.Composer{ID: 1},
				},
			},
			pieceStore:    mockPieceStore{},
			composerStore: mockComposerStore{},
			expected:      nil,
			expectedErr:   content.ErrOperationMismatch,
		},
		{
			name: "piece validation failed",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationCreate,
					Data:      content.Piece{ComposerID: 1},
				},
				Composer: model.ComposerIntent{
					Operation: model.OperationSelect,
					Data:      content.Composer{ID: 1},
				},
			},
			pieceStore:    mockPieceStore{},
			composerStore: mockComposerStore{},
			expected:      nil,
			expectedErr:   content.ErrInvalidResource,
		},
		{
			name: "success",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationCreate,
					Data:      content.Piece{Title: "Foo Sonata", ComposerID: 1},
				},
				Composer: model.ComposerIntent{
					Operation: model.OperationSelect,
					Data:      content.Composer{ID: 1},
				},
			},
			pieceStore: mockPieceStore{
				createdPieces: []*content.Piece{
					{ID: 1, Title: "Foo Sonata", ComposerID: 1},
				},
			},
			composerStore: mockComposerStore{
				composers: map[int]*content.Composer{1: {ID: 1}},
			},
			expected:    &content.Piece{ID: 1, Title: "Foo Sonata", ComposerID: 1},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := PieceService{
				db: mockDB{},
				newPieceStore: func(db store.Executor) PieceStore {
					return &tt.pieceStore
				},
				newComposerStore: func(db store.Executor) ComposerStore {
					return &tt.composerStore
				},
			}

			piece, err := svc.Create(testContext(), tt.cmd)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, piece)
			}
		})
	}
}

func TestPieceService_Update(t *testing.T) {
	tests := []struct {
		name          string
		cmd           model.PieceCommand
		pieceStore    mockPieceStore
		composerStore mockComposerStore
		expected      *content.Piece
		expectedErr   error
	}{
		{
			name: "operation mismatch",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationCreate,
					Data:      content.Piece{ID: 1, Title: "Foo Sonata", ComposerID: 1},
				},
				Composer: model.ComposerIntent{
					Operation: model.OperationSelect,
					Data:      content.Composer{ID: 1},
				},
			},
			pieceStore:    mockPieceStore{},
			composerStore: mockComposerStore{},
			expected:      nil,
			expectedErr:   content.ErrOperationMismatch,
		},
		{
			name: "piece validation failed",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationUpdate,
					Data:      content.Piece{ID: 1},
				},
				Composer: model.ComposerIntent{
					Operation: model.OperationSelect,
					Data:      content.Composer{ID: 1},
				},
			},
			pieceStore:    mockPieceStore{},
			composerStore: mockComposerStore{},
			expected:      nil,
			expectedErr:   content.ErrInvalidResource,
		},
		{
			name: "success",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationUpdate,
					Data:      content.Piece{ID: 1, Title: "Foo Sonata", ComposerID: 1},
				},
				Composer: model.ComposerIntent{
					Operation: model.OperationSelect,
					Data:      content.Composer{ID: 1},
				},
			},
			pieceStore: mockPieceStore{
				pieces: map[int]*content.Piece{
					1: {ID: 1, Title: "Foo Sonata", ComposerID: 1},
				},
			},
			composerStore: mockComposerStore{
				composers: map[int]*content.Composer{1: {ID: 1}},
			},
			expected:    &content.Piece{ID: 1, Title: "Foo Sonata", ComposerID: 1},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := PieceService{
				db: mockDB{},
				newPieceStore: func(db store.Executor) PieceStore {
					return &tt.pieceStore
				},
				newComposerStore: func(db store.Executor) ComposerStore {
					return &tt.composerStore
				},
			}

			piece, err := svc.Update(testContext(), tt.cmd)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, piece)
			}
		})
	}
}

func TestPieceService_Delete(t *testing.T) {
	tests := []struct {
		name        string
		store       mockPieceStore
		expectedErr error
	}{
		{
			name: "piece protected",
			store: mockPieceStore{
				detailedPieces: map[int]*model.PieceWithDetails{
					1: {Piece: content.Piece{ID: 1, Title: "foo"}, ProgrammeCount: 1},
				},
			},
			expectedErr: content.ErrPieceProtected,
		},
		{
			name: "success",
			store: mockPieceStore{
				detailedPieces: map[int]*model.PieceWithDetails{
					1: {Piece: content.Piece{ID: 1, Title: "foo"}, ProgrammeCount: 0},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := PieceService{
				newPieceStore: func(db store.Executor) PieceStore {
					return &tt.store
				},
			}

			err := svc.Delete(testContext(), 1)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPieceResolver_Run(t *testing.T) {
	tests := []struct {
		name        string
		intent      model.PieceIntent
		store       mockPieceStore
		expectedErr error
	}{
		{
			name: "invalid operation",
			intent: model.PieceIntent{
				Operation: model.Operation("DELETE"),
				Data: content.Piece{
					ID:         1,
					Title:      "Foo Sonata",
					ComposerID: 1,
				},
			},
			store: mockPieceStore{
				pieces: map[int]*content.Piece{
					1: {ID: 1, Title: "Foo Sonata", ComposerID: 1},
				},
			},
			expectedErr: model.ErrInvalidOperation,
		},
		{
			name: "select success",
			intent: model.PieceIntent{
				Operation: model.OperationSelect,
				Data: content.Piece{
					ID:         1,
					Title:      "Foo Sonata",
					ComposerID: 1,
				},
			},
			store: mockPieceStore{
				pieces: map[int]*content.Piece{
					1: {ID: 1, Title: "Foo Sonata", ComposerID: 1},
				},
			},
			expectedErr: nil,
		},
		{
			name: "create success",
			intent: model.PieceIntent{
				Operation: model.OperationCreate,
				Data: content.Piece{
					Title:      "Foo Sonata",
					ComposerID: 1,
				},
			},
			store: mockPieceStore{
				createdPieces: []*content.Piece{
					{ID: 1, Title: "Foo Sonata", ComposerID: 1},
				},
			},
			expectedErr: nil,
		},
		{
			name: "update success",
			intent: model.PieceIntent{
				Operation: model.OperationUpdate,
				Data: content.Piece{
					ID:         1,
					Title:      "Foo Sonata",
					ComposerID: 1,
				},
			},
			store: mockPieceStore{
				pieces: map[int]*content.Piece{
					1: {ID: 1, Title: "Foo Sonata", ComposerID: 1},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resolver := newPieceResolver(&tt.store)

			_, err := resolver.run(testContext(), tt.intent)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

type mockPieceStore struct {
	pieces         map[int]*content.Piece
	detailedPieces map[int]*model.PieceWithDetails
	createdPieces  []*content.Piece
	err            error
}

func (s *mockPieceStore) Get(
	ctx context.Context,
	id int,
) (*content.Piece, error) {
	return s.pieces[id], s.err
}

func (s *mockPieceStore) GetWithDetails(
	ctx context.Context,
	id int,
) (*model.PieceWithDetails, error) {
	return s.detailedPieces[id], s.err
}

func (s *mockPieceStore) ListWithDetails(
	ctx context.Context,
) ([]model.PieceWithDetails, error) {
	if s.err != nil {
		return nil, s.err
	}

	pieces := make([]model.PieceWithDetails, 0, len(s.detailedPieces))
	for _, p := range s.detailedPieces {
		pieces = append(pieces, *p)
	}

	slices.SortFunc(pieces, func(a, b model.PieceWithDetails) int {
		return cmp.Compare(a.Piece.ID, b.Piece.ID)
	})

	return pieces, nil
}

func (s *mockPieceStore) Create(
	ctx context.Context,
	v content.Piece,
) (*content.Piece, error) {
	if s.err != nil {
		return nil, s.err
	}
	p := s.createdPieces[0]
	s.createdPieces = s.createdPieces[1:]
	return p, nil
}

func (s *mockPieceStore) Update(
	ctx context.Context,
	v content.Piece,
) (*content.Piece, error) {
	return s.pieces[v.ID], s.err
}

func (s *mockPieceStore) Delete(
	ctx context.Context,
	id int,
) error {
	return s.err
}
