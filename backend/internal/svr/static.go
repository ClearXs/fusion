package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/pkg/mongodb"
	"github.com/google/wire"
	"golang.org/x/exp/slog"
)

type StaticService struct {
	Cfg        *config.Config
	StaticRepo *repo.StaticRepository
}

var StaticServiceSet = wire.NewSet(wire.Struct(new(StaticService), "*"))

func (s *StaticService) GetAll() []*domain.Static {
	statics, err := s.StaticRepo.FindList(mongodb.NewLogical())
	if err != nil {
		slog.Error("Get static all has error", "err", err)
		return make([]*domain.Static, 0)
	}
	return statics
}
