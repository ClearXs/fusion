package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/pkg/mongodb"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slog"
)

type CustomPageService struct {
	Cfg                  *config.Config
	CustomPageRepository *repo.CustomPageRepository
}

var CustomPageServiceSet = wire.NewSet(wire.Struct(new(CustomPageService), "*"))

func (c *CustomPageService) GetAll() []*domain.CustomPage {
	customPages, err := c.CustomPageRepository.FindList(mongodb.NewLogicalDefault(bson.E{Key: "html", Value: 0}))
	if err != nil {
		slog.Error("get all custom page has error", "err", err)
		return make([]*domain.CustomPage, 0)
	}
	return customPages
}

func (c *CustomPageService) GetByPath(path string) *domain.CustomPage {
	customPage, err := c.CustomPageRepository.FindOne(mongodb.NewLogicalDefault(bson.E{"path", path}))
	if err != nil {
		slog.Error("get custom page by path has error", "err", err, "path", path)
		return nil
	}
	return customPage
}
