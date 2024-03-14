package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"errors"
	"fmt"
	"github.com/google/wire"
)

type TokenUser struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

type AuthService struct {
	Cfg          *config.Config
	UserService  *UserService
	TokenService *TokenService
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
	return &TokenUser{Token: token, User: user}, nil
}
