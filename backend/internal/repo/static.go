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

type StaticRepository struct {
	Cfg *config.Config
	Db  *mongo.Database
}

var StaticRepositorySet = wire.NewSet(wire.Struct(new(StaticRepository), "*"))

func (a *StaticRepository) Save(insert *domain.Static, opts ...*options.InsertOneOptions) (int64, error) {
	coll := a.Db.Collection(StaticCollection)
	return handleSave[domain.Static](coll, insert, opts...)
}

func (a *StaticRepository) SaveMany(inserts []interface{}, opts ...*options.InsertManyOptions) ([]int64, error) {
	coll := a.Db.Collection(StaticCollection)
	return handleSaveMany(coll, inserts, opts...)
}

func (a *StaticRepository) Update(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(StaticCollection)
	return handleUpdate(coll, filter, update, opts...)
}

func (a *StaticRepository) UpdateMany(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(StaticCollection)
	return handleUpdateMany(coll, filter, update, opts...)
}

func (a *StaticRepository) Remove(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(StaticCollection)
	return handleRemove(coll, filter, opts...)
}

func (a *StaticRepository) RemoveMany(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(StaticCollection)
	return handleRemoveMany(coll, filter, opts...)
}

func (a *StaticRepository) FindOne(filter mongodb.Logical, opts ...*options.FindOneOptions) (*domain.Static, error) {
	coll := a.Db.Collection(StaticCollection)
	return handleFindOne[domain.Static](coll, func() *domain.Static { return &domain.Static{} }, filter, opts...)
}

func (a *StaticRepository) FindList(filter mongodb.Logical, opts ...*options.FindOptions) ([]*domain.Static, error) {
	coll := a.Db.Collection(StaticCollection)
	return handleFindList[domain.Static](coll, filter, opts...)
}

func (a *StaticRepository) Count(filter mongodb.Logical, opts ...*options.CountOptions) (int64, error) {
	coll := a.Db.Collection(StaticCollection)
	return handleCount(coll, filter, opts...)
}
