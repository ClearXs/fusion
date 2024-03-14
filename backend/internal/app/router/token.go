package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/svr"
	"cc.allio/fusion/pkg/web"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
)

const TokenPathPrefix = "/api/admin/token"

type TokenRouter struct {
	Cfg          *config.Config
	TokenService *svr.TokenService
}

var TokenRouterSet = wire.NewSet(wire.Struct(new(TokenRouter), "*"))

// GetAllApiTokens
// @Summary get all tokens
// @Schemes
// @Description get all tokens
// @Tags Token
// @Accept json
// @Produce json
// @Success 200 {object} []domain.Token
// @Router /api/admin/token [Get]
func (t *TokenRouter) GetAllApiTokens(c *gin.Context) *R {
	tokens := t.TokenService.FindAllApiToken()
	return Ok(tokens)
}

// CreateApiToken
// @Summary create token
// @Schemes
// @Description create token
// @Tags Token
// @Accept json
// @Produce json
// @Param        tokenCredential   body      credential.TokenCredential   true  "tokenCredential"
// @Success 200 {object} string
// @Router /api/admin/token [Post]
func (t *TokenRouter) CreateApiToken(c *gin.Context) *R {
	if t.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	tokenCredential := &credential.TokenCredential{}
	if err := c.Bind(tokenCredential); err != nil {
		return InternalError(err)
	}
	token, err := t.TokenService.CreateApiToken(tokenCredential.Name)
	if err != nil {
		return InternalError(err)
	}
	return Ok(token)
}

// DeleteApiToken
// @Summary delete token by id
// @Schemes
// @Description delete token by id
// @Tags Token
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /api/admin/token/:id [Delete]
func (t *TokenRouter) DeleteApiToken(c *gin.Context) *R {
	if t.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	id := web.ParseNumberForPath(c, "id", -1)
	success, err := t.TokenService.DisabledTokenById(int64(id))
	if err != nil {
		return InternalError(err)
	}
	return Ok(success)
}

func (t *TokenRouter) Register(r *gin.Engine) {
	r.GET(TokenPathPrefix, Handle(t.GetAllApiTokens))
	r.POST(TokenPathPrefix, Handle(t.CreateApiToken))
	r.DELETE(TokenPathPrefix+"/:id", Handle(t.DeleteApiToken))
}
