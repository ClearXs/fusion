package router

import (
	"bufio"
	"bytes"
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/svr"
	"cc.allio/fusion/pkg/img"
	"cc.allio/fusion/pkg/storage"
	"cc.allio/fusion/pkg/web"
	"context"
	"errors"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"io"
	"mime/multipart"
	"path/filepath"
)

const imagePathPrefix = "/api/admin/img"

type ImgRoute struct {
	Cfg            *config.Config
	StaticService  *svr.StaticService
	SettingService *svr.SettingService
}

var ImgRouterSet = wire.NewSet(wire.Struct(new(ImgRoute), "*"))

// Upload
// @Summary Upload image
// @Schemes
// @Description Upload image
// @Tags Static
// @Accept multipart/form-data
// @Produce json
// @Param       file              formData     file       true     "file"
// @Param       waterMarkText     formData     string     true     "waterMarkText"
// @Param       withWaterMark     formData     bool       true     "withWaterMark"
// @Success 200 {object} bool
// @Router /api/admin/img/upload [POST]
func (i *ImgRoute) Upload(c *gin.Context) *R {
	// create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	header, err := c.FormFile("file")
	if err != nil {
		return InternalError(err)
	}
	file, err := header.Open()
	if err != nil {
		return InternalError(err)
	}

	filename := header.Filename
	size := header.Size
	ext := filepath.Ext(filename)
	// build upload file header
	fileHeader := &storage.FileHeader{
		Filename: filename,
		FilePath: "/" + filename,
		Ext:      ext,
		Size:     uint64(size),
		File:     file,
		Header:   header.Header,
	}

	// determine whether file going on append watermark by file extension and with withWaterMark
	withWaterMark := web.ParseBoolForPath(c, "withWaterMark")
	if format, err := imaging.FormatFromFilename(filename); err == nil && withWaterMark {
		waterMarkText := c.Param("waterMarkText")
		// append watermark
		newFile, err := i.addWatermark(file, waterMarkText, format)
		if err != nil {
			return InternalError(err)
		}
		fileHeader.File = io.NopCloser(newFile)
	}

	// upload
	static, err := i.StaticService.CreateStatic(ctx, fileHeader)
	if err != nil {
		return InternalError(err)
	}
	return Ok(static)
}

// GetAll
// @Summary get all images
// @Schemes
// @Description get all images
// @Tags Static
// @Accept json
// @Produce json
// @Success 200 {object} []domain.Static
// @Router /api/admin/img/all [Get]
func (i *ImgRoute) GetAll(c *gin.Context) *R {
	statics := i.StaticService.GetAll(domain.ImgStaticType)
	return Ok(statics)
}

// DeleteAll
// @Summary delete all images
// @Schemes
// @Description delete all images
// @Tags Static
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /api/admin/img/all/delete [Get]
func (i *ImgRoute) DeleteAll(c *gin.Context) *R {
	if i.Cfg.Demo {
		return Error(401, errors.New("演示站禁止操作！！"))
	}
	// create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	success, err := i.StaticService.DeleteAll(ctx, domain.ImgStaticType)
	if err != nil {
		return InternalError(err)
	}
	return Ok(success)
}

// DeleteBySign
// @Summary delete image by sign
// @Schemes
// @Description delete image by sign
// @Tags Static
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /api/admin/img/:sign [Get]
func (i *ImgRoute) DeleteBySign(c *gin.Context) *R {
	if i.Cfg.Demo {
		return Error(401, errors.New("演示站禁止操作！！"))
	}
	// create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sign := c.Param("sign")
	success, err := i.StaticService.DeleteBySign(ctx, sign)
	if err != nil {
		return InternalError(err)
	}
	return Ok(success)
}

// GetByOption
// @Summary get images by option
// @Schemes
// @Description get images by option
// @Tags Static
// @Accept json
// @Produce json
// @Param        page                 query      int        false       "page"               -1
// @Param        pageSize             query      int        false       "pageSize"            5
// @Param        staticType           query      int        false       "staticType"          "img"
// @Success 200 {object} domain.StaticPageResult
// @Router /api/admin/img [Get]
func (i *ImgRoute) GetByOption(c *gin.Context) *R {
	option := &credential.StaticSearchOption{}
	if err := c.Bind(option); err != nil {
		return InternalError(err)
	}
	statics := i.StaticService.GetByOption(option)
	return Ok(statics)
}

// addWatermark by file and text
func (i *ImgRoute) addWatermark(file multipart.File, text string, format imaging.Format) (io.Reader, error) {
	watermark, err := img.NewWatermark(file)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}

	err = watermark.Encode(buf, text, img.Center, format)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(buf)
	return reader, nil
}

func (i *ImgRoute) Register(r *gin.Engine) {
	r.POST(imagePathPrefix+"/upload", Handle(i.Upload))
	r.GET(imagePathPrefix+"/all", Handle(i.GetAll))
	r.DELETE(imagePathPrefix+"/all/delete", Handle(i.DeleteAll))
	r.DELETE(imagePathPrefix+"/:sign", Handle(i.DeleteBySign))
	r.GET(imagePathPrefix, Handle(i.GetByOption))
}
