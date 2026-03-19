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

func TestProgrammeService_Get(t *testing.T) {
	tests := []struct {
		name           string
		programmeStore mockProgrammeStore
		expected       *model.ProgrammeWithPieces
		expectedErr    error
	}{
		{
			name: "success",
			programmeStore: mockProgrammeStore{
				programmes: map[int]*content.Programme{
					1: {ID: 1, Title: "Foo Programme"},
				},
				pieces: map[int]content.ProgrammePiece{
					1: {
						Piece:    content.Piece{ID: 1, Title: "Foo Sonata", ComposerID: 1},
						Composer: content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
						Sequence: 1,
					},
					2: {
						Piece:    content.Piece{ID: 2, Title: "Bar Toccata", ComposerID: 1},
						Composer: content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
						Sequence: 2,
					},
				},
			},
			expected: &model.ProgrammeWithPieces{
				Programme: &content.Programme{ID: 1, Title: "Foo Programme"},
				Pieces: []content.ProgrammePiece{
					{
						Piece:    content.Piece{ID: 1, Title: "Foo Sonata", ComposerID: 1},
						Composer: content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
						Sequence: 1,
					},
					{
						Piece:    content.Piece{ID: 2, Title: "Bar Toccata", ComposerID: 1},
						Composer: content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
						Sequence: 2,
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ProgrammeService{
				newProgrammeStore: func(db store.Executor) ProgrammeStore {
					return &tt.programmeStore
				},
			}

			programme, err := svc.Get(testContext(), 1)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, programme)
			}
		})
	}
}

func TestProgrammeService_List(t *testing.T) {
	tests := []struct {
		name           string
		programmeStore mockProgrammeStore
		expected       []model.ProgrammeWithDetails
		expectedErr    error
	}{
		{
			name: "success",
			programmeStore: mockProgrammeStore{
				detailedProgrammes: map[int]*model.ProgrammeWithDetails{
					1: {Programme: content.Programme{ID: 1, Title: "Foo Programme"}, PieceCount: 0, EventCount: 0},
					2: {Programme: content.Programme{ID: 2, Title: "Bar Programme"}, PieceCount: 1, EventCount: 0},
					3: {Programme: content.Programme{ID: 3, Title: "Baz Programme"}, PieceCount: 2, EventCount: 1},
				},
			},
			expected: []model.ProgrammeWithDetails{
				{Programme: content.Programme{ID: 1, Title: "Foo Programme"}, PieceCount: 0, EventCount: 0},
				{Programme: content.Programme{ID: 2, Title: "Bar Programme"}, PieceCount: 1, EventCount: 0},
				{Programme: content.Programme{ID: 3, Title: "Baz Programme"}, PieceCount: 2, EventCount: 1},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Parallel()

		svc := ProgrammeService{
			newProgrammeStore: func(db store.Executor) ProgrammeStore {
				return &tt.programmeStore
			},
		}

		programmes, err := svc.List(testContext())

		if tt.expectedErr != nil {
			require.ErrorIs(t, err, tt.expectedErr)
		} else {
			require.NoError(t, err)
			require.Equal(t, tt.expected, programmes)
		}
	}
}

func TestProgrammeService_Create(t *testing.T) {
	tests := []struct {
		name           string
		cmd            model.ProgrammeCommand
		programmeStore ProgrammeStore
		pieceStore     PieceStore
		composerStore  ComposerStore
		expected       *model.ProgrammeWithPieces
		expectedErr    error
	}{}

	for _, tt := range tests {
		t.Parallel()

		svc := ProgrammeService{
			db: mockDB{
				tx: mockTx{
					err: nil,
				},
				err: nil,
			},
			newProgrammeStore: func(db store.Executor) ProgrammeStore {
				return tt.programmeStore
			},
			newPieceStore: func(db store.Executor) PieceStore {
				return tt.pieceStore
			},
			newComposerStore: func(db store.Executor) ComposerStore {
				return tt.composerStore
			},
		}

		programme, err := svc.Create(testContext(), tt.cmd)

		if tt.expectedErr != nil {
			require.ErrorIs(t, err, tt.expectedErr)
		} else {
			require.NoError(t, err)
			require.Equal(t, tt.expected, programme)
		}
	}
}

func TestProgrammeService_Update(t *testing.T) {
	tempID1 := new(int)
	*tempID1 = 1

	tests := []struct {
		name           string
		cmd            model.ProgrammeCommand
		programmeStore mockProgrammeStore
		pieceStore     mockPieceStore
		composerStore  mockComposerStore
		expected       *model.ProgrammeWithPieces
		expectedErr    error
	}{
		{
			name: "operation mismatch",
			cmd: model.ProgrammeCommand{
				Programme: model.ProgrammeIntent{
					Operation: model.OperationCreate,
					Data:      content.Programme{ID: 1, Title: "Foo Programme"},
				},
			},
			programmeStore: mockProgrammeStore{},
			pieceStore:     mockPieceStore{},
			composerStore:  mockComposerStore{},
			expected:       nil,
			expectedErr:    content.ErrOperationMismatch,
		},
		{
			name: "programme validation failed",
			cmd: model.ProgrammeCommand{
				Programme: model.ProgrammeIntent{
					Operation: model.OperationUpdate,
					Data:      content.Programme{ID: 1},
				},
			},
			programmeStore: mockProgrammeStore{},
			pieceStore:     mockPieceStore{},
			composerStore:  mockComposerStore{},
			expected:       nil,
			expectedErr:    content.ErrInvalidResource,
		},
		{
			name: "no pieces",
			cmd: model.ProgrammeCommand{
				Programme: model.ProgrammeIntent{
					Operation: model.OperationUpdate,
					Data:      content.Programme{ID: 1, Title: "Foo Programme"},
				},
				Pieces: []model.PieceCommand{},
			},
			programmeStore: mockProgrammeStore{
				programmes: map[int]*content.Programme{
					1: {ID: 1, Title: "Foo Programme"},
				},
				pieces: map[int]content.ProgrammePiece{},
			},
			pieceStore:    mockPieceStore{},
			composerStore: mockComposerStore{},
			expected: &model.ProgrammeWithPieces{
				Programme: &content.Programme{ID: 1, Title: "Foo Programme"},
				Pieces:    []content.ProgrammePiece{},
			},
			expectedErr: nil,
		},
		{
			name: "single piece: select composer, select piece",
			cmd: model.ProgrammeCommand{
				Programme: model.ProgrammeIntent{
					Operation: model.OperationUpdate,
					Data:      content.Programme{ID: 1, Title: "Foo Programme"},
				},
				Pieces: []model.PieceCommand{
					{
						Piece: model.PieceIntent{
							Operation: model.OperationSelect,
							Data:      content.Piece{ID: 1},
						},
						Composer: model.ComposerIntent{
							Operation: model.OperationSelect,
							Data:      content.Composer{ID: 1},
						},
					},
				},
			},
			programmeStore: mockProgrammeStore{
				programmes: map[int]*content.Programme{
					1: {ID: 1, Title: "Foo Programme"},
				},
				pieces: map[int]content.ProgrammePiece{
					1: {
						Piece:    content.Piece{ID: 1, Title: "Foo Sonata", ComposerID: 1},
						Composer: content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
						Sequence: 1,
					},
				},
			},
			pieceStore: mockPieceStore{
				pieces: map[int]*content.Piece{
					1: {ID: 1, Title: "Foo Sonata", ComposerID: 1},
				},
			},
			composerStore: mockComposerStore{
				composers: map[int]*content.Composer{
					1: {ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
				},
			},
			expected: &model.ProgrammeWithPieces{
				Programme: &content.Programme{ID: 1, Title: "Foo Programme"},
				Pieces: []content.ProgrammePiece{
					{
						Piece:    content.Piece{ID: 1, Title: "Foo Sonata", ComposerID: 1},
						Composer: content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
						Sequence: 1,
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "single piece: create composer, create piece",
			cmd: model.ProgrammeCommand{
				Programme: model.ProgrammeIntent{
					Operation: model.OperationUpdate,
					Data:      content.Programme{ID: 1, Title: "Foo Programme"},
				},
				Pieces: []model.PieceCommand{
					{
						Piece: model.PieceIntent{
							Operation: model.OperationCreate,
							Data:      content.Piece{Title: "Foo Sonata"},
						},
						Composer: model.ComposerIntent{
							Operation: model.OperationCreate,
							TempID:    tempID1,
							Data:      content.Composer{FullName: "Foo Foolington", ShortName: "Foolington"},
						},
					},
				},
			},
			programmeStore: mockProgrammeStore{
				programmes: map[int]*content.Programme{
					1: {ID: 1, Title: "Foo Programme"},
				},
				pieces: map[int]content.ProgrammePiece{
					1: {
						Piece:    content.Piece{ID: 1, Title: "Foo Sonata", ComposerID: 1},
						Composer: content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
						Sequence: 1,
					},
				},
			},
			pieceStore: mockPieceStore{
				createdPieces: []*content.Piece{
					{ID: 1, Title: "Foo Sonata", ComposerID: 1},
				},
			},
			composerStore: mockComposerStore{
				createdComposers: []*content.Composer{
					{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
				},
			},
			expected: &model.ProgrammeWithPieces{
				Programme: &content.Programme{ID: 1, Title: "Foo Programme"},
				Pieces: []content.ProgrammePiece{
					{
						Piece:    content.Piece{ID: 1, Title: "Foo Sonata", ComposerID: 1},
						Composer: content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
						Sequence: 1,
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "multiple pieces: shared composer deduplication",
			cmd: model.ProgrammeCommand{
				Programme: model.ProgrammeIntent{
					Operation: model.OperationUpdate,
					Data:      content.Programme{ID: 1, Title: "Foo Programme"},
				},
				Pieces: []model.PieceCommand{
					{
						Piece: model.PieceIntent{
							Operation: model.OperationCreate,
							Data:      content.Piece{Title: "Foo Sonata"},
						},
						Composer: model.ComposerIntent{
							Operation: model.OperationCreate,
							TempID:    tempID1,
							Data:      content.Composer{FullName: "Foo Foolington", ShortName: "Foolington"},
						},
					},
					{
						Piece: model.PieceIntent{
							Operation: model.OperationCreate,
							Data:      content.Piece{Title: "Bar Toccata"},
						},
						Composer: model.ComposerIntent{
							Operation: model.OperationCreate,
							TempID:    tempID1,
							Data:      content.Composer{FullName: "Foo Foolington", ShortName: "Foolington"},
						},
					},
				},
			},
			programmeStore: mockProgrammeStore{
				programmes: map[int]*content.Programme{
					1: {ID: 1, Title: "Foo Programme"},
				},
				pieces: map[int]content.ProgrammePiece{
					1: {
						Piece:    content.Piece{ID: 1, Title: "Foo Sonata", ComposerID: 1},
						Composer: content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
						Sequence: 1,
					},
					2: {
						Piece:    content.Piece{ID: 2, Title: "Bar Toccata", ComposerID: 1},
						Composer: content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
						Sequence: 2,
					},
				},
			},
			pieceStore: mockPieceStore{
				createdPieces: []*content.Piece{
					{ID: 1, Title: "Foo Sonata", ComposerID: 1},
					{ID: 2, Title: "Bar Toccata", ComposerID: 1},
				},
			},
			composerStore: mockComposerStore{
				createdComposers: []*content.Composer{
					{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
				},
			},
			expected: &model.ProgrammeWithPieces{
				Programme: &content.Programme{ID: 1, Title: "Foo Programme"},
				Pieces: []content.ProgrammePiece{
					{
						Piece:    content.Piece{ID: 1, Title: "Foo Sonata", ComposerID: 1},
						Composer: content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
						Sequence: 1,
					},
					{
						Piece:    content.Piece{ID: 2, Title: "Bar Toccata", ComposerID: 1},
						Composer: content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
						Sequence: 2,
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "multiple pieces: mixed operations",
			cmd: model.ProgrammeCommand{
				Programme: model.ProgrammeIntent{
					Operation: model.OperationUpdate,
					Data:      content.Programme{ID: 1, Title: "Foo Programme"},
				},
				Pieces: []model.PieceCommand{
					{
						Piece: model.PieceIntent{
							Operation: model.OperationSelect,
							Data:      content.Piece{ID: 1},
						},
						Composer: model.ComposerIntent{
							Operation: model.OperationCreate,
							TempID:    tempID1,
							Data:      content.Composer{FullName: "Foo Foolington", ShortName: "Foolington"},
						},
					},
					{
						Piece: model.PieceIntent{
							Operation: model.OperationCreate,
							Data:      content.Piece{Title: "Bar Toccata"},
						},
						Composer: model.ComposerIntent{
							Operation: model.OperationSelect,
							Data:      content.Composer{ID: 2},
						},
					},
				},
			},
			programmeStore: mockProgrammeStore{
				programmes: map[int]*content.Programme{
					1: {ID: 1, Title: "Foo Programme"},
				},
				pieces: map[int]content.ProgrammePiece{
					1: {
						Piece:    content.Piece{ID: 1, Title: "Foo Sonata", ComposerID: 1},
						Composer: content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
						Sequence: 1,
					},
					2: {
						Piece:    content.Piece{ID: 2, Title: "Bar Toccata", ComposerID: 2},
						Composer: content.Composer{ID: 2, FullName: "Bar Bartholomew", ShortName: "Bartholomew"},
						Sequence: 2,
					},
				},
			},
			pieceStore: mockPieceStore{
				pieces: map[int]*content.Piece{
					1: {ID: 1, Title: "Foo Sonata", ComposerID: 1},
				},
				createdPieces: []*content.Piece{
					{ID: 2, Title: "Bar Toccata", ComposerID: 2},
				},
			},
			composerStore: mockComposerStore{
				createdComposers: []*content.Composer{
					{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
				},
				composers: map[int]*content.Composer{
					2: {ID: 2, FullName: "Bar Bartholomew", ShortName: "Bartholomew"},
				},
			},
			expected: &model.ProgrammeWithPieces{
				Programme: &content.Programme{ID: 1, Title: "Foo Programme"},
				Pieces: []content.ProgrammePiece{
					{
						Piece:    content.Piece{ID: 1, Title: "Foo Sonata", ComposerID: 1},
						Composer: content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
						Sequence: 1,
					},
					{
						Piece:    content.Piece{ID: 2, Title: "Bar Toccata", ComposerID: 2},
						Composer: content.Composer{ID: 2, FullName: "Bar Bartholomew", ShortName: "Bartholomew"},
						Sequence: 2,
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ProgrammeService{
				db: mockDB{},
				newProgrammeStore: func(db store.Executor) ProgrammeStore {
					return &tt.programmeStore
				},
				newPieceStore: func(db store.Executor) PieceStore {
					return &tt.pieceStore
				},
				newComposerStore: func(db store.Executor) ComposerStore {
					return &tt.composerStore
				},
			}

			programme, err := svc.Update(testContext(), tt.cmd)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, programme)
			}
		})
	}
}

func TestProgrammeService_Delete(t *testing.T) {
	tests := []struct {
		name           string
		programmeStore mockProgrammeStore
		expectedErr    error
	}{
		{
			name: "programme protected",
			programmeStore: mockProgrammeStore{
				detailedProgrammes: map[int]*model.ProgrammeWithDetails{
					1: {Programme: content.Programme{ID: 1, Title: "Foo Programme"}, EventCount: 1},
				},
			},
			expectedErr: content.ErrProgrammeProtected,
		},
		{
			name: "success",
			programmeStore: mockProgrammeStore{
				detailedProgrammes: map[int]*model.ProgrammeWithDetails{
					1: {Programme: content.Programme{ID: 1, Title: "Foo Programme"}, EventCount: 0},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ProgrammeService{
				newProgrammeStore: func(db store.Executor) ProgrammeStore {
					return &tt.programmeStore
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

func TestProgrammePieceResolver_Run(t *testing.T) {
	tempID1 := new(int)
	*tempID1 = 1

	tests := []struct {
		name           string
		cmd            model.PieceCommand
		programmeStore mockProgrammeStore
		pieceStore     mockPieceStore
		composerStore  mockComposerStore
		expectedErr    error
	}{
		{
			name: "select composer, select piece",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationSelect,
					Data:      content.Piece{ID: 1},
				},
				Composer: model.ComposerIntent{
					Operation: model.OperationSelect,
					Data:      content.Composer{ID: 1},
				},
			},
			programmeStore: mockProgrammeStore{},
			pieceStore: mockPieceStore{
				pieces: map[int]*content.Piece{
					1: {ID: 1, Title: "Foo Sonata", ComposerID: 1},
				},
			},
			composerStore: mockComposerStore{
				composers: map[int]*content.Composer{
					1: {ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
				},
			},
			expectedErr: nil,
		},
		{
			name: "create composer, create piece",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationCreate,
					Data:      content.Piece{Title: "Foo Sonata"},
				},
				Composer: model.ComposerIntent{
					Operation: model.OperationCreate,
					TempID:    tempID1,
					Data:      content.Composer{FullName: "Foo Foolington", ShortName: "Foolington"},
				},
			},
			programmeStore: mockProgrammeStore{},
			pieceStore: mockPieceStore{
				createdPieces: []*content.Piece{
					{ID: 1, Title: "Foo Sonata", ComposerID: 1},
				},
			},
			composerStore: mockComposerStore{
				createdComposers: []*content.Composer{
					{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
				},
			},
			expectedErr: nil,
		},
		{
			name: "create composer deduplication",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationCreate,
					Data:      content.Piece{Title: "Foo Sonata"},
				},
				Composer: model.ComposerIntent{
					Operation: model.OperationCreate,
					TempID:    tempID1,
					Data:      content.Composer{FullName: "Foo Foolington", ShortName: "Foolington"},
				},
			},
			programmeStore: mockProgrammeStore{},
			pieceStore: mockPieceStore{
				createdPieces: []*content.Piece{
					{ID: 1, Title: "Foo Sonata", ComposerID: 1},
				},
			},
			composerStore: mockComposerStore{
				createdComposers: []*content.Composer{
					{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resolver := newProgrammePieceResolver(
				&tt.programmeStore,
				newPieceResolver(&tt.pieceStore),
				newComposerResolver(&tt.composerStore),
			)

			_, err := resolver.run(testContext(), tt.cmd)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

type mockProgrammeStore struct {
	programmes         map[int]*content.Programme
	detailedProgrammes map[int]*model.ProgrammeWithDetails
	pieces             map[int]content.ProgrammePiece
	err                error
}

func (s *mockProgrammeStore) Get(
	ctx context.Context,
	id int,
) (*content.Programme, error) {
	return s.programmes[id], s.err
}

func (s *mockProgrammeStore) GetWithDetails(
	ctx context.Context,
	id int,
) (*model.ProgrammeWithDetails, error) {
	return s.detailedProgrammes[id], s.err
}

func (s *mockProgrammeStore) ListWithDetails(
	ctx context.Context,
) ([]model.ProgrammeWithDetails, error) {
	if s.err != nil {
		return nil, s.err
	}

	programmes := make([]model.ProgrammeWithDetails, 0, len(s.detailedProgrammes))
	for _, p := range s.detailedProgrammes {
		programmes = append(programmes, *p)
	}

	slices.SortFunc(programmes, func(a, b model.ProgrammeWithDetails) int {
		return cmp.Compare(a.Programme.ID, b.Programme.ID)
	})

	return programmes, nil
}

func (s *mockProgrammeStore) ListPieces(
	ctx context.Context,
	id int,
) ([]content.ProgrammePiece, error) {
	if s.err != nil {
		return nil, s.err
	}

	pieces := make([]content.ProgrammePiece, 0, len(s.pieces))
	for _, p := range s.pieces {
		pieces = append(pieces, p)
	}

	slices.SortFunc(pieces, func(a, b content.ProgrammePiece) int {
		return cmp.Compare(a.Sequence, b.Sequence)
	})

	return pieces, nil
}

func (s *mockProgrammeStore) UpdatePieces(
	ctx context.Context,
	id int,
	ids []int,
) ([]content.ProgrammePiece, error) {
	if s.err != nil {
		return nil, s.err
	}

	pieces := make([]content.ProgrammePiece, 0, len(s.pieces))
	for _, p := range s.pieces {
		pieces = append(pieces, p)
	}

	slices.SortFunc(pieces, func(a, b content.ProgrammePiece) int {
		return cmp.Compare(a.Sequence, b.Sequence)
	})

	return pieces, nil
}

func (s *mockProgrammeStore) Create(
	ctx context.Context,
	p content.Programme,
) (*content.Programme, error) {
	return s.programmes[0], s.err
}

func (s *mockProgrammeStore) Update(
	ctx context.Context,
	p content.Programme,
) (*content.Programme, error) {
	return s.programmes[p.ID], s.err
}

func (s *mockProgrammeStore) Delete(
	ctx context.Context,
	id int,
) error {
	return s.err
}
