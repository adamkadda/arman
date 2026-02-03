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

type mockComposerStore struct {
	composers         []content.Composer
	composer          *content.Composer
	detailedComposers []model.ComposerWithDetails
	detailedComposer  *model.ComposerWithDetails
	err               error
	getErr            error
	deleteErr         error
}

func (s mockComposerStore) Get(ctx context.Context, id int) (*content.Composer, error) {
	return s.composer, s.err
}

func (s mockComposerStore) GetWithDetails(ctx context.Context, id int) (*model.ComposerWithDetails, error) {
	return s.detailedComposer, s.getErr
}

func (s mockComposerStore) ListWithDetails(ctx context.Context) ([]model.ComposerWithDetails, error) {
	return s.detailedComposers, s.err
}

func (s mockComposerStore) Create(ctx context.Context, c content.Composer) (*content.Composer, error) {
	return s.composer, s.err
}

func (s mockComposerStore) Update(ctx context.Context, c content.Composer) (*content.Composer, error) {
	return s.composer, s.err
}

func (s mockComposerStore) Delete(ctx context.Context, id int) error {
	return s.deleteErr
}

func TestComposerService_Get(t *testing.T) {
	tests := []struct {
		name     string
		composer *content.Composer
		err      error
		wantErr  bool
	}{
		{"composer.get success", &content.Composer{FullName: "foo"}, nil, false},
		{"composer.get error", nil, errors.New("oops"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ComposerService{
				newComposerStore: func(db store.Executor) ComposerStore {
					return mockComposerStore{
						composer: tt.composer,
						err:      tt.err,
					}
				},
			}

			composer, err := svc.Get(testContext(), 1)

			if tt.wantErr {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, composer)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.composer, composer)
			}
		})
	}
}

func TestComposerService_List(t *testing.T) {
	tests := []struct {
		name      string
		composers []model.ComposerWithDetails
		err       error
		wantErr   bool
	}{
		{"composer.list success", []model.ComposerWithDetails{
			{
				Composer:   content.Composer{FullName: "foo"},
				PieceCount: 0,
			},
			{
				Composer:   content.Composer{FullName: "bar"},
				PieceCount: 1,
			},
			{
				Composer:   content.Composer{FullName: "baz"},
				PieceCount: 2,
			},
		}, nil, false},
		{"composer.list error", nil, errors.New("oops"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ComposerService{
				newComposerStore: func(db store.Executor) ComposerStore {
					return mockComposerStore{
						detailedComposers: tt.composers,
						err:               tt.err,
					}
				},
			}

			composers, err := svc.List(testContext())

			if tt.wantErr {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, composers)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.composers, composers)
			}
		})
	}
}

func TestComposerService_Create(t *testing.T) {
	tests := []struct {
		name     string
		cmd      model.UpsertComposerCommand
		composer *content.Composer
		err      error
		wantErr  bool
	}{
		{
			"operation mismatch",
			model.UpsertComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationUpdate,
					Data: content.Composer{
						FullName:  "Foo Bar",
						ShortName: "Bar",
					},
				},
			},
			nil,
			content.ErrOperationMismatch,
			true,
		},
		{
			"invalid input composer",
			model.UpsertComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationCreate,
					Data:      content.Composer{},
				},
			},
			nil,
			content.ErrInvalidResource,
			true,
		},
		{
			"store error",
			model.UpsertComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationCreate,
					Data: content.Composer{
						FullName:  "Foo Bar",
						ShortName: "Bar",
					},
				},
			},
			nil,
			ErrFoo,
			true,
		},
		{
			"success",
			model.UpsertComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationCreate,
					Data: content.Composer{
						FullName:  "Foo Bar",
						ShortName: "Bar",
					},
				},
			},
			&content.Composer{
				ID:        1,
				FullName:  "Foo Bar",
				ShortName: "Bar",
			},
			nil,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ComposerService{
				newComposerStore: func(db store.Executor) ComposerStore {
					return mockComposerStore{
						composer: tt.composer,
						err:      tt.err,
					}
				},
			}

			composer, err := svc.Create(testContext(), tt.cmd)

			if tt.wantErr {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, composer)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.composer, composer)
			}
		})
	}
}

func TestComposerService_Update(t *testing.T) {
	tests := []struct {
		name     string
		cmd      model.UpsertComposerCommand
		composer *content.Composer
		err      error
		wantErr  bool
	}{
		{
			"operation mismatch",
			model.UpsertComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationCreate,
					Data: content.Composer{
						FullName:  "Foo Bar",
						ShortName: "Bar",
					},
				},
			},
			nil,
			content.ErrOperationMismatch,
			true,
		},
		{
			"invalid input composer",
			model.UpsertComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationUpdate,
					Data:      content.Composer{},
				},
			},
			nil,
			content.ErrInvalidResource,
			true,
		},
		{
			"store error",
			model.UpsertComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationUpdate,
					Data: content.Composer{
						FullName:  "Foo Bar",
						ShortName: "Bar",
					},
				},
			},
			nil,
			ErrFoo,
			true,
		},
		{
			"success",
			model.UpsertComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationUpdate,
					Data: content.Composer{
						ID:        1,
						FullName:  "Foo Bar",
						ShortName: "Bar",
					},
				},
			},
			&content.Composer{
				ID:        1,
				FullName:  "Foo Bar",
				ShortName: "Bar",
			},
			nil,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ComposerService{
				newComposerStore: func(db store.Executor) ComposerStore {
					return mockComposerStore{
						composer: tt.composer,
						err:      tt.err,
					}
				},
			}

			composer, err := svc.Update(testContext(), tt.cmd)

			if tt.wantErr {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, composer)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.composer, composer)
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
			name: "composer.delete success",
			composer: &model.ComposerWithDetails{
				Composer:   content.Composer{FullName: "foo"},
				PieceCount: 0,
			},
			getErr:        nil,
			deleteErr:     nil,
			expectedError: nil,
		},
		{
			name:          "composer.get_with_details error",
			composer:      nil,
			getErr:        ErrGet,
			deleteErr:     nil,
			expectedError: ErrGet,
		},
		{
			name: "composer.delete blocked",
			composer: &model.ComposerWithDetails{
				Composer:   content.Composer{FullName: "foo"},
				PieceCount: 1,
			},
			getErr:        nil,
			deleteErr:     nil,
			expectedError: content.ErrComposerProtected,
		},
		{
			name: "composer.delete error",
			composer: &model.ComposerWithDetails{
				Composer:   content.Composer{FullName: "foo"},
				PieceCount: 0,
			},
			getErr:        nil,
			deleteErr:     ErrDelete,
			expectedError: ErrDelete,
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

			err := svc.Delete(testContext(), 1)

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
			"select success",
			model.ComposerIntent{
				Operation: model.OperationSelect,
				Data: content.Composer{
					ID:        1,
					FullName:  "Foo Bar Baz",
					ShortName: "Baz",
				},
			},
			nil,
		},
		{
			"create success",
			model.ComposerIntent{
				Operation: model.OperationCreate,
				Data: content.Composer{
					FullName:  "Foo Bar Baz",
					ShortName: "Baz",
				},
			},
			nil,
		},
		{
			"update success",
			model.ComposerIntent{
				Operation: model.OperationUpdate,
				Data: content.Composer{
					ID:        1,
					FullName:  "Foo Bar Baz",
					ShortName: "Baz",
				},
			},
			nil,
		},
		{
			"invalid operation",
			model.ComposerIntent{
				Operation: model.Operation("DELETE"),
				Data: content.Composer{
					ID:        1,
					FullName:  "Foo Bar Baz",
					ShortName: "Baz",
				},
			},
			model.ErrInvalidOperation,
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
