package service

import (
	"context"
	"testing"

	"github.com/adamkadda/arman/internal/cms/model"
	"github.com/adamkadda/arman/internal/cms/store"
	"github.com/adamkadda/arman/internal/content"
	"github.com/stretchr/testify/require"
)

func TestComposerService_Get(t *testing.T) {
	tests := []struct {
		name             string
		expectedComposer *content.Composer
		expectedErr      error
	}{
		{
			name:             "store error",
			expectedComposer: nil,
			expectedErr:      ErrGet,
		},
		{
			name: "success",
			expectedComposer: &content.Composer{
				ID:        1,
				FullName:  "Foo Foolington",
				ShortName: "Foolington",
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ComposerService{
				newComposerStore: func(db store.Executor) ComposerStore {
					return mockComposerStore{
						composer: tt.expectedComposer,
						err:      tt.expectedErr,
					}
				},
			}

			composer, err := svc.Get(testContext(), 1)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedComposer, composer)
			}
		})
	}
}

func TestComposerService_List(t *testing.T) {
	tests := []struct {
		name              string
		expectedComposers []model.ComposerWithDetails
		storeErr          error
		expectedErr       error
	}{
		{
			name:              "store error",
			expectedComposers: nil,
			storeErr:          ErrFoo,
			expectedErr:       ErrFoo,
		},
		{
			name: "success",
			expectedComposers: []model.ComposerWithDetails{
				{
					Composer: content.Composer{
						ID:        1,
						FullName:  "Foo Foolington",
						ShortName: "Foolington",
					},
					PieceCount: 0,
				},
				{
					Composer: content.Composer{
						ID:        2,
						FullName:  "Bar Bartholomew",
						ShortName: "Bartholomew",
					},
					PieceCount: 1,
				},
				{
					Composer: content.Composer{
						ID:        3,
						FullName:  "Baz Bazura",
						ShortName: "Bazura",
					},
					PieceCount: 2,
				},
			},
			storeErr:    nil,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ComposerService{
				newComposerStore: func(db store.Executor) ComposerStore {
					return mockComposerStore{
						detailedComposers: tt.expectedComposers,
						err:               tt.storeErr,
					}
				},
			}

			composers, err := svc.List(testContext())

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedComposers, composers)
			}
		})
	}
}

func TestComposerService_Create(t *testing.T) {
	tempID := new(int)
	*tempID = 1

	tests := []struct {
		name             string
		cmd              model.ComposerCommand
		expectedComposer *content.Composer
		storeErr         error
		expectedErr      error
	}{
		{
			name: "operation mismatch",
			cmd: model.ComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationUpdate,
					TempID:    tempID,
					Data: content.Composer{
						FullName:  "Foo Foolington",
						ShortName: "Foolington",
					},
				},
			},
			expectedComposer: nil,
			storeErr:         content.ErrOperationMismatch,
			expectedErr:      content.ErrOperationMismatch,
		},
		{
			name: "invalid input composer",
			cmd: model.ComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationCreate,
					TempID:    tempID,
					Data:      content.Composer{},
				},
			},
			expectedComposer: nil,
			storeErr:         nil,
			expectedErr:      content.ErrInvalidResource,
		},
		{
			name: "store error",
			cmd: model.ComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationCreate,
					TempID:    tempID,
					Data: content.Composer{
						FullName:  "Foo Foolington",
						ShortName: "Foolington",
					},
				},
			},
			expectedComposer: nil,
			storeErr:         ErrFoo,
			expectedErr:      ErrFoo,
		},
		{
			name: "success",
			cmd: model.ComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationCreate,
					TempID:    tempID,
					Data: content.Composer{
						FullName:  "Foo Foolington",
						ShortName: "Foolington",
					},
				},
			},
			expectedComposer: &content.Composer{
				ID:        1,
				FullName:  "Foo Foolington",
				ShortName: "Foolington",
			},
			storeErr:    nil,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ComposerService{
				newComposerStore: func(db store.Executor) ComposerStore {
					return mockComposerStore{
						composer: tt.expectedComposer,
						err:      tt.storeErr,
					}
				},
			}

			composer, err := svc.Create(testContext(), tt.cmd)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedComposer, composer)
			}
		})
	}
}

func TestComposerService_Update(t *testing.T) {
	tests := []struct {
		name             string
		cmd              model.ComposerCommand
		expectedComposer *content.Composer
		storeErr         error
		expectedErr      error
	}{
		{
			name: "operation mismatch",
			cmd: model.ComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationCreate,
					Data: content.Composer{
						ID:        1,
						FullName:  "Foo Foolington",
						ShortName: "Foolington",
					},
				},
			},
			expectedComposer: nil,
			storeErr:         nil,
			expectedErr:      content.ErrOperationMismatch,
		},
		{
			name: "invalid input composer",
			cmd: model.ComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationUpdate,
					Data: content.Composer{
						ID: 1,
					},
				},
			},
			expectedComposer: nil,
			storeErr:         nil,
			expectedErr:      content.ErrInvalidResource,
		},
		{
			name: "store error",
			cmd: model.ComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationUpdate,
					Data: content.Composer{
						ID:        1,
						FullName:  "Foo Foolington",
						ShortName: "Foolington",
					},
				},
			},
			expectedComposer: nil,
			storeErr:         ErrFoo,
			expectedErr:      ErrFoo,
		},
		{
			name: "success",
			cmd: model.ComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationUpdate,
					Data: content.Composer{
						ID:        1,
						FullName:  "Foo Foolington",
						ShortName: "Foolington",
					},
				},
			},
			expectedComposer: &content.Composer{
				ID:        1,
				FullName:  "Foo Foolington",
				ShortName: "Foolington",
			},
			storeErr:    nil,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ComposerService{
				newComposerStore: func(db store.Executor) ComposerStore {
					return mockComposerStore{
						composer: tt.expectedComposer,
						err:      tt.expectedErr,
					}
				},
			}

			composer, err := svc.Update(testContext(), tt.cmd)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedComposer, composer)
			}
		})
	}
}

func TestComposerService_Delete(t *testing.T) {
	tests := []struct {
		name          string
		composer      *model.ComposerWithDetails
		getErr        error
		deleteErr     error
		expectedError error
	}{
		{
			name:          "get error",
			composer:      nil,
			getErr:        ErrGet,
			deleteErr:     nil,
			expectedError: ErrGet,
		},
		{
			name: "composer protected",
			composer: &model.ComposerWithDetails{
				Composer: content.Composer{
					ID:        2,
					FullName:  "Bar Bartholomew",
					ShortName: "Bartholomew",
				},
				PieceCount: 1,
			},
			getErr:        nil,
			deleteErr:     nil,
			expectedError: content.ErrComposerProtected,
		},
		{
			name: "delete error",
			composer: &model.ComposerWithDetails{
				Composer: content.Composer{
					ID:        2,
					FullName:  "Bar Bartholomew",
					ShortName: "Bartholomew",
				},
				PieceCount: 0,
			},
			getErr:        nil,
			deleteErr:     ErrDelete,
			expectedError: ErrDelete,
		},
		{
			name: "success",
			composer: &model.ComposerWithDetails{
				Composer: content.Composer{
					ID:        2,
					FullName:  "Bar Bartholomew",
					ShortName: "Bartholomew",
				},
				PieceCount: 0,
			},
			getErr:        nil,
			deleteErr:     nil,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ComposerService{
				newComposerStore: func(db store.Executor) ComposerStore {
					return mockComposerStore{
						detailedComposer: tt.composer,
						getErr:           tt.getErr,
						deleteErr:        tt.deleteErr,
					}
				},
			}

			err := svc.Delete(testContext(), 2)

			if tt.expectedError != nil {
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestComposerResolver_Run(t *testing.T) {
	tests := []struct {
		name        string
		intent      model.ComposerIntent
		expectedErr error
	}{
		{
			name: "invalid operation",
			intent: model.ComposerIntent{
				Operation: model.Operation("DELETE"),
				Data: content.Composer{
					ID: 1,
				},
			},
			expectedErr: model.ErrInvalidOperation,
		},
		{
			name: "select success",
			intent: model.ComposerIntent{
				Operation: model.OperationSelect,
				Data: content.Composer{
					ID:        1,
					FullName:  "Foo Bar Baz",
					ShortName: "Baz",
				},
			},
			expectedErr: nil,
		},
		{
			name: "create success",
			intent: model.ComposerIntent{
				Operation: model.OperationCreate,
				Data: content.Composer{
					FullName:  "Foo Bar Baz",
					ShortName: "Baz",
				},
			},
			expectedErr: nil,
		},
		{
			name: "update success",
			intent: model.ComposerIntent{
				Operation: model.OperationUpdate,
				Data: content.Composer{
					ID:        1,
					FullName:  "Foo Bar Baz",
					ShortName: "Baz",
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resolver := newComposerResolver(mockComposerStore{})

			_, err := resolver.run(testContext(), tt.intent)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

type mockComposerStore struct {
	composers         []content.Composer
	composer          *content.Composer
	detailedComposers []model.ComposerWithDetails
	detailedComposer  *model.ComposerWithDetails
	err               error
	getErr            error
	deleteErr         error
}

func (s mockComposerStore) Get(
	ctx context.Context,
	id int,
) (*content.Composer, error) {
	return s.composer, s.err
}

func (s mockComposerStore) GetWithDetails(
	ctx context.Context,
	id int,
) (*model.ComposerWithDetails, error) {
	return s.detailedComposer, s.getErr
}

func (s mockComposerStore) ListWithDetails(
	ctx context.Context,
) ([]model.ComposerWithDetails, error) {
	return s.detailedComposers, s.err
}

func (s mockComposerStore) Create(
	ctx context.Context,
	c content.Composer,
) (*content.Composer, error) {
	return s.composer, s.err
}

func (s mockComposerStore) Update(
	ctx context.Context,
	c content.Composer,
) (*content.Composer, error) {
	return s.composer, s.err
}

func (s mockComposerStore) Delete(
	ctx context.Context,
	id int,
) error {
	return s.deleteErr
}
