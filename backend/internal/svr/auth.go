package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/apm"
	"cc.allio/fusion/internal/domain"
	"errors"
	"fmt"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
)

type TokenUser struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

type AuthService struct {
	Cfg          *config.Config
	UserService  *UserService
	TokenService *TokenService
	LogService   *LogService
	Logger       *apm.Logger
}

var AuthServiceSet = wire.NewSet(wire.Struct(new(AuthService), "*"))

func (auth *AuthService) Login(username string, password string) (*TokenUser, error) {
	user := auth.UserService.GetUserByUsernameAndPassword(username, password)
	if user == nil {
		return nil, errors.New(fmt.Sprintf("user by username: [%s] and password: [%s] not found instance, please check", username, password))
	}
	token, err := auth.TokenService.CreateToken(user.Id, user.Name)
	if err != nil {
		return nil, err
	}
	// log
	auth.Logger.Record(bson.D{{"success", true}, {"logType", apm.BusinessLogType}, {"username", username}, {"businessType", LoginLogType}})
	return &TokenUser{Token: token, User: user}, nil
}
