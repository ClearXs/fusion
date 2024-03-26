package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/apm"
	"cc.allio/fusion/internal/credential"
	"fmt"
	"github.com/google/wire"
)

type BusinessType = string

const (
	LoginLogType       = "login"
	LogoutLogType      = "logout"
	RunPipelineLogType = "runPipeline"
)

type LogService struct {
	Cfg    *config.Config
	Logger *apm.Logger
}

var LogServiceSet = wire.NewSet(wire.Struct(new(LogService), "*"))

func (l *LogService) GetByOption(option *credential.LogSearchOption) (*apm.QueryResult, error) {
	// build search option
	cfg := l.Cfg
	loggerConfig := cfg.Log.Apm
	logger := l.Logger
	search := apm.NewSearch()
	search.Page(option.Page, option.PageSize)
	sql := fmt.Sprintf("select * from %s", loggerConfig.Stream)
	if option.LogType != "" {
		sql = sql + fmt.Sprintf("where logType = '%s'", option.LogType)
	}
	search.SQL(sql)

	// find log
	formatter, err := logger.Find(search)
	if err != nil {
		return nil, err
	}

	// format
	result, err := formatter.Format()
	if err != nil {
		return nil, err
	}
	return &result, nil
}
