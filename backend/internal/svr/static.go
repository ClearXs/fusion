package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/pkg/mongodb"
	"cc.allio/fusion/pkg/storage"
	"context"
	"github.com/google/wire"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slog"
)

type StaticService struct {
	Cfg         *config.Config
	FileService *FileService
	StaticRepo  *repo.StaticRepository
}

var StaticServiceSet = wire.NewSet(wire.Struct(new(StaticService), "*"))

func (s *StaticService) CreateStatic(ctx context.Context, header *storage.FileHeader) (*credential.FilePathCredential, error) {
	static, err := s.FileService.Upload(ctx, header)
	if err != nil {
		return nil, err
	}

	_, err = s.StaticRepo.Save(static)
	if err != nil {
		return nil, err
	}
	return &credential.FilePathCredential{Src: static.RealPath, IsNew: true}, nil
}

func (s *StaticService) GetAll(storageType domain.StaticType) []*domain.Static {
	filter := mongodb.NewLogical()
	if storageType != "" {
		filter.Append(bson.E{Key: "staticType", Value: storageType})
	}
	statics, err := s.StaticRepo.FindList(filter)
	if err != nil {
		slog.Error("Get static all has error", "err", err)
		return make([]*domain.Static, 0)
	}
	return statics
}

func (s *StaticService) GetByOption(option *credential.StaticSearchOption) *domain.StaticPageResult {
	filter := mongodb.NewLogical()
	if lo.IsNotEmpty(option.StaticType) {
		filter.Append(bson.E{Key: "staticType", Value: option.StaticType})
	}
	opt := options.FindOptions{}
	if option.PageSize != -1 {
		skip := int64(option.PageSize*option.Page - option.PageSize)
		limit := int64(option.PageSize)
		opt.Skip = &skip
		opt.Limit = &limit
	}

	statics, err := s.StaticRepo.FindList(filter, &opt)
	if err != nil {
		slog.Error("Failed to static page list by option", "err", err, "option", option)
		return &domain.StaticPageResult{}
	}
	count, err := s.StaticRepo.Count(filter)
	if err != nil {
		slog.Error("Failed to count to static page list", "err", err, "option", option)
	}

	return &domain.StaticPageResult{Data: statics, Total: count}
}

func (s *StaticService) DeleteBySign(ctx context.Context, sign string) (bool, error) {
	filter := mongodb.NewLogicalDefault(bson.E{Key: "sign", Value: sign})
	static, err := s.StaticRepo.FindOne(filter)
	if err != nil {
		return false, err
	}
	// remove file in filesystem
	err = s.FileService.Delete(ctx, []string{static.RealPath})
	if err != nil {
		return false, err
	}
	// remove file in db
	return s.StaticRepo.Remove(filter)
}

// DeleteAll delete file in filesystem and db
func (s *StaticService) DeleteAll(ctx context.Context, staticType domain.StaticType) (bool, error) {
	statics := s.GetAll(staticType)

	// remove file in filesystem
	filepath := lo.Map[*domain.Static, string](statics, func(item *domain.Static, index int) string { return item.RealPath })
	err := s.FileService.Delete(ctx, filepath)
	if err != nil {
		return false, err
	}

	// remove file in db
	return s.StaticRepo.RemoveMany(mongodb.NewLogicalDefault(bson.E{Key: "staticType", Value: staticType}))
}
