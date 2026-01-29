package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/adamkadda/arman/internal/cms/store"
	"github.com/adamkadda/arman/internal/content"
	"github.com/adamkadda/arman/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BiographyService struct {
	pool *pgxpool.Pool
}

func NewBiographyService(pool *pgxpool.Pool) *BiographyService {
	return &BiographyService{
		pool: pool,
	}
}

func (s *BiographyService) Get(
	ctx context.Context,
	variant content.BiographyVariant,
) (*content.Biography, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "biography.get"),
		slog.String("variant", string(variant)),
	)

	logger.Info(
		"get composer",
	)

	if err := variant.Validate(); err != nil {
		logger.Warn(
			"validate variant failed",
			slog.String("reason", reason(err)),
		)

		return nil, fmt.Errorf("%w: %s", content.ErrInvalidBiographyVariant, err)
	}

	biographyStore := store.NewBiographyStore(s.pool)

	biography, err := biographyStore.Get(ctx, variant)
	if err != nil {
		logger.Error(
			"get biography failed",
			slog.String("step", "biography.get"),
			slog.Any("error", err),
		)

		return nil, err
	}

	return biography, nil
}

func (s *BiographyService) Update(
	ctx context.Context,
	b content.Biography,
) (*content.Biography, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "biography.get"),
		slog.String("variant", string(b.Variant)),
	)

	logger.Info(
		"get composer",
	)

	if err := b.Variant.Validate(); err != nil {
		logger.Warn(
			"validate variant failed",
			slog.String("reason", reason(err)),
		)

		return nil, fmt.Errorf("%w: %s", content.ErrInvalidBiographyVariant, err)
	}

	biographyStore := store.NewBiographyStore(s.pool)

	biography, err := biographyStore.Update(ctx, b)
	if err != nil {
		logger.Error(
			"get biography failed",
			slog.String("step", "biography.get"),
			slog.Any("error", err),
		)

		return nil, err
	}

	return biography, nil
}
