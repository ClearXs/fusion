package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/internal/token"
	"cc.allio/fusion/pkg/mongodb"
	"cc.allio/fusion/pkg/util"
	"errors"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slog"
)

const AdminId uint64 = 0

type UserService struct {
	Cfg            *config.Config
	UserRepository *repo.UserRepository
}

var UserServiceSet = wire.NewSet(wire.Struct(new(UserService), "*"))

func (userSvr *UserService) GetUser() *domain.User {
	user, err := userSvr.UserRepository.FindOne(mongodb.NewLogical().Append(bson.E{Key: "id", Value: AdminId}))
	if err != nil {
		slog.Error("Get user has error", "err", err)
		return &domain.User{}
	}
	return user
}

func (userSvr *UserService) getUser() (*domain.User, error) {
	user, err := userSvr.UserRepository.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "id", Value: AdminId}))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userSvr *UserService) GetUserList() []*domain.User {
	userList, err := userSvr.UserRepository.FindList(mongodb.NewLogical())
	if err != nil {
		slog.Error(err.Error())
		return make([]*domain.User, 0)
	}
	return userList
}

func (userSvr *UserService) GetUserByUsernameAndPassword(username string, password string) *domain.User {
	user, err := userSvr.UserRepository.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "name", Value: username}))
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	if user == nil {
		return nil
	}
	encryptPassword := token.EncryptPassword(user.Name, password, user.Salt)
	result, err := userSvr.UserRepository.FindOne(mongodb.NewLogicalDefaultArray(bson.D{{"name", username}, {"password", encryptPassword}}))
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	if result == nil {
		return nil
	}
	defer func() {
		// update salt
		newSalt := token.MakeSalt()
		userSvr.UserRepository.Update(mongodb.NewLogicalDefault(bson.E{Key: "id", Value: result.Id}), bson.D{{"salt", newSalt}})
	}()
	return result
}

func (userSvr *UserService) UpdateUser(user *domain.UpdateUser) (bool, error) {
	curUser, err := userSvr.getUser()
	if err != nil {
		return false, nil
	}
	return userSvr.UserRepository.Update(
		mongodb.NewLogicalDefault(bson.E{Key: "id", Value: curUser.Id}),
		bson.D{{"name", user.Name}, {"nickname", user.Nickname}, {"password", token.EncryptPassword(user.Name, user.Password, curUser.Salt)}},
	)
}

func (userSvr *UserService) GetAllCollaborators() []*domain.User {
	collaborators, err := userSvr.UserRepository.FindList(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: domain.CollaborateUserType}))
	if err != nil {
		slog.Error("get system all collaborators has error", "err", err)
		return make([]*domain.User, 0)
	}
	return collaborators
}

func (userSvr *UserService) DeleteCollaborator(id int) bool {
	removed, err := userSvr.UserRepository.Remove(mongodb.NewLogicalDefault(bson.E{Key: "id", Value: id}))
	if err != nil {
		slog.Error("remove collaborator has err", "err", err)
		return false
	}
	return removed
}

func (userSvr *UserService) InsertCollaborator(collaborator *domain.User) (bool, error) {
	user, err := userSvr.getCollaboratorByName(collaborator.Name)
	if err != nil {
		return false, err
	}
	if user != nil {
		return false, errors.New("exist collaborator")
	}
	save, err := userSvr.UserRepository.Save(collaborator)
	if err != nil {
		return false, err
	}
	return save > 0, nil
}

func (userSvr *UserService) UpdateCollaborator(collaborator *domain.User) (bool, error) {
	user, err := userSvr.getCollaboratorByName(collaborator.Name)
	if err != nil {
		return false, err
	}
	if user != nil {
		return false, errors.New("collaborator non existent")
	}

	salt := token.MakeSalt()
	password := token.EncryptPassword(collaborator.Name, collaborator.Password, salt)
	collaborator.Password = password
	collaborator.Salt = salt
	elements := util.ToBsonElements(collaborator)
	filter := mongodb.NewLogicalOrDefaultArray(bson.D{{"id", user.Id}, {"type", domain.CollaborateUserType}})
	return userSvr.UserRepository.Update(filter, elements)
}

func (userSvr *UserService) getCollaboratorByName(name string) (*domain.User, error) {
	return userSvr.UserRepository.FindOne(mongodb.NewLogicalDefaultArray(bson.D{{"name", name}, {"type", domain.CollaborateUserType}}))
}
