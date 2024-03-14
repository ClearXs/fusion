package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/pkg/mongodb"
	"cc.allio/fusion/pkg/utils"
	"errors"
	"github.com/google/wire"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slog"
	"strings"
	"time"
)

type DraftService struct {
	Cfg        *config.Config
	DraftRepo  *repo.DraftRepository
	ArticleSvr *ArticleService
}

var DraftServiceSet = wire.NewSet(wire.Struct(new(DraftService), "*"))

func (d *DraftService) GetAll() []*domain.Draft {
	drafts, err := d.DraftRepo.FindList(mongodb.NewLogical())
	if err != nil {
		slog.Error("Get draft all has err", "err", err)
		return make([]*domain.Draft, 0)
	}
	return drafts
}

func (d *DraftService) GetByOption(option *credential.DraftSearchOptionCredential) *domain.DraftPageResult {
	filter := mongodb.NewLogical()
	// delete
	filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(DeleteFilter))
	opt := options.FindOptions{}
	// tags
	if lo.IsNotEmpty(option.Tags) {
		tags := strings.Split(option.Tags, ",")
		tagFilters := make([]bson.E, 0)
		for _, tag := range tags {
			tagFilters = append(tagFilters, bson.E{Key: "tags", Value: bson.D{{"$regex", tag}, {"$options", "i"}}})
		}
		if len(tagFilters) > 0 {
			filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(tagFilters))
		}
	}
	// category
	if lo.IsNotEmpty(option.Category) {
		filter.Append(bson.E{Key: "category", Value: bson.D{{"$regex", option.Category}, {"$options", "i"}}})
	}

	// title
	if lo.IsNotEmpty(option.Title) {
		filter.Append(bson.E{Key: "title", Value: bson.D{{"$regex", option.Title}, {"$options", "i"}}})
	}

	// time
	if lo.IsNotEmpty(option.StartTime) || lo.IsNotEmpty(option.EndTime) {
		timeFilter := make([]bson.E, 0)
		if lo.IsNotEmpty(option.StartTime) {
			startTime, err := time.Parse(time.DateTime, option.StartTime)
			if err == nil {
				timeFilter = append(timeFilter, bson.E{Key: "$gte", Value: startTime})
			}
		}
		if lo.IsNotEmpty(option.EndTime) {
			endTime, err := time.Parse(time.DateTime, option.EndTime)
			if err == nil {
				timeFilter = append(timeFilter, bson.E{Key: "$lte", Value: endTime})
			}
		}
		if len(timeFilter) > 0 {
			filter.Append(bson.E{Key: "createdAt", Value: timeFilter})
		}
	}
	if option.PageSize != -1 {
		skip := int64(option.PageSize*option.Page - option.PageSize)
		limit := int64(option.PageSize)
		opt.Skip = &skip
		opt.Limit = &limit
	}

	drafts, err := d.DraftRepo.FindList(filter, &opt)
	if err != nil {
		slog.Error("find drafts by options has error", "err", err, "options", option)
		return &domain.DraftPageResult{}
	}

	count, err := d.DraftRepo.Count(filter)
	if err != nil {
		slog.Error("count for draft by options has error", "err", err, "options", option)
	}
	return &domain.DraftPageResult{Drafts: drafts, Total: count}
}

func (d *DraftService) GetById(id int64) (*domain.Draft, error) {
	filter := mongodb.NewLogical()
	// delete
	filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(DeleteFilter))
	filter.Append(bson.E{Key: "id", Value: id})
	return d.DraftRepo.FindOne(filter)
}

func (d *DraftService) DeleteById(id int64) (bool, error) {
	return d.DraftRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "id", Value: id}), bson.D{{"deleted", true}})
}

func (d *DraftService) UpdateById(id int64, draft *domain.Draft) (bool, error) {
	update := utils.ToBsonElements(draft)
	return d.DraftRepo.Update(mongodb.NewLogicalOrDefault(bson.E{Key: "id", Value: id}), update)
}

func (d *DraftService) Create(draft *domain.Draft) (*domain.Draft, error) {
	nextId, err := d.getNextId()
	if err != nil {
		return nil, err
	}
	draft.Id = nextId
	saved, err := d.DraftRepo.Save(draft)
	if err != nil {
		return nil, err
	}
	if saved < 0 {
		return nil, errors.New("create draft error")
	}
	return draft, nil
}

func (d *DraftService) Publish(id int64, option *credential.DraftPublishCredential) (*domain.Article, error) {
	draft, err := d.GetById(id)
	if err != nil {
		return nil, err
	}
	article := &domain.Article{
		Title:     draft.Title,
		Content:   draft.Content,
		Tags:      draft.Tags,
		Category:  draft.Category,
		Author:    draft.Author,
		Hidden:    option.Hidden,
		Pathname:  option.Pathname,
		Private:   option.Private,
		Password:  option.Password,
		Copyright: option.Copyright,
	}
	return d.ArticleSvr.Create(article)
}

func (d *DraftService) getNextId() (int64, error) {
	result, err := d.DraftRepo.FindOne(mongodb.NewLogical(), &options.FindOneOptions{Sort: bson.E{Key: "id", Value: -1}})
	if err != nil {
		return -1, err
	}
	return result.Id + 1, nil
}
