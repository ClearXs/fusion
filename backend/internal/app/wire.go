//go:build wireinject
// +build wireinject

package app

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/apm"
	"cc.allio/fusion/internal/app/router"
	"cc.allio/fusion/internal/event"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/internal/svr"
	"cc.allio/fusion/pkg/mongodb"
	"context"
	"github.com/asaskevich/EventBus"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitApp(ctx context.Context, cfg *config.Config) (*App, func(), error) {
	panic(wire.Build(
		New,
		mongodbConnect,
		createEventBus(),
		apm.LoggerSet,
		repo.RepositorySet,
		svr.ServiceSet,
		router.Set,
		event.IsrEventBusSet,
		event.ScriptEngineSet,
	))
}

func mongodbConnect(cfg *config.Config) (*mongo.Database, func(), error) {
	return mongodb.Connect(&cfg.Mongodb)
}

func createEventBus() EventBus.Bus {
	return &EventBus.New()
}
