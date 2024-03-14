/*
 *
 *  ____ _  _ ____ _ ____ _  _
 * |___ |  | [__  | |  | |\ |
 * |    |__| ___] | |__| | \|
 *
 * generate by https://patorjk.com/software/taag/#p=testall&f=Standard&t=Fusion
 */
package main

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/app"
	"cc.allio/fusion/pkg/logger"
	"cc.allio/fusion/pkg/middleware"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slog"
	"os"
	"strconv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error("failed get config by config.yml", "err", err)
	}
	// init logger
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		slog.Error("parse log level error", "system log level", cfg.Log.Level, "err", err)
	} else {
		logrus.SetLevel(level)
	}
	log := slog.New(logger.NewLogrusHandler(logrus.StandardLogger()))
	slog.SetDefault(log)

	// setup app
	a, cleanup, err := app.InitApp(ctx, cfg)
	if err != nil {
		slog.Error("init app err", "err", err)
		cleanup()
		cancel()
		panic(err)
	}

	// init gin
	r := gin.New()
	// middleware setup
	r.Use(gin.Recovery())
	r.Use(middleware.Logging())

	// init router
	a.Router.Init(r)

	// startup gin server
	addr := ":" + strconv.Itoa(cfg.Server.Port)
	slog.Info("\U0001FAE7 start server...", "address", addr)
	err = r.Run(addr)
	if err != nil {
		slog.Error("start server failed", "err", err)
		cleanup()
		cancel()
		panic(err)
	}
}
