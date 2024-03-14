package app

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/app/router"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/internal/svr"
	"github.com/asaskevich/EventBus"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Cfg        *config.Config
	Router     *router.Router
	Svr        *svr.Service
	Repository *repo.Repository
	Database   *mongo.Database
	Bus        *EventBus.EventBus
}

func New(
	cfg *config.Config,
	router *router.Router,
	svr *svr.Service,
	repository *repo.Repository,
	database *mongo.Database,
) *App {
	return &App{
		Cfg:        cfg,
		Svr:        svr,
		Router:     router,
		Repository: repository,
		Database:   database,
	}
}
