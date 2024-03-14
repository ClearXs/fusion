package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/svr"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
)

const RewardPathPrefix = "/api/admin/meta/reward"

type RewardRouter struct {
	Cfg         *config.Config
	MetaService *svr.MetaService
}

var RewardRouterSet = wire.NewSet(wire.Struct(new(RewardRouter), "*"))

// GetReward
// @Summary get reward list
// @Schemes
// @Description get reward list
// @Tags Reward
// @Accept json
// @Produce json
// @Success 200 {object} []domain.RewardItem
// @Router /api/admin/meta/reward [Get]
func (re *RewardRouter) GetReward(c *gin.Context) *R {
	rewards := re.MetaService.GetRewards()
	return Ok(rewards)
}

// UpdateReward
// @Summary update reward
// @Schemes
// @Description update reward
// @Tags Reward
// @Accept json
// @Produce json
// @Param        reward   body      domain.RewardItem   true  "reward"
// @Success 200 {object} bool
// @Router /api/admin/meta/reward [Put]
func (re *RewardRouter) UpdateReward(c *gin.Context) *R {
	if re.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	reward := &domain.RewardItem{}
	if err := c.Bind(reward); err != nil {
		return InternalError(err)
	}
	updated, err := re.MetaService.AddOrUpdateReward(reward)
	if err != nil {
		return InternalError(err)
	}
	return Ok(updated)
}

// CreateReward
// @Summary create reward
// @Schemes
// @Description create reward
// @Tags Reward
// @Accept json
// @Produce json
// @Param        reward   body      domain.RewardItem   true  "reward"
// @Success 200 {object} bool
// @Router /api/admin/meta/reward [Post]
func (re *RewardRouter) CreateReward(c *gin.Context) *R {
	if re.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	reward := &domain.RewardItem{}
	if err := c.Bind(reward); err != nil {
		return InternalError(err)
	}
	added, err := re.MetaService.AddOrUpdateReward(reward)
	if err != nil {
		return InternalError(err)
	}
	return Ok(added)
}

// DeleteReward
// @Summary delete reward by reward name
// @Schemes
// @Description delete reward by reward name
// @Tags Reward
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /api/admin/meta/reward/:name [Delete]
func (re *RewardRouter) DeleteReward(c *gin.Context) *R {
	if re.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	name := c.Query("name")
	deleted, err := re.MetaService.DeleteRewardByName(name)
	if err != nil {
		return InternalError(err)
	}
	return Ok(deleted)
}

func (re *RewardRouter) Register(r *gin.Engine) {
	r.GET(RewardPathPrefix, Handle(re.GetReward))
	r.PUT(RewardPathPrefix, Handle(re.UpdateReward))
	r.POST(RewardPathPrefix, Handle(re.CreateReward))
	r.DELETE(RewardPathPrefix+"/:name", Handle(re.DeleteReward))
}
