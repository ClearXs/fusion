package app

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/app/router"
	"cc.allio/fusion/internal/event"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/internal/svr"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Cfg        *config.Config
	Router     *router.Router
	Svr        *svr.Service
	Repository *repo.Repository
	Database   *mongo.Database
	Isr        *event.IsrEventBus
}

func New(
	cfg *config.Config,
	router *router.Router,
	svr *svr.Service,
	repository *repo.Repository,
	database *mongo.Database,
	isr *event.IsrEventBus,
) *App {
	return &App{
		Cfg:        cfg,
		Svr:        svr,
		Router:     router,
		Repository: repository,
		Database:   database,
		Isr:        isr,
	}
}
