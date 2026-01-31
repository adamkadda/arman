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

type mockVenueStore struct {
	venues         []content.Venue
	venue          *content.Venue
	detailedVenue  *models.VenueWithDetails
	detailedVenues []models.VenueWithDetails
	err            error
	getErr         error
	deleteErr      error
}

func (s mockVenueStore) Get(ctx context.Context, id int) (*content.Venue, error) {
	return s.venue, s.err
}

func (s mockVenueStore) GetWithDetails(ctx context.Context, id int) (*models.VenueWithDetails, error) {
	return s.detailedVenue, s.getErr
}

func (s mockVenueStore) ListWithDetails(ctx context.Context) ([]models.VenueWithDetails, error) {
	return s.detailedVenues, s.err
}

func (s mockVenueStore) Create(ctx context.Context, v content.Venue) (*content.Venue, error) {
	return s.venue, s.err
}

func (s mockVenueStore) Update(ctx context.Context, v content.Venue) (*content.Venue, error) {
	return s.venue, s.err
}

func (s mockVenueStore) Delete(ctx context.Context, id int) error {
	return s.deleteErr
}

func TestVenueService_Get(t *testing.T) {
	tests := []struct {
		name    string
		venue   *content.Venue
		err     error
		wantErr bool
	}{
		{"venue.get success", &content.Venue{Name: "foo"}, nil, false},
		{"venue.get error", nil, errors.New("oops"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := VenueService{
				newVenueStore: func(db store.Executor) VenueStore {
					return mockVenueStore{
						venue: tt.venue,
						err:   tt.err,
					}
				},
			}

			venue, err := svc.Get(testContext(), 1)

			if tt.wantErr {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, venue)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.venue, venue)
			}
		})
	}
}

func TestVenueService_List(t *testing.T) {
	tests := []struct {
		name    string
		venues  []models.VenueWithDetails
		err     error
		wantErr bool
	}{
		{"venue.list success", []models.VenueWithDetails{
			{
				Venue:      content.Venue{Name: "foo"},
				EventCount: 0,
			},
			{
				Venue:      content.Venue{Name: "bar"},
				EventCount: 1,
			},
			{
				Venue:      content.Venue{Name: "baz"},
				EventCount: 3,
			},
		}, nil, false},
		{"venue.list error", nil, errors.New("oops"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := VenueService{
				newVenueStore: func(db store.Executor) VenueStore {
					return mockVenueStore{
						detailedVenues: tt.venues,
						err:            tt.err,
					}
				},
			}

			venues, err := svc.List(testContext())

			if tt.wantErr {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, venues)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.venues, venues)
			}
		})
	}
}

func TestVenueService_Create(t *testing.T) {
	tests := []struct {
		name    string
		venue   content.Venue
		err     error
		wantErr bool
	}{
		{"venue.create success", content.Venue{
			Name:         "Foo Hall",
			FullAddress:  "22 Bar Street. Baz Town",
			ShortAddress: "22 Bar Street.",
		}, nil, false},
		{"venue.create rejected", content.Venue{}, content.ErrInvalidResource, true},
		{"venue.create error", content.Venue{
			Name:         "Foo Hall",
			FullAddress:  "22 Bar Street. Baz Town",
			ShortAddress: "22 Bar Street.",
		}, errors.New("oops"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := VenueService{
				newVenueStore: func(db store.Executor) VenueStore {
					return mockVenueStore{
						venue: &tt.venue,
						err:   tt.err,
					}
				},
			}

			venue, err := svc.Create(testContext(), tt.venue)

			if tt.wantErr {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, venue)
			} else {
				require.NoError(t, err)
				require.Equal(t, &tt.venue, venue)
			}
		})
	}
}

func TestVenueService_Update(t *testing.T) {
	tests := []struct {
		name    string
		venue   content.Venue
		err     error
		wantErr bool
	}{
		{"venue.update success", content.Venue{
			Name:         "Foo Hall",
			FullAddress:  "22 Bar Street. Baz Town",
			ShortAddress: "22 Bar Street.",
		}, nil, false},
		{"venue.update rejected", content.Venue{}, content.ErrInvalidResource, true},
		{"venue.update error", content.Venue{
			Name:         "Foo Hall",
			FullAddress:  "22 Bar Street. Baz Town",
			ShortAddress: "22 Bar Street.",
		}, errors.New("oops"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := VenueService{
				newVenueStore: func(db store.Executor) VenueStore {
					return mockVenueStore{
						venue: &tt.venue,
						err:   tt.err,
					}
				},
			}

			venue, err := svc.Update(testContext(), tt.venue)

			if tt.wantErr {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, venue)
			} else {
				require.NoError(t, err)
				require.Equal(t, &tt.venue, venue)
			}
		})
	}
}

func TestVenueService_Delete(t *testing.T) {
	tests := []struct {
		name          string
		venue         *models.VenueWithDetails
		getErr        error
		deleteErr     error
		expectedError error
	}{
		{
			name: "venue.delete success",
			venue: &models.VenueWithDetails{
				Venue:      content.Venue{Name: "foo"},
				EventCount: 0,
			},
			getErr:        nil,
			deleteErr:     nil,
			expectedError: nil,
		},
		{
			name:          "venue.get_with_details error",
			venue:         nil,
			getErr:        ErrGet,
			deleteErr:     nil,
			expectedError: ErrGet,
		},
		{
			name: "venue.delete blocked",
			venue: &models.VenueWithDetails{
				Venue:      content.Venue{Name: "foo"},
				EventCount: 1,
			},
			getErr:        nil,
			deleteErr:     nil,
			expectedError: content.ErrVenueProtected,
		},
		{
			name: "venue.delete error",
			venue: &models.VenueWithDetails{
				Venue:      content.Venue{Name: "foo"},
				EventCount: 0,
			},
			getErr:        nil,
			deleteErr:     ErrDelete,
			expectedError: ErrDelete,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := VenueService{
				newVenueStore: func(db store.Executor) VenueStore {
					return mockVenueStore{
						detailedVenue: tt.venue,
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
