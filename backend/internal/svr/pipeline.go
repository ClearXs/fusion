package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/pkg/mongodb"
	"cc.allio/fusion/pkg/util"
	"errors"
	"github.com/google/wire"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slog"
)

type PipelineService struct {
	Cfg                *config.Config
	PipelineRepository *repo.PipelineRepository
}

var events = []domain.EventKey{
	domain.LoginEvent,
	domain.LogoutEvent,
	domain.BeforeUpdateArticleEvent,
	domain.AfterUpdateArticleEvent,
	domain.DeleteArticleEvent,
	domain.BeforeUpdateDraftEvent,
	domain.AfterUpdateDraftEvent,
	domain.DeleteDraftEvent,
	domain.UpdateSiteInfoEvent,
	domain.ManualTriggerEvent,
}

var PipelineServiceSet = wire.NewSet(wire.Struct(new(PipelineService), "*"))

func (p *PipelineService) GetAll() []*domain.Pipeline {
	pipelines, err := p.PipelineRepository.FindList(mongodb.NewLogicalDefaultLogical(mongodb.NewLogicalDefaultArray(DeleteFilter)))
	if err != nil {
		slog.Error("get all pipelines has error", "err", err)
		return make([]*domain.Pipeline, 0)
	}
	return pipelines
}

func (p *PipelineService) GetPipelineById(id int64) (*domain.Pipeline, error) {
	pipeline, err := p.PipelineRepository.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "id", Value: id}))
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

func (p *PipelineService) GetPipelineByEventKey(eventKey domain.EventKey) []*domain.Pipeline {
	filter := mongodb.NewLogical()
	filter.Append(bson.E{Key: "eventName", Value: eventKey})
	filter.AppendLogical(mongodb.NewLogicalDefaultArray(DeleteFilter))
	pipelines, err := p.PipelineRepository.FindList(filter)
	if err != nil {
		slog.Error("get pipelines by event key has error", "err", err, "eventKey", eventKey)
		return make([]*domain.Pipeline, 0)
	}
	return pipelines
}

func (p *PipelineService) CreatePipeline(pipeline *domain.Pipeline) (bool, error) {
	if !lo.Contains(events, pipeline.EventName) {
		return false, errors.New("the event not found in system support event")
	}
	saved, err := p.PipelineRepository.Save(pipeline)
	if err != nil {
		return false, err
	}
	return saved > 0, nil
}

func (p *PipelineService) UpdatePipeline(pipeline *domain.Pipeline) (bool, error) {
	data := util.ToBsonElements(pipeline)
	return p.PipelineRepository.Update(mongodb.NewLogicalDefault(bson.E{Key: "id", Value: pipeline.Id}), data)
}

func (p *PipelineService) DeletePipeline(id int64) (bool, error) {
	return p.PipelineRepository.Update(mongodb.NewLogicalDefault(bson.E{Key: "id", Value: id}), bson.D{{"deleted", true}})
}
