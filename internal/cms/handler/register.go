package handler

import (
	"net/http"

	"github.com/adamkadda/arman/internal/cms/service"
	"github.com/adamkadda/arman/pkg/logging"
	"github.com/adamkadda/arman/pkg/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(
	pool *pgxpool.Pool,
) http.Handler {
	layers := []middleware.Middleware{
		logging.Middleware(),
	}

	// TODO: Append middleware to layers based on app stage.

	stack := middleware.NewStack(layers...)

	router := http.NewServeMux()

	venueService := service.NewVenueService(pool)
	venueHandler := NewVenueHandler(venueService)
	venueHandler.Register(router)

	composerService := service.NewComposerService(pool)
	composerHandler := NewComposerHandler(composerService)
	composerHandler.Register(router)

	pieceService := service.NewPieceService(pool)
	pieceHandler := NewPieceHandler(pieceService)
	pieceHandler.Register(router)

	programmeService := service.NewProgrammeService(pool)
	programmeHandler := NewProgrammeHandler(programmeService)
	programmeHandler.Register(router)

	eventService := service.NewEventService(pool)
	eventHandler := NewEventHandler(eventService)
	eventHandler.Register(router)

	biographyService := service.NewBiographyService(pool)
	biographyHandler := NewBiographyHandler(biographyService)
	biographyHandler.Register(router)

	// TODO: Create and register contact details routes.

	// TODO: Create and register media (blob) routes.

	// TODO: Create and register authentication routes.

	return stack(router)
}
