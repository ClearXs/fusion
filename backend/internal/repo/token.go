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

type TokenRepository struct {
	Cfg *config.Config
	Db  *mongo.Database
}

var TokenRepositorySet = wire.NewSet(wire.Struct(new(TokenRepository), "*"))

func (a *TokenRepository) Save(insert *domain.Token, opts ...*options.InsertOneOptions) (uint64, error) {
	coll := a.Db.Collection(TokenCollection)
	return handleSave[domain.Token](coll, insert, opts...)
}

func (a *TokenRepository) SaveMany(inserts []interface{}, opts ...*options.InsertManyOptions) ([]uint64, error) {
	coll := a.Db.Collection(TokenCollection)
	return handleSaveMany(coll, inserts, opts...)
}

func (a *TokenRepository) Update(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(TokenCollection)
	return handleUpdate(coll, filter, update, opts...)
}

func (a *TokenRepository) UpdateMany(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(TokenCollection)
	return handleUpdateMany(coll, filter, update, opts...)
}

func (a *TokenRepository) Remove(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(TokenCollection)
	return handleRemove(coll, filter, opts...)
}

func (a *TokenRepository) RemoveMany(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(TokenCollection)
	return handleRemoveMany(coll, filter, opts...)
}

func (a *TokenRepository) FindOne(filter mongodb.Logical, opts ...*options.FindOneOptions) (*domain.Token, error) {
	coll := a.Db.Collection(TokenCollection)
	return handleFindOne[domain.Token](coll, func() *domain.Token { return &domain.Token{} }, filter, opts...)
}

func (a *TokenRepository) FindList(filter mongodb.Logical, opts ...*options.FindOptions) ([]*domain.Token, error) {
	coll := a.Db.Collection(TokenCollection)
	return handleFindList[domain.Token](coll, filter, opts...)
}

func (a *TokenRepository) Count(filter mongodb.Logical, opts ...*options.CountOptions) (int64, error) {
	coll := a.Db.Collection(TokenCollection)
	return handleCount(coll, filter, opts...)
}
