package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/repo"
	token2 "cc.allio/fusion/internal/token"
	"cc.allio/fusion/pkg/mongodb"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"
	"time"
)

type TokenService struct {
	Cfg         *config.Config
	SettingsSvr *SettingService
	TokenRepo   *repo.TokenRepository
}

var TokenServiceSet = wire.NewSet(wire.Struct(new(TokenService), "*"))

func (tokenSvr *TokenService) FindAllApiToken() []*domain.Token {
	tokens, err := tokenSvr.TokenRepo.FindList(mongodb.NewLogicalDefault(bson.E{Key: "id", Value: AdminId}))
	if err != nil {
		slog.Error("find all api token has error", "err", err)
		return make([]*domain.Token, 0)
	}
	return tokens
}

func (tokenSvr *TokenService) CreateApiToken(name string) (string, error) {
	return tokenSvr.CreateToken(AdminId, name)
}

func (tokenSvr *TokenService) CreateToken(userId int64, username string) (string, error) {
	slog.Info("created token", "username", username)
	loginSettings := tokenSvr.SettingsSvr.FindLoginSetting()
	expireIn := loginSettings.ExpiresIn
	exp := time.Now().Add(time.Duration(expireIn))
	tokenSigned, err := token2.CreateJwtToken(token2.Claims{
		Id:   userId,
		Name: username,
		Exp:  float64(exp.Unix()),
	}, tokenSvr.Cfg.Token.SignedKey)
	if err != nil {
		return "", err
	}
	defer func() {
		t := &domain.Token{
			UserId:    userId,
			Token:     tokenSigned,
			ExpiresIn: expireIn,
			CreatedAt: time.Now(),
			Disabled:  false,
		}
		_, err := tokenSvr.TokenRepo.Save(t)
		if err != nil {
			slog.Error("token persistence failed", "err", err, "token", t)
		}
	}()
	return tokenSigned, nil
}

func (tokenSvr *TokenService) DisabledToken(token string) (bool, error) {
	return tokenSvr.TokenRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "token", Value: token}), bson.D{{"disabled", true}})
}

func (tokenSvr *TokenService) DisabledAll() (bool, error) {
	return tokenSvr.TokenRepo.UpdateMany(mongodb.NewLogicalDefault(bson.E{Key: "disabled", Value: false}), primitive.D{{"disabled", true}})
}

func (tokenSvr *TokenService) DisabledTokenById(id int64) (bool, error) {
	return tokenSvr.TokenRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "id", Value: id}), bson.D{{"disabled", true}})
}
