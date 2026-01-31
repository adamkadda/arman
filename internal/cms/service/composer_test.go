package service

import (
	"context"
	"errors"
	"testing"

	"github.com/adamkadda/arman/internal/cms/models"
	"github.com/adamkadda/arman/internal/cms/store"
	"github.com/adamkadda/arman/internal/content"
	"github.com/stretchr/testify/require"
)

type mockComposerStore struct {
	composers         []content.Composer
	composer          *content.Composer
	detailedComposers []models.ComposerWithDetails
	detailedComposer  *models.ComposerWithDetails
	err               error
}

func (s mockComposerStore) Get(ctx context.Context, id int) (*content.Composer, error) {
	return s.composer, s.err
}

func (s mockComposerStore) GetWithDetails(ctx context.Context, id int) (*models.ComposerWithDetails, error) {
	return s.detailedComposer, s.err
}

func (s mockComposerStore) ListWithDetails(ctx context.Context) ([]models.ComposerWithDetails, error) {
	return s.detailedComposers, s.err
}

func (s mockComposerStore) Create(ctx context.Context, c content.Composer) (*content.Composer, error) {
	return s.composer, s.err
}

func (s mockComposerStore) Update(ctx context.Context, c content.Composer) (*content.Composer, error) {
	return s.composer, s.err
}

func (s mockComposerStore) Delete(ctx context.Context, id int) error {
	return s.err
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
				require.Error(t, err)
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
		composers []models.ComposerWithDetails
		err       error
		wantErr   bool
	}{
		{"composer.list success", []models.ComposerWithDetails{
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
				require.Error(t, err)
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
		composer content.Composer
		err      error
		wantErr  bool
	}{
		{"composer.create success", content.Composer{
			FullName:  "Foo",
			ShortName: "Bar",
		}, nil, false},
		{"composer.create rejected", content.Composer{}, content.ErrInvalidResource, true},
		{"composer.create error", content.Composer{}, errors.New("oops"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ComposerService{
				newComposerStore: func(db store.Executor) ComposerStore {
					return mockComposerStore{
						composer: &tt.composer,
						err:      tt.err,
					}
				},
			}

			composer, err := svc.Create(testContext(), tt.composer)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, composer)
			} else {
				require.NoError(t, err)
				require.Equal(t, &tt.composer, composer)
			}
		})
	}
}

func TestComposerService_Update(t *testing.T) {
	tests := []struct {
		name     string
		composer content.Composer
		err      error
		wantErr  bool
	}{
		{"composer.update success", content.Composer{
			FullName:  "Foo",
			ShortName: "Bar",
		}, nil, false},
		{"composer.update rejected", content.Composer{}, content.ErrInvalidResource, true},
		{"composer.update error", content.Composer{}, errors.New("oops"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ComposerService{
				newComposerStore: func(db store.Executor) ComposerStore {
					return mockComposerStore{
						composer: &tt.composer,
						err:      tt.err,
					}
				},
			}

			composer, err := svc.Update(testContext(), tt.composer)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, composer)
			} else {
				require.NoError(t, err)
				require.Equal(t, &tt.composer, composer)
			}
		})
	}
}

func TestComposerService_Delete(t *testing.T) {
	tests := []struct {
		name     string
		composer *models.ComposerWithDetails
		err      error
		wantErr  bool
	}{
		{"composer.delete success", &models.ComposerWithDetails{
			Composer:   content.Composer{FullName: "foo"},
			PieceCount: 0,
		}, nil, false},
		{"composer.get_with_details error", &models.ComposerWithDetails{}, errors.New("oops"), true},
		{"composer.delete blocked", &models.ComposerWithDetails{
			Composer:   content.Composer{FullName: "foo"},
			PieceCount: 1,
		}, content.ErrComposerProtected, true},
		{"composer.delete error", &models.ComposerWithDetails{}, errors.New("oops"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := ComposerService{
				newComposerStore: func(db store.Executor) ComposerStore {
					return mockComposerStore{
						detailedComposer: tt.composer,
						err:              tt.err,
					}
				},
			}

			err := svc.Delete(testContext(), 1)

			if tt.wantErr {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
