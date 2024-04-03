package svr

import (
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	DeleteFilter = bson.D{{"deleted", false}}
	HiddenFilter = bson.D{{"hidden", false}}
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
	FileService       *FileService
	LogService        *LogService
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
	FileServiceSet,
	LogServiceSet,
	wire.Struct(new(Service), "*"),
)
