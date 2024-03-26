package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/event"
	"cc.allio/fusion/internal/svr"
	"cc.allio/fusion/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

const pipelinePathPrefix = "/api/admin/pipeline"

type PipelineRoute struct {
	Cfg             *config.Config
	PipelineService *svr.PipelineService
	Script          *event.ScriptEngine
}

var PipelineRouterSet = wire.NewSet(wire.Struct(new(PipelineRoute), "*"))

// GetAllPipelines
// @Summary get all pipeline
// @Schemes
// @Description get all pipeline
// @Tags Pipeline
// @Accept json
// @Produce json
// @Success 200 {object} []domain.Pipeline
// @Router /api/admin/pipeline [Get]
func (p *PipelineRoute) GetAllPipelines(c *gin.Context) *R {
	pipelines := p.PipelineService.GetAll()
	return Ok(pipelines)
}

// GetPipelineConfig
// @Summary get system default pipeline
// @Schemes
// @Description get system default pipeline
// @Tags Pipeline
// @Accept json
// @Produce json
// @Success 200 {object} []domain.EventItem
// @Router /api/admin/pipeline/config [Get]
func (p *PipelineRoute) GetPipelineConfig(c *gin.Context) *R {
	return Ok(domain.SystemEvents)
}

// GetPipelineById
// @Summary get pipeline by id
// @Schemes
// @Description get pipeline by id
// @Tags Pipeline
// @Accept json
// @Produce json
// @Success 200 {object} domain.Pipeline
// @Router /api/admin/pipeline/:id [Get]
func (p *PipelineRoute) GetPipelineById(c *gin.Context) *R {
	id := web.ParseNumberForPath(c, "id", -1)
	pipeline, err := p.PipelineService.GetPipelineById(int64(id))
	if err != nil {
		return InternalError(err)
	}
	return Ok(pipeline)
}

// CreatePipeline
// @Summary create pipeline
// @Schemes
// @Description create pipeline
// @Tags Pipeline
// @Accept json
// @Produce json
// @Param        pipeline   body      domain.Pipeline   true  "pipeline"
// @Success 200 {object} bool
// @Router /api/admin/pipeline [Post]
func (p *PipelineRoute) CreatePipeline(c *gin.Context) *R {
	pipeline := &domain.Pipeline{}
	if err := c.Bind(pipeline); err != nil {
		return InternalError(err)
	}
	saved, err := p.PipelineService.CreatePipeline(pipeline)
	if err != nil {
		return InternalError(err)
	}
	return Ok(saved)
}

// DeletePipeline
// @Summary delete pipeline
// @Schemes
// @Description delete pipeline
// @Tags Pipeline
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /api/admin/pipeline/:id [Delete]
func (p *PipelineRoute) DeletePipeline(c *gin.Context) *R {
	id := web.ParseNumberForPath(c, "id", -1)
	success, err := p.PipelineService.DeletePipeline(int64(id))
	if err != nil {
		return InternalError(err)
	}
	return Ok(success)
}

// UpdatePipeline
// @Summary update pipeline
// @Schemes
// @Description update pipeline
// @Tags Pipeline
// @Accept json
// @Produce json
// @Param        pipeline   body      domain.Pipeline   true  "pipeline"
// @Success 200 {object} bool
// @Router /api/admin/pipeline [Put]
func (p *PipelineRoute) UpdatePipeline(c *gin.Context) *R {
	pipeline := &domain.Pipeline{}
	if err := c.Bind(pipeline); err != nil {
		return InternalError(err)
	}
	success, err := p.PipelineService.UpdatePipeline(pipeline)
	if err != nil {
		return InternalError(err)
	}
	return Ok(success)
}

// TriggerPipeline
// @Summary trigger pipeline
// @Schemes
// @Description trigger pipeline
// @Tags Pipeline
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /api/admin/pipeline/trigger/:id [Get]
func (p *PipelineRoute) TriggerPipeline(c *gin.Context) *R {
	id := web.ParseNumberForPath(c, "id", -1)
	triggerPipeline := &credential.TriggerPipelineCredential{}
	if err := c.Bind(triggerPipeline); err != nil {
		return InternalError(err)
	}
	p.Script.DispatchById(int64(id), triggerPipeline.Input)
	return Ok(true)
}

func (p *PipelineRoute) Register(r *gin.Engine) {
	r.GET(pipelinePathPrefix, Handle(p.GetAllPipelines))
	r.GET(pipelinePathPrefix+"/config", Handle(p.GetPipelineConfig))
	r.GET(pipelinePathPrefix+"/:id", Handle(p.GetPipelineById))

	r.POST(pipelinePathPrefix, Handle(p.CreatePipeline))
	r.DELETE(pipelinePathPrefix+"/:id", Handle(p.DeletePipeline))
	r.PUT(pipelinePathPrefix+"/:id", Handle(p.UpdatePipeline))

	r.POST(pipelinePathPrefix+"/trigger/:id", Handle(p.TriggerPipeline))
}
