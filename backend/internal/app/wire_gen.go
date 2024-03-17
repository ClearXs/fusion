// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/app/router"
	"cc.allio/fusion/internal/event"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/internal/svr"
	"cc.allio/fusion/pkg/mongodb"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// Injectors from wire.go:

func InitApp(ctx context.Context, cfg *config.Config) (*App, func(), error) {
	database, cleanup, err := mongodbConnect(cfg)
	if err != nil {
		return nil, nil, err
	}
	metaRepository := &repo.MetaRepository{
		Cfg: cfg,
		Db:  database,
	}
	viewerRepository := &repo.ViewerRepository{
		Cfg: cfg,
		Db:  database,
	}
	viewerService := &svr.ViewerService{
		Cfg:        cfg,
		ViewerRepo: viewerRepository,
	}
	visitRepository := &repo.VisitRepository{
		Cfg: cfg,
		Db:  database,
	}
	articleRepository := &repo.ArticleRepository{
		Cfg: cfg,
		Db:  database,
	}
	metaService := &svr.MetaService{
		Cfg:           cfg,
		MetaRepo:      metaRepository,
		ViewerService: viewerService,
		VisitRepo:     visitRepository,
		ArticleRepo:   articleRepository,
	}
	userRepository := &repo.UserRepository{
		Cfg: cfg,
		Db:  database,
	}
	userService := &svr.UserService{
		Cfg:            cfg,
		UserRepository: userRepository,
	}
	settingsRepository := &repo.SettingsRepository{
		Cfg: cfg,
		Db:  database,
	}
	settingService := &svr.SettingService{
		Cfg:         cfg,
		SettingRepo: settingsRepository,
	}
	tokenRepository := &repo.TokenRepository{
		Cfg: cfg,
		Db:  database,
	}
	tokenService := &svr.TokenService{
		Cfg:         cfg,
		SettingsSvr: settingService,
		TokenRepo:   tokenRepository,
	}
	authService := &svr.AuthService{
		Cfg:          cfg,
		UserService:  userService,
		TokenService: tokenService,
	}
	categoryRepository := &repo.CategoryRepository{
		Cfg: cfg,
		Db:  database,
	}
	articleService := &svr.ArticleService{
		Cfg:          cfg,
		MetaService:  metaService,
		ArticleRepo:  articleRepository,
		CategoryRepo: categoryRepository,
	}
	visitService := &svr.VisitService{
		Cfg:            cfg,
		VisitRepo:      visitRepository,
		ArticleService: articleService,
	}
	categoryService := &svr.CategoryService{
		Cfg:          cfg,
		CategoryRepo: categoryRepository,
		ArticleSvr:   articleService,
	}
	tagService := &svr.TagService{
		Cfg:        cfg,
		ArticleSvr: articleService,
	}
	analysisService := &svr.AnalysisService{
		Cfg:         cfg,
		MetaSvr:     metaService,
		ArticleSvr:  articleService,
		ViewerSvr:   viewerService,
		VisitSvr:    visitService,
		TagSvr:      tagService,
		CategorySvr: categoryService,
	}
	draftRepository := &repo.DraftRepository{
		Cfg: cfg,
		Db:  database,
	}
	draftService := &svr.DraftService{
		Cfg:        cfg,
		DraftRepo:  draftRepository,
		ArticleSvr: articleService,
	}
	staticRepository := &repo.StaticRepository{
		Cfg: cfg,
		Db:  database,
	}
	staticService := &svr.StaticService{
		Cfg:        cfg,
		StaticRepo: staticRepository,
	}
	svrSettingService := svr.SettingService{
		Cfg:         cfg,
		SettingRepo: settingsRepository,
	}
	caddyService := &svr.CaddyService{
		Cfg:             cfg,
		SettingsService: svrSettingService,
	}
	customPageRepository := &repo.CustomPageRepository{
		Cfg: cfg,
		Db:  database,
	}
	customPageService := &svr.CustomPageService{
		Cfg:                  cfg,
		CustomPageRepository: customPageRepository,
	}
	service := &svr.Service{
		UserService:       userService,
		TokenService:      tokenService,
		AuthService:       authService,
		MetaService:       metaService,
		VisitService:      visitService,
		ViewerService:     viewerService,
		ArticleService:    articleService,
		CategoryService:   categoryService,
		TagService:        tagService,
		AnalysisService:   analysisService,
		DraftService:      draftService,
		SettingsService:   settingService,
		StaticService:     staticService,
		CaddyService:      caddyService,
		CustomPageService: customPageService,
	}
	isrEventBus := event.NewIsrEventBus(service)
	aboutRouter := &router.AboutRouter{
		Cfg:     cfg,
		MetaSvr: metaService,
		Isr:     isrEventBus,
	}
	analysisRouter := &router.AnalysisRouter{
		Cfg:         cfg,
		AnalysisSvr: analysisService,
	}
	articleRouter := &router.ArticleRouter{
		Cfg:        cfg,
		ArticleSvr: articleService,
		Isr:        isrEventBus,
	}
	authRouter := &router.AuthRouter{
		Cfg:      cfg,
		AuthSvr:  authService,
		UserSvr:  userService,
		TokenSvr: tokenService,
	}
	backupRouter := &router.BackupRouter{
		Cfg:         cfg,
		UserSvr:     userService,
		MetaSvr:     metaService,
		VisitSvr:    visitService,
		ViewerSvr:   viewerService,
		ArticleSvr:  articleService,
		CategorySvr: categoryService,
		TagSvr:      tagService,
		DraftSvr:    draftService,
		SettingSvr:  settingService,
		StaticSvr:   staticService,
	}
	userRouter := &router.UserRouter{
		Cfg:     cfg,
		UserSvr: userService,
	}
	caddyRouter := &router.CaddyRouter{
		Cfg:             cfg,
		SettingsService: settingService,
		CaddyService:    caddyService,
	}
	categoryRouter := &router.CategoryRouter{
		Cfg:             cfg,
		CategoryService: categoryService,
		Isr:             isrEventBus,
	}
	collaboratorRouter := &router.CollaboratorRouter{
		Cfg:          cfg,
		UserService:  userService,
		MetaService:  metaService,
		TokenService: tokenService,
	}
	customPageRouter := &router.CustomPageRouter{
		Cfg:           cfg,
		StaticService: staticService,
	}
	draftRouter := &router.DraftRouter{
		Cfg:          cfg,
		DraftService: draftService,
		Isr:          isrEventBus,
	}
	linkRouter := &router.LinkRouter{
		Cfg:         cfg,
		MetaService: metaService,
	}
	menuRouter := &router.MenuRouter{
		Cfg:             cfg,
		SettingsService: settingService,
	}
	metaRouter := &router.MetaRouter{
		Cfg:         cfg,
		MetaService: metaService,
	}
	rewardRouter := &router.RewardRouter{
		Cfg:         cfg,
		MetaService: metaService,
	}
	settingRouter := &router.SettingRouter{
		Cfg:            cfg,
		SettingService: settingService,
	}
	siteRouter := &router.SiteRouter{
		Cfg:         cfg,
		MetaService: metaService,
	}
	socialRouter := &router.SocialRouter{
		Cfg:         cfg,
		MetaService: metaService,
	}
	tagRouter := &router.TagRouter{
		Cfg:        cfg,
		TagService: tagService,
	}
	tokenRouter := &router.TokenRouter{
		Cfg:          cfg,
		TokenService: tokenService,
	}
	publicRouter := &router.PublicRouter{
		Cfg:               cfg,
		ArticleService:    articleService,
		TagService:        tagService,
		MetaService:       metaService,
		ViewerService:     viewerService,
		VisitService:      visitService,
		SettingService:    settingService,
		CustomPageService: customPageService,
		CategoryService:   categoryService,
	}
	routerRouter := &router.Router{
		AboutRouter:        aboutRouter,
		AnalysisRouter:     analysisRouter,
		ArticleRouter:      articleRouter,
		AuthRouter:         authRouter,
		BackupRouter:       backupRouter,
		UserRouter:         userRouter,
		CaddyRouter:        caddyRouter,
		CategoryRouter:     categoryRouter,
		CollaboratorRouter: collaboratorRouter,
		CustomPageRouter:   customPageRouter,
		DraftRouter:        draftRouter,
		LinkRouter:         linkRouter,
		MenuRouter:         menuRouter,
		MetaRouter:         metaRouter,
		RewardRouter:       rewardRouter,
		SettingRouter:      settingRouter,
		SiteRouter:         siteRouter,
		SocialRouter:       socialRouter,
		TagRouter:          tagRouter,
		TokenRouter:        tokenRouter,
		PublicRouter:       publicRouter,
	}
	pipelineRepository := &repo.PipelineRepository{
		Cfg: cfg,
		Db:  database,
	}
	repository := &repo.Repository{
		ArticleRepository:    articleRepository,
		CategoryRepository:   categoryRepository,
		DraftRepository:      draftRepository,
		MetaRepository:       metaRepository,
		SettingsRepository:   settingsRepository,
		StaticRepository:     staticRepository,
		TokenRepository:      tokenRepository,
		UserRepository:       userRepository,
		ViewerRepository:     viewerRepository,
		VisitRepository:      visitRepository,
		CustomPageRepository: customPageRepository,
		PipelineRepository:   pipelineRepository,
	}
	app := New(cfg, routerRouter, service, repository, database, isrEventBus)
	return app, func() {
		cleanup()
	}, nil
}

// wire.go:

func mongodbConnect(cfg *config.Config) (*mongo.Database, func(), error) {
	return mongodb.Connect(&cfg.Mongodb)
}
