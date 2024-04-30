package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/svr"
	"cc.allio/fusion/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

const CollaboratorPathPrefix = "/api/admin/collaborator"

type CollaboratorRoute struct {
	Cfg          *config.Config
	UserService  *svr.UserService
	MetaService  *svr.MetaService
	TokenService *svr.TokenService
}

var CollaboratorRouterSet = wire.NewSet(wire.Struct(new(CollaboratorRoute), "*"))

// GetAllCollaborators
// @Summary get all collaborator
// @Schemes
// @Description get all collaborator
// @Tags Collaborator
// @Accept json
// @Produce json
// @Success 200 {object} []domain.User
// @Router /api/admin/collaborator [Get]
func (collaborator *CollaboratorRoute) GetAllCollaborators(c *gin.Context) *R {
	collaborators := collaborator.UserService.GetAllCollaborators()
	return Ok(collaborators)
}

// GetAllCollaboratorsList
// @Summary get all collaborator of list
// @Schemes
// @Description get all collaborator of list
// @Tags Collaborator
// @Accept json
// @Produce json
// @Success 200 {object} []domain.User
// @Router /api/admin/collaborator/list [Get]
func (collaborator *CollaboratorRoute) GetAllCollaboratorsList(c *gin.Context) *R {
	siteInfo := collaborator.MetaService.GetSiteInfo()
	admin := collaborator.UserService.GetUser()
	collaborators := collaborator.UserService.GetAllCollaborators()
	admin.Id = 0
	admin.Nickname = siteInfo.Author
	allCollaborators := append([]*domain.User{admin}, collaborators...)
	return Ok(allCollaborators)
}

// DeleteCollaborator
// @Summary delete collaborator by id
// @Schemes
// @Description delete collaborator by id
// @Tags Collaborator
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /api/admin/collaborator/:id [Delete]
func (collaborator *CollaboratorRoute) DeleteCollaborator(c *gin.Context) *R {
	id := web.ParseNumberForPath(c, "id", -1)
	successed := collaborator.UserService.DeleteCollaborator(id)
	return Ok(successed)
}

// CreateCollaborator
// @Summary create collaborator
// @Schemes
// @Description create collaborator
// @Tags Collaborator
// @Accept json
// @Produce json
// @Param        user   body      domain.User   true  "user"
// @Success 200 {object} bool
// @Router /api/admin/collaborator [Post]
func (collaborator *CollaboratorRoute) CreateCollaborator(c *gin.Context) *R {
	collaboratorUser := &domain.User{}
	if err := c.Bind(collaboratorUser); err != nil {
		return InternalError(err)
	}
	successed, err := collaborator.UserService.InsertCollaborator(collaboratorUser)
	if err != nil {
		return InternalError(err)
	}
	return Ok(successed)
}

// UpdateCollaborator
// @Summary update collaborator
// @Schemes
// @Description update collaborator
// @Tags Collaborator
// @Accept json
// @Produce json
// @Param        user   body      domain.User   true  "user"
// @Success 200 {object} bool
// @Router /api/admin/collaborator [Put]
func (collaborator *CollaboratorRoute) UpdateCollaborator(c *gin.Context) *R {
	collaboratorUser := &domain.User{}
	if err := c.Bind(collaboratorUser); err != nil {
		return InternalError(err)
	}
	successed, err := collaborator.UserService.UpdateCollaborator(collaboratorUser)
	if err != nil {
		return InternalError(err)
	}
	return Ok(successed)
}

func (collaborator *CollaboratorRoute) Register(r *gin.Engine) {
	r.GET(CollaboratorPathPrefix, Handle(collaborator.GetAllCollaborators))

	r.GET(CollaboratorPathPrefix+"/list", Handle(collaborator.GetAllCollaboratorsList))
	r.POST(CollaboratorPathPrefix, Handle(collaborator.CreateCollaborator))
	r.PUT(CollaboratorPathPrefix, Handle(collaborator.UpdateCollaborator))
	r.DELETE(CollaboratorPathPrefix+"/:id", Handle(collaborator.DeleteCollaborator))
}
