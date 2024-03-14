package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/svr"
	"cc.allio/fusion/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/wire"
	"os"
)

const BackupPathPrefix = "/api/admin/backup"

type BackupRouter struct {
	Cfg         *config.Config
	UserSvr     *svr.UserService
	MetaSvr     *svr.MetaService
	VisitSvr    *svr.VisitService
	ViewerSvr   *svr.ViewerService
	ArticleSvr  *svr.ArticleService
	CategorySvr *svr.CategoryService
	TagSvr      *svr.TagService
	DraftSvr    *svr.DraftService
	SettingSvr  *svr.SettingService
	StaticSvr   *svr.StaticService
}

var BackupRouterSet = wire.NewSet(wire.Struct(new(BackupRouter), "*"))

// ExportBackup
// @Summary 导出系统数据
// @Schemes
// @Description 导出系统数据
// @Tags Backup
// @Accept json
// @Produce json
// @Success 200 {object} nil
// @Router /api/admin/backup/export [Get]
func (router *BackupRouter) ExportBackup(c *gin.Context) *R {
	backup := &Backup{}
	{
		articles := router.ArticleSvr.GetAll("admin", true, false)
		backup.Articles = articles
	}
	{
		categories := router.CategorySvr.GetAllCategories()
		backup.Categories = categories
	}
	{
		tags := router.TagSvr.GetAllTags(true)
		backup.Tags = tags
	}
	{
		meta := router.MetaSvr.GetMeta()
		backup.Meta = meta
	}
	{
		drafts := router.DraftSvr.GetAll()
		backup.Drafts = drafts
	}
	{
		user := router.UserSvr.GetUser()
		backup.User = user
	}
	{
		viewers := router.ViewerSvr.GetAll()
		backup.Viewers = viewers
	}
	{
		visits := router.VisitSvr.GetAll()
		backup.Visits = visits
	}
	{
		statics := router.StaticSvr.GetAll()
		backup.Static = statics
	}
	{
		staticSetting := router.SettingSvr.FindStaticSetting()
		backup.Setting.Static = staticSetting
	}
	jsonBytes, err := json.Marshal(backup)
	if err != nil {
		return InternalError(err)
	}
	err = os.WriteFile("backup.json", jsonBytes, 0777)
	if err != nil {
		return InternalError(err)
	}
	c.FileAttachment("backup.json", "backup.json")
	// 删除文件
	defer func() {
		os.Remove("backup.json")
	}()
	return nil
}

// ImportBackup
// @Summary 导入系统数据
// @Schemes
// @Description 导入系统数据
// @Tags Backup
// @Accept json
// @Produce json
// @Param        file              formData      file        true       "file"
// @Success 200 {object} nil
// @Router /api/admin/backup/import [Post]
func (router *BackupRouter) ImportBackup(c *gin.Context) *R {
	filePart, err := c.FormFile("file")
	if err != nil {
		return InternalError(err)
	}
	file, err := filePart.Open()
	if err != nil {
		return InternalError(err)
	}
	body, err := web.ReadMultipartFile(file)
	if err != nil {
		return InternalError(err)
	}
	backup := &Backup{}
	err = json.Unmarshal(body, backup)
	if err != nil {
		return InternalError(err)
	}
	// TODO 需要对系统内每一个写入导入方法
	return nil
}

func (router *BackupRouter) Register(r *gin.Engine) {
	r.GET(BackupPathPrefix+"/export", Handle(router.ExportBackup))
	r.POST(BackupPathPrefix+"/import", Handle(router.ImportBackup))
}

type Backup struct {
	Articles   []*domain.Article  `json:"articles"`
	Categories []*domain.Category `json:"categories"`
	Tags       []string           `json:"tags"`
	Meta       *domain.Meta       `json:"meta"`
	Drafts     []*domain.Draft    `json:"drafts"`
	User       *domain.User       `json:"user"`
	Viewers    []*domain.Viewer   `json:"viewers"`
	Visits     []*domain.Visit    `json:"visits"`
	Static     []*domain.Static   `json:"static"`
	Setting    struct {
		Static *domain.StaticSetting `json:"static"`
	} `json:"setting"`
}
