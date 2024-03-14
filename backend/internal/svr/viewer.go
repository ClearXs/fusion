package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/pkg/mongodb"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slog"
	"time"
)

type ViewerService struct {
	Cfg        *config.Config
	ViewerRepo *repo.ViewerRepository
}

var ViewerServiceSet = wire.NewSet(wire.Struct(new(ViewerService), "*"))

func (v *ViewerService) GetAll() []*domain.Viewer {
	viewers, err := v.ViewerRepo.FindList(mongodb.NewLogical())
	if err != nil {
		slog.Error("Get viewer all has error", "err", err)
		return make([]*domain.Viewer, 0)
	}
	return viewers
}

func (v *ViewerService) GetViewerGrid(num int64) *domain.ViewerGrid {
	curDate := time.Now()
	gridTotal := make([]domain.DateViewer, 0)
	temArr := make([]domain.DateViewer, 0)
	today := domain.DataViewer{Viewer: 0, Visited: 0}
	lastDay := domain.DataViewer{Viewer: 0, Visited: 0}
	for i := 0; i < int(num); i++ {
		last := curDate.AddDate(0, 0, -1*i).Format(time.DateOnly)
		lastDayData, err := v.ViewerRepo.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "date", Value: last}))
		if err != nil {
			panic(err)
		}
		if i == 0 && lastDayData != nil {
			today.Viewer = lastDayData.Viewer
			today.Visited = lastDayData.Visited
		}

		if i == 1 {
			if lastDayData != nil {
				lastDay.Viewer = lastDayData.Viewer
				lastDay.Visited = lastDayData.Visited
			}
			if today.Viewer == 0 {
				// 如果今天没数据，那今天的就和昨天的一样吧。这样新增就都是 0
				today.Viewer = lastDayData.Viewer
				today.Visited = lastDayData.Visited
			}
		}

		if lastDayData != nil {
			temArr = append(temArr, domain.DateViewer{Date: last, Viewer: lastDayData.Viewer, Visited: lastDayData.Visited})
			if i != int(num)+1 {
				gridTotal = append(gridTotal, domain.DateViewer{Date: last, Viewer: lastDayData.Viewer, Visited: lastDayData.Visited})
			}
		}
	}

	gridEachDay := make([]domain.DateViewer, 0)
	pre := &temArr[0]
	for i := 1; i < len(temArr); i++ {
		curObj := &temArr[i]
		if curObj != nil {
			if pre != nil {
				gridEachDay = append(gridEachDay, domain.DateViewer{Date: curObj.Date, Viewer: curObj.Viewer - pre.Viewer, Visited: curObj.Visited - pre.Visited})
			} else {
				gridEachDay = append(gridEachDay, domain.DateViewer{Date: curObj.Date, Viewer: curObj.Viewer, Visited: curObj.Visited})
			}
		}
		pre = curObj
	}

	return &domain.ViewerGrid{
		Grid: domain.GridDateView{
			Total: gridTotal,
			Each:  gridEachDay,
		},
		Add: domain.DataViewer{
			Viewer:  today.Viewer - lastDay.Viewer,
			Visited: today.Viewer - lastDay.Viewer,
		},
		Now: domain.DataViewer{
			Viewer:  today.Viewer,
			Visited: today.Visited,
		},
	}
}
