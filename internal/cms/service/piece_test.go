package service

import (
	"context"
	"testing"

	"github.com/adamkadda/arman/internal/cms/model"
	"github.com/adamkadda/arman/internal/cms/store"
	"github.com/adamkadda/arman/internal/content"
	"github.com/stretchr/testify/require"
)

func TestPieceService_Get(t *testing.T) {
	tests := []struct {
		name          string
		expectedPiece *content.Piece
		expectedErr   error
	}{
		{
			name:          "store error",
			expectedPiece: nil,
			expectedErr:   ErrGet,
		},
		{
			name: "success",
			expectedPiece: &content.Piece{
				ID:         1,
				Title:      "Foo Sonata",
				ComposerID: 1,
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := PieceService{
				newPieceStore: func(db store.Executor) PieceStore {
					return mockPieceStore{
						piece: tt.expectedPiece,
						err:   tt.expectedErr,
					}
				},
			}

			piece, err := svc.Get(testContext(), 1)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedPiece, piece)
			}
		})
	}
}

func TestPieceService_List(t *testing.T) {
	tests := []struct {
		name           string
		expectedPieces []model.PieceWithDetails
		expectedErr    error
	}{
		{
			name:           "store error",
			expectedPieces: nil,
			expectedErr:    ErrFoo,
		},
		{
			name: "success",
			expectedPieces: []model.PieceWithDetails{
				{
					Piece: content.Piece{
						ID:         1,
						Title:      "Foo Sonata",
						ComposerID: 1,
					},
					ProgrammeCount: 0,
				},
				{
					Piece: content.Piece{
						ID:         2,
						Title:      "Bar Toccata",
						ComposerID: 2,
					},
					ProgrammeCount: 1,
				},
				{
					Piece: content.Piece{
						ID:         3,
						Title:      "Baz Prelude",
						ComposerID: 3,
					},
					ProgrammeCount: 2,
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
					return mockPieceStore{
						detailedPieces: tt.expectedPieces,
						err:            tt.expectedErr,
					}
				},
			}

			pieces, err := svc.List(testContext())

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedPieces, pieces)
			}
		})
	}
}

func TestPieceService_Create(t *testing.T) {
	tests := []struct {
		name             string
		cmd              model.PieceCommand
		expectedPiece    *content.Piece
		beginErr         error
		commitErr        error
		pieceStoreErr    error
		composerStoreErr error
		expectedErr      error
	}{
		{
			name:             "tx begin failed",
			cmd:              model.PieceCommand{},
			expectedPiece:    nil,
			beginErr:         ErrTxBegin,
			commitErr:        nil,
			pieceStoreErr:    nil,
			composerStoreErr: nil,
			expectedErr:      ErrTxBegin,
		},
		{
			name: "operation mismatch",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationUpdate,
					Data: content.Piece{
						Title:      "Foo Sonata",
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
			expectedPiece:    nil,
			beginErr:         nil,
			commitErr:        nil,
			pieceStoreErr:    nil,
			composerStoreErr: nil,
			expectedErr:      content.ErrOperationMismatch,
		},
		{
			name: "piece validation failed",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationCreate,
					Data: content.Piece{
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
			expectedPiece:    nil,
			beginErr:         nil,
			commitErr:        nil,
			pieceStoreErr:    nil,
			composerStoreErr: nil,
			expectedErr:      content.ErrInvalidResource,
		},
		{
			name: "composer resolver error",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationCreate,
					Data: content.Piece{
						Title:      "Foo Sonata",
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
			expectedPiece:    nil,
			beginErr:         nil,
			commitErr:        nil,
			pieceStoreErr:    nil,
			composerStoreErr: ErrGet,
			expectedErr:      ErrGet,
		},
		{
			name: "piece store error",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationCreate,
					Data: content.Piece{
						Title:      "Foo Sonata",
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
			expectedPiece:    nil,
			beginErr:         nil,
			commitErr:        nil,
			pieceStoreErr:    ErrFoo,
			composerStoreErr: nil,
			expectedErr:      ErrFoo,
		},
		{
			name: "commit tx failed",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationCreate,
					Data: content.Piece{
						Title:      "Foo Sonata",
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
			expectedPiece:    nil,
			beginErr:         nil,
			commitErr:        ErrTxCommit,
			pieceStoreErr:    nil,
			composerStoreErr: nil,
			expectedErr:      ErrTxCommit,
		},
		{
			name: "success",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationCreate,
					Data: content.Piece{
						Title:      "Foo Sonata",
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
			expectedPiece: &content.Piece{
				ID:         1,
				Title:      "Foo Sonata",
				ComposerID: 1,
			},
			beginErr:         nil,
			commitErr:        nil,
			pieceStoreErr:    nil,
			composerStoreErr: nil,
			expectedErr:      nil,
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
						piece: tt.expectedPiece,
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
				require.Equal(t, tt.expectedPiece, piece)
			}
		})
	}
}

func TestPieceService_Update(t *testing.T) {
	tests := []struct {
		name             string
		cmd              model.PieceCommand
		expectedPiece    *content.Piece
		beginErr         error
		commitErr        error
		pieceStoreErr    error
		composerStoreErr error
		expectedErr      error
	}{
		{
			name:             "tx begin failed",
			cmd:              model.PieceCommand{},
			expectedPiece:    nil,
			beginErr:         ErrTxBegin,
			commitErr:        nil,
			pieceStoreErr:    nil,
			composerStoreErr: nil,
			expectedErr:      ErrTxBegin,
		},
		{
			name: "operation mismatch",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationCreate,
					Data: content.Piece{
						ID:         1,
						Title:      "Foo Sonata",
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
			expectedPiece:    nil,
			beginErr:         nil,
			commitErr:        nil,
			pieceStoreErr:    nil,
			composerStoreErr: nil,
			expectedErr:      content.ErrOperationMismatch,
		},
		{
			name: "piece validation failed",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationUpdate,
					Data: content.Piece{
						ID: 1,
					},
				},
				Composer: model.ComposerIntent{
					Operation: model.OperationSelect,
					Data: content.Composer{
						ID: 1,
					},
				},
			},
			expectedPiece:    nil,
			beginErr:         nil,
			commitErr:        nil,
			pieceStoreErr:    nil,
			composerStoreErr: nil,
			expectedErr:      content.ErrInvalidResource,
		},
		{
			name: "composer resolver error",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationUpdate,
					Data: content.Piece{
						ID:         1,
						Title:      "Foo Sonata",
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
			expectedPiece:    nil,
			beginErr:         nil,
			commitErr:        nil,
			pieceStoreErr:    nil,
			composerStoreErr: ErrGet,
			expectedErr:      ErrGet,
		},
		{
			name: "piece store error",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationUpdate,
					Data: content.Piece{
						ID:         1,
						Title:      "Foo Sonata",
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
			expectedPiece:    nil,
			beginErr:         nil,
			commitErr:        nil,
			pieceStoreErr:    ErrFoo,
			composerStoreErr: nil,
			expectedErr:      ErrFoo,
		},
		{
			name: "tx commit failed",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationUpdate,
					Data: content.Piece{
						ID:         1,
						Title:      "Foo Sonata",
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
			expectedPiece:    nil,
			beginErr:         nil,
			commitErr:        ErrTxCommit,
			pieceStoreErr:    nil,
			composerStoreErr: nil,
			expectedErr:      ErrTxCommit,
		},
		{
			name: "success",
			cmd: model.PieceCommand{
				Piece: model.PieceIntent{
					Operation: model.OperationUpdate,
					Data: content.Piece{
						ID:         1,
						Title:      "Foo Sonata",
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
			expectedPiece: &content.Piece{
				ID:         1,
				Title:      "Foo Sonata",
				ComposerID: 1,
			},
			beginErr:         nil,
			commitErr:        nil,
			pieceStoreErr:    nil,
			composerStoreErr: nil,
			expectedErr:      nil,
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
						piece: &tt.cmd.Piece.Data,
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

			piece, err := svc.Update(testContext(), tt.cmd)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedPiece, piece)
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
			name:          "get error",
			piece:         nil,
			getErr:        ErrGet,
			deleteErr:     nil,
			expectedError: ErrGet,
		},
		{
			name: "piece protected",
			piece: &model.PieceWithDetails{
				Piece:          content.Piece{Title: "foo"},
				ProgrammeCount: 1,
			},
			getErr:        nil,
			deleteErr:     nil,
			expectedError: content.ErrPieceProtected,
		},
		{
			name: "delete error",
			piece: &model.PieceWithDetails{
				Piece:          content.Piece{Title: "foo"},
				ProgrammeCount: 0,
			},
			getErr:        nil,
			deleteErr:     ErrDelete,
			expectedError: ErrDelete,
		},
		{
			name: "success",
			piece: &model.PieceWithDetails{
				Piece:          content.Piece{Title: "foo"},
				ProgrammeCount: 0,
			},
			getErr:        nil,
			deleteErr:     nil,
			expectedError: nil,
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

func TestPieceResolver_Run(t *testing.T) {
	tests := []struct {
		name        string
		intent      model.PieceIntent
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
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resolver := newPieceResolver(mockPieceStore{})

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
