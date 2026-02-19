package service

import (
	"context"
	"testing"

	"github.com/adamkadda/arman/internal/cms/model"
	"github.com/adamkadda/arman/internal/cms/store"
	"github.com/adamkadda/arman/internal/content"
	"github.com/stretchr/testify/require"
)

func TestVenueService_Get(t *testing.T) {
	tests := []struct {
		name          string
		expectedVenue *content.Venue
		expectedErr   error
	}{
		{
			name: "success",
			expectedVenue: &content.Venue{
				ID:           1,
				Name:         "Foo Hall",
				FullAddress:  "11 Foo St. Foo City",
				ShortAddress: "11 Foo St.",
			},
			expectedErr: nil,
		},
		{
			name:          "store error",
			expectedVenue: nil,
			expectedErr:   ErrGet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := VenueService{
				newVenueStore: func(db store.Executor) VenueStore {
					return mockVenueStore{
						venue: tt.expectedVenue,
						err:   tt.expectedErr,
					}
				},
			}

			venue, err := svc.Get(testContext(), 1)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedVenue, venue)
			}
		})
	}
}

func TestVenueService_List(t *testing.T) {
	tests := []struct {
		name           string
		expectedVenues []model.VenueWithDetails
		expectedErr    error
	}{
		{
			name: "success",
			expectedVenues: []model.VenueWithDetails{
				{
					Venue: content.Venue{
						ID:           1,
						Name:         "Foo Hall",
						FullAddress:  "11 Foo St. Foo City",
						ShortAddress: "11 Foo St.",
					},
					EventCount: 0,
				},
				{
					Venue: content.Venue{
						ID:           2,
						Name:         "Bar Hall",
						FullAddress:  "22 Bar St. Bar City",
						ShortAddress: "22 Bar St.",
					},
					EventCount: 4,
				},
				{
					Venue: content.Venue{
						ID:           3,
						Name:         "Baz Hall",
						FullAddress:  "33 Baz St. Baz City",
						ShortAddress: "33 Baz St.",
					},
					EventCount: 9,
				},
			},
			expectedErr: nil,
		},
		{
			name:           "store error",
			expectedVenues: nil,
			expectedErr:    ErrFoo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := VenueService{
				newVenueStore: func(db store.Executor) VenueStore {
					return mockVenueStore{
						detailedVenues: tt.expectedVenues,
						err:            tt.expectedErr,
					}
				},
			}

			venues, err := svc.List(testContext())

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedVenues, venues)
			}
		})
	}
}

func TestVenueService_Create(t *testing.T) {
	tests := []struct {
		name        string
		cmd         model.VenueCommand
		venue       *content.Venue
		storeErr    error
		expectedErr error
	}{
		{
			name: "operation mismatch",
			cmd: model.VenueCommand{
				Venue: model.VenueIntent{
					Operation: model.OperationUpdate,
					Data: content.Venue{
						Name:         "Foo Hall",
						FullAddress:  "11 Foo St. Foo City",
						ShortAddress: "11 Foo St.",
					},
				},
			},
			venue:       nil,
			storeErr:    nil,
			expectedErr: content.ErrOperationMismatch,
		},
		{
			name: "invalid input venue",
			cmd: model.VenueCommand{
				Venue: model.VenueIntent{
					Operation: model.OperationCreate,
					Data:      content.Venue{},
				},
			},
			venue:       nil,
			storeErr:    nil,
			expectedErr: content.ErrInvalidResource,
		},
		{
			name: "store error",
			cmd: model.VenueCommand{
				Venue: model.VenueIntent{
					Operation: model.OperationCreate,
					Data: content.Venue{
						Name:         "Foo Hall",
						FullAddress:  "11 Foo St. Foo City",
						ShortAddress: "11 Foo St.",
					},
				},
			},
			venue:       nil,
			storeErr:    ErrFoo,
			expectedErr: ErrFoo,
		},
		{
			name: "success",
			cmd: model.VenueCommand{
				Venue: model.VenueIntent{
					Operation: model.OperationCreate,
					Data: content.Venue{
						Name:         "Foo Hall",
						FullAddress:  "11 Foo St. Foo City",
						ShortAddress: "11 Foo St.",
					},
				},
			},
			venue: &content.Venue{
				ID:           1,
				Name:         "Foo Hall",
				FullAddress:  "11 Foo St. Foo City",
				ShortAddress: "11 Foo St.",
			},
			storeErr:    nil,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := VenueService{
				newVenueStore: func(db store.Executor) VenueStore {
					return mockVenueStore{
						venue: tt.venue,
						err:   tt.storeErr,
					}
				},
			}

			venue, err := svc.Create(testContext(), tt.cmd)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.venue, venue)
			}
		})
	}
}

func TestVenueService_Update(t *testing.T) {
	tests := []struct {
		name        string
		cmd         model.VenueCommand
		venue       *content.Venue
		storeErr    error
		expectedErr error
	}{
		{
			name: "operation mismatch",
			cmd: model.VenueCommand{
				Venue: model.VenueIntent{
					Operation: model.OperationCreate,
					Data: content.Venue{
						ID:           1,
						Name:         "Foo Hall",
						FullAddress:  "11 Foo St. Foo City",
						ShortAddress: "11 Foo St.",
					},
				},
			},
			venue:       nil,
			storeErr:    nil,
			expectedErr: content.ErrOperationMismatch,
		},
		{
			name: "invalid input venue",
			cmd: model.VenueCommand{
				Venue: model.VenueIntent{
					Operation: model.OperationUpdate,
					Data: content.Venue{
						ID: 1,
					},
				},
			},
			venue:       nil,
			storeErr:    nil,
			expectedErr: content.ErrInvalidResource,
		},
		{
			name: "store error",
			cmd: model.VenueCommand{
				Venue: model.VenueIntent{
					Operation: model.OperationUpdate,
					Data: content.Venue{
						ID:           1,
						Name:         "Foo Hall",
						FullAddress:  "11 Foo St. Foo City",
						ShortAddress: "11 Foo St.",
					},
				},
			},
			venue:       nil,
			storeErr:    ErrFoo,
			expectedErr: ErrFoo,
		},
		{
			name: "success",
			cmd: model.VenueCommand{
				Venue: model.VenueIntent{
					Operation: model.OperationUpdate,
					Data: content.Venue{
						ID:           1,
						Name:         "Foo Hall",
						FullAddress:  "22 Bar St. Baz Town",
						ShortAddress: "22 Bar St.",
					},
				},
			},
			venue: &content.Venue{
				ID:           1,
				Name:         "Foo Hall",
				FullAddress:  "11 Foo St. Foo City",
				ShortAddress: "11 Foo St.",
			},
			storeErr:    nil,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := VenueService{
				newVenueStore: func(db store.Executor) VenueStore {
					return mockVenueStore{
						venue: tt.venue,
						err:   tt.storeErr,
					}
				},
			}

			venue, err := svc.Update(testContext(), tt.cmd)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.venue, venue)
			}
		})
	}
}

func TestVenueService_Delete(t *testing.T) {
	tests := []struct {
		name        string
		venue       *model.VenueWithDetails
		getErr      error
		deleteErr   error
		expectedErr error
	}{
		{
			name:        "get error",
			venue:       nil,
			getErr:      ErrGet,
			deleteErr:   nil,
			expectedErr: ErrGet,
		},
		{
			name: "venue protected",
			venue: &model.VenueWithDetails{
				Venue: content.Venue{
					ID:           2,
					Name:         "Bar Hall",
					FullAddress:  "22 Bar St. Bar City",
					ShortAddress: "22 Bar St.",
				},
				EventCount: 4,
			},
			getErr:      nil,
			deleteErr:   nil,
			expectedErr: content.ErrVenueProtected,
		},
		{
			name: "delete error",
			venue: &model.VenueWithDetails{
				Venue: content.Venue{
					ID:           1,
					Name:         "Foo Hall",
					FullAddress:  "11 Foo St. Foo City",
					ShortAddress: "11 Foo St.",
				},
				EventCount: 0,
			},
			getErr:      nil,
			deleteErr:   ErrDelete,
			expectedErr: ErrDelete,
		},
		{
			name: "success",
			venue: &model.VenueWithDetails{
				Venue: content.Venue{
					ID:           1,
					Name:         "Foo Hall",
					FullAddress:  "11 Foo St. Foo City",
					ShortAddress: "11 Foo St.",
				},
				EventCount: 0,
			},
			getErr:      nil,
			deleteErr:   nil,
			expectedErr: nil,
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

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestVenueResolver_Run(t *testing.T) {
	tests := []struct {
		name        string
		intent      model.VenueIntent
		expectedErr error
	}{
		{
			name: "invalid operation",
			intent: model.VenueIntent{
				Operation: model.Operation("DELETE"),
				Data: content.Venue{
					ID: 1,
				},
			},
			expectedErr: model.ErrInvalidOperation,
		},
		{
			name: "select success",
			intent: model.VenueIntent{
				Operation: model.OperationSelect,
				Data: content.Venue{
					ID:           1,
					Name:         "Foo Hall",
					FullAddress:  "22 Foo St. Bar City",
					ShortAddress: "22 Foo St.",
				},
			},
			expectedErr: nil,
		},
		{
			name: "create success",
			intent: model.VenueIntent{
				Operation: model.OperationCreate,
				Data: content.Venue{
					Name:         "Foo Hall",
					FullAddress:  "22 Foo St. Bar City",
					ShortAddress: "22 Foo St.",
				},
			},
			expectedErr: nil,
		},
		{
			name: "update success",
			intent: model.VenueIntent{
				Operation: model.OperationUpdate,
				Data: content.Venue{
					ID:           1,
					Name:         "Foo Hall",
					FullAddress:  "22 Foo St. Bar City",
					ShortAddress: "22 Foo St.",
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resolver := newVenueResolver(mockVenueStore{})

			_, err := resolver.run(testContext(), tt.intent)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

type mockVenueStore struct {
	venues         []content.Venue
	venue          *content.Venue
	detailedVenue  *model.VenueWithDetails
	detailedVenues []model.VenueWithDetails
	err            error
	getErr         error
	deleteErr      error
}

func (s mockVenueStore) Get(
	ctx context.Context,
	id int,
) (*content.Venue, error) {
	return s.venue, s.err
}

func (s mockVenueStore) GetWithDetails(
	ctx context.Context,
	id int,
) (*model.VenueWithDetails, error) {
	return s.detailedVenue, s.getErr
}

func (s mockVenueStore) ListWithDetails(
	ctx context.Context,
) ([]model.VenueWithDetails, error) {
	return s.detailedVenues, s.err
}

func (s mockVenueStore) Create(
	ctx context.Context,
	v content.Venue,
) (*content.Venue, error) {
	return s.venue, s.err
}

func (s mockVenueStore) Update(
	ctx context.Context,
	v content.Venue,
) (*content.Venue, error) {
	return s.venue, s.err
}

func (s mockVenueStore) Delete(
	ctx context.Context,
	id int,
) error {
	return s.deleteErr
}
