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

type MetaRepository struct {
	Cfg *config.Config
	Db  *mongo.Database
}

var MetaRepositorySet = wire.NewSet(wire.Struct(new(MetaRepository), "*"))

func (a *MetaRepository) Save(insert *domain.Meta, opts ...*options.InsertOneOptions) (string, error) {
	coll := a.Db.Collection(MetaCollection)
	return handleSave[domain.Meta](coll, insert, opts...)
}

func (a *MetaRepository) SaveMany(inserts []interface{}, opts ...*options.InsertManyOptions) ([]string, error) {
	coll := a.Db.Collection(MetaCollection)
	return handleSaveMany(coll, inserts, opts...)
}

func (a *MetaRepository) Update(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(MetaCollection)
	return handleUpdate(coll, filter, update, opts...)
}

func (a *MetaRepository) UpdateMany(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(MetaCollection)
	return handleUpdateMany(coll, filter, update, opts...)
}

func (a *MetaRepository) Remove(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(MetaCollection)
	return handleRemove(coll, filter, opts...)
}

func (a *MetaRepository) RemoveMany(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(MetaCollection)
	return handleRemoveMany(coll, filter, opts...)
}

func (a *MetaRepository) FindOne(filter mongodb.Logical, opts ...*options.FindOneOptions) (*domain.Meta, error) {
	coll := a.Db.Collection(MetaCollection)
	return handleFindOne[domain.Meta](coll, func() *domain.Meta { return &domain.Meta{} }, filter, opts...)
}

func (a *MetaRepository) FindList(filter mongodb.Logical, opts ...*options.FindOptions) ([]*domain.Meta, error) {
	coll := a.Db.Collection(MetaCollection)
	return handleFindList[domain.Meta](coll, filter, opts...)
}

func (a *MetaRepository) Count(filter mongodb.Logical, opts ...*options.CountOptions) (int64, error) {
	coll := a.Db.Collection(MetaCollection)
	return handleCount(coll, filter, opts...)
}
