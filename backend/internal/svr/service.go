package svr

import (
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	DeleteFilter = bson.D{{"deleted", false}, {"deleted", primitive.E{Key: "$exists", Value: false}}}
	HiddenFilter = bson.D{{"hidden", false}, {"hidden", primitive.E{Key: "$exists", Value: false}}}
)

type Service struct {
	UserService       *UserService
	TokenService      *TokenService
	AuthService       *AuthService
	MetaService       *MetaService
	VisitService      *VisitService
	ViewerService     *ViewerService
	ArticleService    *ArticleService
	CategoryService   *CategoryService
	TagService        *TagService
	AnalysisService   *AnalysisService
	DraftService      *DraftService
	SettingsService   *SettingService
	StaticService     *StaticService
	CaddyService      *CaddyService
	CustomPageService *CustomPageService
	PipelineService   *PipelineService
}

var ServiceSet = wire.NewSet(
	AnalysisServiceSet,
	ArticleServiceSet,
	AuthServiceSet,
	CategoryServiceSet,
	DraftServiceSet,
	MetaServiceSet,
	SettingServiceSet,
	StaticServiceSet,
	TagServiceSet,
	TokenServiceSet,
	UserServiceSet,
	ViewerServiceSet,
	VisitServiceSet,
	CaddyServiceSet,
	CustomPageServiceSet,
	PipelineServiceSet,
	wire.Struct(new(Service), "*"),
)
