package repo

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/pkg/mongodb"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	Cfg *config.Config
	Db  *mongo.Database
}

var UserRepositorySet = wire.NewSet(wire.Struct(new(UserRepository), "*"))

func (a *UserRepository) Save(insert *domain.User, opts ...*options.InsertOneOptions) (string, error) {
	coll := a.Db.Collection(UserCollection)
	return handleSave[domain.User](coll, insert, opts...)
}

func (a *UserRepository) SaveMany(inserts []interface{}, opts ...*options.InsertManyOptions) ([]string, error) {
	coll := a.Db.Collection(UserCollection)
	return handleSaveMany(coll, inserts, opts...)
}

func (a *UserRepository) Update(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(UserCollection)
	return handleUpdate(coll, filter, update, opts...)
}

func (a *UserRepository) UpdateMany(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(UserCollection)
	return handleUpdateMany(coll, filter, update, opts...)
}

func (a *UserRepository) Remove(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(UserCollection)
	return handleRemove(coll, filter, opts...)
}

func (a *UserRepository) RemoveMany(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(UserCollection)
	return handleRemoveMany(coll, filter, opts...)
}

func (a *UserRepository) FindOne(filter mongodb.Logical, opts ...*options.FindOneOptions) (*domain.User, error) {
	coll := a.Db.Collection(UserCollection)
	return handleFindOne[domain.User](coll, func() *domain.User { return &domain.User{} }, filter, opts...)
}

func (a *UserRepository) FindList(filter mongodb.Logical, opts ...*options.FindOptions) ([]*domain.User, error) {
	coll := a.Db.Collection(UserCollection)
	return handleFindList[domain.User](coll, filter, opts...)
}

func (a *UserRepository) Count(filter mongodb.Logical, opts ...*options.CountOptions) (int64, error) {
	coll := a.Db.Collection(UserCollection)
	return handleCount(coll, filter, opts...)
}
