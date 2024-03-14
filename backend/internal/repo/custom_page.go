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

type CustomPageRepository struct {
	Cfg *config.Config
	Db  *mongo.Database
}

var CustomPageRepositorySet = wire.NewSet(wire.Struct(new(CustomPageRepository), "*"))

func (a *CustomPageRepository) Save(insert *domain.CustomPage, opts ...*options.InsertOneOptions) (int64, error) {
	coll := a.Db.Collection(CustomPageCollection)
	return handleSave[domain.CustomPage](coll, insert, opts...)
}

func (a *CustomPageRepository) SaveMany(inserts []interface{}, opts ...*options.InsertManyOptions) ([]int64, error) {
	coll := a.Db.Collection(CustomPageCollection)
	return handleSaveMany(coll, inserts, opts...)
}

func (a *CustomPageRepository) Update(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(CustomPageCollection)
	return handleUpdate(coll, filter, update, opts...)
}

func (a *CustomPageRepository) UpdateMany(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(CustomPageCollection)
	return handleUpdateMany(coll, filter, update, opts...)
}

func (a *CustomPageRepository) Remove(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(CustomPageCollection)
	return handleRemove(coll, filter, opts...)
}

func (a *CustomPageRepository) RemoveMany(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(CustomPageCollection)
	return handleRemoveMany(coll, filter, opts...)
}

func (a *CustomPageRepository) FindOne(filter mongodb.Logical, opts ...*options.FindOneOptions) (*domain.CustomPage, error) {
	coll := a.Db.Collection(CustomPageCollection)
	return handleFindOne[domain.CustomPage](coll, func() *domain.CustomPage { return &domain.CustomPage{} }, filter, opts...)
}

func (a *CustomPageRepository) FindList(filter mongodb.Logical, opts ...*options.FindOptions) ([]*domain.CustomPage, error) {
	coll := a.Db.Collection(CustomPageCollection)
	return handleFindList[domain.CustomPage](coll, filter, opts...)
}

func (a *CustomPageRepository) Count(filter mongodb.Logical, opts ...*options.CountOptions) (int64, error) {
	coll := a.Db.Collection(CustomPageCollection)
	return handleCount(coll, filter, opts...)
}
