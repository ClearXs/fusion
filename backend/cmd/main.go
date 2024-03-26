package main

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/app"
	"cc.allio/fusion/pkg/env"
	"cc.allio/fusion/pkg/logger"
	"context"
	"fmt"
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

	// setup app
	a, cleanup, err := app.InitApp(ctx, cfg)
	if err != nil {
		slog.Error("init app err", "err", err)
		cleanup()
		cancel()
		panic(err)
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
	standardLogger := logrus.StandardLogger()
	standardLogger.SetOutput(a.Logger)
	log := slog.New(logger.NewLogrusHandler(standardLogger))
	slog.SetDefault(log)

	// init gin
	r := gin.New()

	// init router
	a.Router.Init(r)

	a.Isr.ActiveAll("trigger incremental rendering by startup")

	// startup gin server
	addr := ":" + strconv.Itoa(cfg.Server.Port)
	fmt.Println(`
    ____ _  _ ____ _ ____ _  _
    |___ |  | [__  | |  | |\ |
    |    |__| ___] | |__| | \|

    God’s Blessing on This Wonderful World!
    version: ` + env.Version + `     address: ` + addr + `
`)
	err = r.Run(addr)
	if err != nil {
		slog.Error("start server failed", "err", err)
		cleanup()
		cancel()
		panic(err)
	}
}
