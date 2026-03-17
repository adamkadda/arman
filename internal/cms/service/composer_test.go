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

func TestComposerService_Get(t *testing.T) {
	tests := []struct {
		name        string
		store       mockComposerStore
		expected    *content.Composer
		expectedErr error
	}{
		{
			name: "store error",
			store: mockComposerStore{
				err: ErrGet,
			},
			expected:    nil,
			expectedErr: ErrGet,
		},
		{
			name: "success",
			store: mockComposerStore{
				composers: map[int]*content.Composer{
					1: {ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
				},
			},
			expected:    &content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ComposerService{
				newComposerStore: func(db store.Executor) ComposerStore {
					return tt.store
				},
			}

			composer, err := svc.Get(testContext(), 1)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, composer)
			}
		})
	}
}

func TestComposerService_List(t *testing.T) {
	tests := []struct {
		name        string
		store       mockComposerStore
		expected    []model.ComposerWithDetails
		expectedErr error
	}{
		{
			name: "store error",
			store: mockComposerStore{
				err: ErrFoo,
			},
			expected:    nil,
			expectedErr: ErrFoo,
		},
		{
			name: "success",
			store: mockComposerStore{
				detailedComposers: map[int]*model.ComposerWithDetails{
					1: {Composer: content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"}, PieceCount: 0},
					2: {Composer: content.Composer{ID: 2, FullName: "Bar Bartholomew", ShortName: "Bartholomew"}, PieceCount: 1},
					3: {Composer: content.Composer{ID: 3, FullName: "Baz Bazura", ShortName: "Bazura"}, PieceCount: 2},
				},
			},
			expected: []model.ComposerWithDetails{
				{Composer: content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"}, PieceCount: 0},
				{Composer: content.Composer{ID: 2, FullName: "Bar Bartholomew", ShortName: "Bartholomew"}, PieceCount: 1},
				{Composer: content.Composer{ID: 3, FullName: "Baz Bazura", ShortName: "Bazura"}, PieceCount: 2},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ComposerService{
				newComposerStore: func(db store.Executor) ComposerStore {
					return tt.store
				},
			}

			composers, err := svc.List(testContext())

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, composers)
			}
		})
	}
}

func TestComposerService_Create(t *testing.T) {
	tempID := new(int)
	*tempID = 1

	tests := []struct {
		name        string
		cmd         model.ComposerCommand
		store       mockComposerStore
		expected    *content.Composer
		expectedErr error
	}{
		{
			name: "operation mismatch",
			cmd: model.ComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationUpdate,
					TempID:    tempID,
					Data:      content.Composer{FullName: "Foo Foolington", ShortName: "Foolington"},
				},
			},
			store:       mockComposerStore{},
			expected:    nil,
			expectedErr: content.ErrOperationMismatch,
		},
		{
			name: "invalid composer",
			cmd: model.ComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationCreate,
					TempID:    tempID,
					Data:      content.Composer{},
				},
			},
			store:       mockComposerStore{},
			expected:    nil,
			expectedErr: content.ErrInvalidResource,
		},
		{
			name: "success",
			cmd: model.ComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationCreate,
					TempID:    tempID,
					Data:      content.Composer{FullName: "Foo Foolington", ShortName: "Foolington"},
				},
			},
			store: mockComposerStore{
				composers: map[int]*content.Composer{
					0: {ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
				},
			},
			expected:    &content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ComposerService{
				newComposerStore: func(db store.Executor) ComposerStore {
					return tt.store
				},
			}

			composer, err := svc.Create(testContext(), tt.cmd)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, composer)
			}
		})
	}
}

func TestComposerService_Update(t *testing.T) {
	tests := []struct {
		name        string
		cmd         model.ComposerCommand
		store       mockComposerStore
		expected    *content.Composer
		expectedErr error
	}{
		{
			name: "operation mismatch",
			cmd: model.ComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationCreate,
					Data:      content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
				},
			},
			store:       mockComposerStore{},
			expected:    nil,
			expectedErr: content.ErrOperationMismatch,
		},
		{
			name: "invalid composer",
			cmd: model.ComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationUpdate,
					Data:      content.Composer{ID: 1},
				},
			},
			store:       mockComposerStore{},
			expected:    nil,
			expectedErr: content.ErrInvalidResource,
		},
		{
			name: "success",
			cmd: model.ComposerCommand{
				Composer: model.ComposerIntent{
					Operation: model.OperationUpdate,
					Data:      content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
				},
			},
			store: mockComposerStore{
				composers: map[int]*content.Composer{
					1: {ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
				},
			},
			expected:    &content.Composer{ID: 1, FullName: "Foo Foolington", ShortName: "Foolington"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ComposerService{
				newComposerStore: func(db store.Executor) ComposerStore {
					return tt.store
				},
			}

			composer, err := svc.Update(testContext(), tt.cmd)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, composer)
			}
		})
	}
}

func TestComposerService_Delete(t *testing.T) {
	tests := []struct {
		name        string
		store       mockComposerStore
		expectedErr error
	}{
		{
			name: "get error",
			store: mockComposerStore{
				err: ErrGet,
			},
			expectedErr: ErrGet,
		},
		{
			name: "composer protected",
			store: mockComposerStore{
				detailedComposers: map[int]*model.ComposerWithDetails{
					2: {Composer: content.Composer{ID: 2, FullName: "Bar Bartholomew", ShortName: "Bartholomew"}, PieceCount: 1},
				},
			},
			expectedErr: content.ErrComposerProtected,
		},
		{
			name: "success",
			store: mockComposerStore{
				detailedComposers: map[int]*model.ComposerWithDetails{
					2: {Composer: content.Composer{ID: 2, FullName: "Bar Bartholomew", ShortName: "Bartholomew"}, PieceCount: 0},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ComposerService{
				newComposerStore: func(db store.Executor) ComposerStore {
					return tt.store
				},
			}

			err := svc.Delete(testContext(), 2)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
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
				Data:      content.Composer{ID: 1},
			},
			expectedErr: model.ErrInvalidOperation,
		},
		{
			name: "select success",
			intent: model.ComposerIntent{
				Operation: model.OperationSelect,
				Data:      content.Composer{ID: 1, FullName: "Foo Bar Baz", ShortName: "Baz"},
			},
			expectedErr: nil,
		},
		{
			name: "create success",
			intent: model.ComposerIntent{
				Operation: model.OperationCreate,
				Data:      content.Composer{FullName: "Foo Bar Baz", ShortName: "Baz"},
			},
			expectedErr: nil,
		},
		{
			name: "update success",
			intent: model.ComposerIntent{
				Operation: model.OperationUpdate,
				Data:      content.Composer{ID: 1, FullName: "Foo Bar Baz", ShortName: "Baz"},
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
	composers         map[int]*content.Composer
	detailedComposers map[int]*model.ComposerWithDetails
	err               error
}

func (s mockComposerStore) Get(
	ctx context.Context,
	id int,
) (*content.Composer, error) {
	return s.composers[id], s.err
}

func (s mockComposerStore) GetWithDetails(
	ctx context.Context,
	id int,
) (*model.ComposerWithDetails, error) {
	return s.detailedComposers[id], s.err
}

func (s mockComposerStore) ListWithDetails(
	ctx context.Context,
) ([]model.ComposerWithDetails, error) {
	if s.err != nil {
		return nil, s.err
	}

	composers := make([]model.ComposerWithDetails, 0, len(s.detailedComposers))
	for _, c := range s.detailedComposers {
		composers = append(composers, *c)
	}

	slices.SortFunc(composers, func(a, b model.ComposerWithDetails) int {
		return cmp.Compare(a.Composer.ID, b.Composer.ID)
	})

	return composers, nil
}

func (s mockComposerStore) Create(
	ctx context.Context,
	c content.Composer,
) (*content.Composer, error) {
	return s.composers[0], s.err
}

func (s mockComposerStore) Update(
	ctx context.Context,
	c content.Composer,
) (*content.Composer, error) {
	return s.composers[c.ID], s.err
}

func (s mockComposerStore) Delete(
	ctx context.Context,
	id int,
) error {
	return s.err
}
