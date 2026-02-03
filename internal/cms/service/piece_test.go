package service

import (
	"context"
	"errors"
	"testing"

	"github.com/adamkadda/arman/internal/cms/model"
	"github.com/adamkadda/arman/internal/cms/store"
	"github.com/adamkadda/arman/internal/content"
	"github.com/stretchr/testify/require"
)

type mockPieceStore struct {
	piece          *content.Piece
	detailedPieces []model.PieceWithDetails
	detailedPiece  *model.PieceWithDetails
	err            error
	getErr         error
	deleteErr      error
}

func (s mockPieceStore) Get(
	ctx context.Context,
	id int,
) (*content.Piece, error) {
	return s.piece, s.err
}

func (s mockPieceStore) GetWithDetails(
	ctx context.Context,
	id int,
) (*model.PieceWithDetails, error) {
	return s.detailedPiece, s.getErr
}

func (s mockPieceStore) ListWithDetails(
	ctx context.Context,
) ([]model.PieceWithDetails, error) {
	return s.detailedPieces, s.err
}

func (s mockPieceStore) Create(
	ctx context.Context,
	v content.Piece,
) (*content.Piece, error) {
	return s.piece, s.err
}

func (s mockPieceStore) Update(
	ctx context.Context,
	v content.Piece,
) (*content.Piece, error) {
	return s.piece, s.err
}

func (s mockPieceStore) Delete(
	ctx context.Context,
	id int,
) error {
	return s.deleteErr
}

func TestPieceService_Get(t *testing.T) {
	tests := []struct {
		name    string
		piece   *content.Piece
		err     error
		wantErr bool
	}{
		{"piece.get success", &content.Piece{Title: "foo"}, nil, false},
		{"piece.get error", nil, errors.New("oops"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := PieceService{
				newPieceStore: func(db store.Executor) PieceStore {
					return mockPieceStore{
						piece: tt.piece,
						err:   tt.err,
					}
				},
			}

			piece, err := svc.Get(testContext(), 1)

			if tt.wantErr {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, piece)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.piece, piece)
			}
		})
	}
}

func TestPieceService_List(t *testing.T) {
	tests := []struct {
		name    string
		pieces  []model.PieceWithDetails
		err     error
		wantErr bool
	}{
		{"piece.list success", []model.PieceWithDetails{
			{
				Piece:          content.Piece{Title: "foo"},
				ProgrammeCount: 0,
			},
			{
				Piece:          content.Piece{Title: "bar"},
				ProgrammeCount: 1,
			},
			{
				Piece:          content.Piece{Title: "baz"},
				ProgrammeCount: 2,
			},
		}, nil, false},
		{"piece.list error", nil, errors.New("oops"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := PieceService{
				newPieceStore: func(db store.Executor) PieceStore {
					return mockPieceStore{
						detailedPieces: tt.pieces,
						err:            tt.err,
					}
				},
			}

			pieces, err := svc.List(testContext())

			if tt.wantErr {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, pieces)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.pieces, pieces)
			}
		})
	}
}

func TestPieceService_Create(t *testing.T) {
	tests := []struct {
		name             string
		cmd              model.UpsertPieceCommand
		piece            *content.Piece
		beginErr         error
		commitErr        error
		pieceStoreErr    error
		composerStoreErr error
		expectedErr      error
	}{
		{
			"begin tx failed",
			model.UpsertPieceCommand{},
			nil,
			ErrTxBegin,
			nil,
			nil,
			nil,
			ErrTxBegin,
		},
		{
			"operation mismatch",
			model.UpsertPieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationUpdate,
				},
				Composer: model.ComposerIntent{},
			},
			nil,
			nil,
			nil,
			nil,
			nil,
			content.ErrOperationMismatch,
		},
		{
			"invalid piece",
			model.UpsertPieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationCreate,
					Data: content.Piece{
						ComposerID: 1,
					},
				},
				Composer: model.ComposerIntent{},
			},
			nil,
			nil,
			nil,
			nil,
			nil,
			content.ErrInvalidResource,
		},
		{
			"composer resolver err",
			model.UpsertPieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationCreate,
					Data: content.Piece{
						Title:      "Foo",
						ComposerID: 1,
					},
				},
				Composer: model.ComposerIntent{
					Operation: model.OperationSelect,
					Data: content.Composer{
						ID: 1,
					},
				},
			},
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		},
		{
			"store error",
			model.UpsertPieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationCreate,
					Data: content.Piece{
						Title:      "Foo",
						ComposerID: 1,
					},
				},
				Composer: model.ComposerIntent{
					Operation: model.OperationSelect,
					Data: content.Composer{
						ID: 1,
					},
				},
			},
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		},
		{
			"commit tx failed",
			model.UpsertPieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationCreate,
					Data: content.Piece{
						Title:      "Foo",
						ComposerID: 1,
					},
				},
				Composer: model.ComposerIntent{
					Operation: model.OperationSelect,
					Data: content.Composer{
						ID: 1,
					},
				},
			},
			nil,
			nil,
			ErrTxCommit,
			nil,
			nil,
			ErrTxCommit,
		},
		{
			"success",
			model.UpsertPieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationCreate,
					Data: content.Piece{
						Title:      "Foo",
						ComposerID: 1,
					},
				},
				Composer: model.ComposerIntent{
					Operation: model.OperationSelect,
					Data: content.Composer{
						ID: 1,
					},
				},
			},
			&content.Piece{
				ID:         1,
				Title:      "Foo",
				ComposerID: 1,
			},
			nil,
			nil,
			nil,
			nil,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := PieceService{
				db: mockDB{
					tx: mockTx{
						err: tt.commitErr,
					},
					err: tt.beginErr,
				},
				newPieceStore: func(db store.Executor) PieceStore {
					return mockPieceStore{
						piece: tt.piece,
						err:   tt.pieceStoreErr,
					}
				},
				newComposerStore: func(db store.Executor) ComposerStore {
					return mockComposerStore{
						composer: &tt.cmd.Composer.Data,
						err:      tt.composerStoreErr,
					}
				},
			}

			piece, err := svc.Create(testContext(), tt.cmd)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.piece, piece)
			}
		})
	}
}

func TestPieceService_Update(t *testing.T) {
	tests := []struct {
		name    string
		cmd     model.UpsertPieceCommand
		piece   *content.Piece
		err     error
		wantErr bool
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := PieceService{
				newPieceStore: func(db store.Executor) PieceStore {
					return mockPieceStore{
						piece: tt.piece,
						err:   tt.err,
					}
				},
			}

			piece, err := svc.Update(testContext(), tt.cmd)

			if tt.wantErr {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, piece)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.piece, piece)
			}
		})
	}
}

func TestPieceService_Delete(t *testing.T) {
	tests := []struct {
		name          string
		piece         *model.PieceWithDetails
		getErr        error
		deleteErr     error
		expectedError error
	}{
		{
			name: "piece.delete success",
			piece: &model.PieceWithDetails{
				Piece:          content.Piece{Title: "foo"},
				ProgrammeCount: 0,
			},
			getErr:        nil,
			deleteErr:     nil,
			expectedError: nil,
		},
		{
			name:          "piece.get_with_details error",
			piece:         nil,
			getErr:        ErrGet,
			deleteErr:     nil,
			expectedError: ErrGet,
		},
		{
			name: "piece.delete blocked",
			piece: &model.PieceWithDetails{
				Piece:          content.Piece{Title: "foo"},
				ProgrammeCount: 1,
			},
			getErr:        nil,
			deleteErr:     nil,
			expectedError: content.ErrPieceProtected,
		},
		{
			name: "piece.delete error",
			piece: &model.PieceWithDetails{
				Piece:          content.Piece{Title: "foo"},
				ProgrammeCount: 0,
			},
			getErr:        nil,
			deleteErr:     ErrDelete,
			expectedError: ErrDelete,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := PieceService{
				newPieceStore: func(db store.Executor) PieceStore {
					return mockPieceStore{
						detailedPiece: tt.piece,
						getErr:        tt.getErr,
						deleteErr:     tt.deleteErr,
					}
				},
			}

			err := svc.Delete(testContext(), 1)

			if tt.expectedError != nil {
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
