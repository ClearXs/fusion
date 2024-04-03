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

type CategoryRepository struct {
	Cfg *config.Config
	Db  *mongo.Database
}

var CategoryRepositorySet = wire.NewSet(wire.Struct(new(CategoryRepository), "*"))

func (a *CategoryRepository) Save(insert *domain.Category, opts ...*options.InsertOneOptions) (string, error) {
	coll := a.Db.Collection(CategoryCollection)
	return handleSave[domain.Category](coll, insert, opts...)
}

func (a *CategoryRepository) SaveMany(inserts []interface{}, opts ...*options.InsertManyOptions) ([]string, error) {
	coll := a.Db.Collection(CategoryCollection)
	return handleSaveMany(coll, inserts, opts...)
}

func (a *CategoryRepository) Update(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(CategoryCollection)
	return handleUpdate(coll, filter, update, opts...)
}

func (a *CategoryRepository) UpdateMany(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(CategoryCollection)
	return handleUpdateMany(coll, filter, update, opts...)
}

func (a *CategoryRepository) Remove(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(CategoryCollection)
	return handleRemove(coll, filter, opts...)
}

func (a *CategoryRepository) RemoveMany(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(CategoryCollection)
	return handleRemoveMany(coll, filter, opts...)
}

func (a *CategoryRepository) FindOne(filter mongodb.Logical, opts ...*options.FindOneOptions) (*domain.Category, error) {
	coll := a.Db.Collection(CategoryCollection)
	return handleFindOne[domain.Category](coll, func() *domain.Category { return &domain.Category{} }, filter, opts...)
}

func (a *CategoryRepository) FindList(filter mongodb.Logical, opts ...*options.FindOptions) ([]*domain.Category, error) {
	coll := a.Db.Collection(CategoryCollection)
	return handleFindList[domain.Category](coll, filter, opts...)
}

func (a *CategoryRepository) Count(filter mongodb.Logical, opts ...*options.CountOptions) (int64, error) {
	coll := a.Db.Collection(CategoryCollection)
	return handleCount(coll, filter, opts...)
}
