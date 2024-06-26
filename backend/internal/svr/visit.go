package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/pkg/mongodb"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slog"
	"strconv"
)

type VisitService struct {
	Cfg            *config.Config
	VisitRepo      *repo.VisitRepository
	ArticleService *ArticleService
}

var VisitServiceSet = wire.NewSet(wire.Struct(new(VisitService), "*"))

func (v *VisitService) GetAll() []*domain.Visit {
	visits, err := v.VisitRepo.FindList(mongodb.NewLogical())
	if err != nil {
		slog.Error("Get visit all has error", "err", err)
		return make([]*domain.Visit, 0)
	}
	return visits
}

func (v *VisitService) GetLastVisitItem() *domain.Visit {
	limit := int64(1)
	opt := options.FindOptions{Sort: bson.E{Key: "lastVisitedTime", Value: -1}, Limit: &limit}
	visits, err := v.VisitRepo.FindList(mongodb.NewLogicalDefaultArray(bson.D{{"lastVisitedTime", bson.D{{"$exists", true}}}}), &opt)
	if err != nil {
		slog.Error("Find visited has err", "err", err.Error())
		return nil
	}
	if visits != nil && len(visits) > 0 {
		return visits[0]
	}
	return nil
}

func (v *VisitService) GetByArticleIdOrPathname(idOrPathname string) *domain.Visit {
	id, err := strconv.Atoi(idOrPathname)
	var pathname string
	if err != nil || id == 0 {
		pathname = "/about"
	} else if err == nil {
		pathname = "/post/" + string(rune(id))
	}
	visit, err := v.VisitRepo.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "pathname", Value: pathname}))
	if err != nil {
		return nil
	}
	return visit
}
