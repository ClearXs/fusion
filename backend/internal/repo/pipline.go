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

type PipelineRepository struct {
	Cfg *config.Config
	Db  *mongo.Database
}

var PipelineRepositorySet = wire.NewSet(wire.Struct(new(PipelineRepository), "*"))

func (a *PipelineRepository) Save(insert *domain.Pipeline, opts ...*options.InsertOneOptions) (uint64, error) {
	coll := a.Db.Collection(PipelineCollection)
	return handleSave[domain.Pipeline](coll, insert, opts...)
}

func (a *PipelineRepository) SaveMany(inserts []interface{}, opts ...*options.InsertManyOptions) ([]uint64, error) {
	coll := a.Db.Collection(PipelineCollection)
	return handleSaveMany(coll, inserts, opts...)
}

func (a *PipelineRepository) Update(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(PipelineCollection)
	return handleUpdate(coll, filter, update, opts...)
}

func (a *PipelineRepository) UpdateMany(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(PipelineCollection)
	return handleUpdateMany(coll, filter, update, opts...)
}

func (a *PipelineRepository) Remove(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(PipelineCollection)
	return handleRemove(coll, filter, opts...)
}

func (a *PipelineRepository) RemoveMany(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(PipelineCollection)
	return handleRemoveMany(coll, filter, opts...)
}

func (a *PipelineRepository) FindOne(filter mongodb.Logical, opts ...*options.FindOneOptions) (*domain.Pipeline, error) {
	coll := a.Db.Collection(PipelineCollection)
	return handleFindOne[domain.Pipeline](coll, func() *domain.Pipeline { return &domain.Pipeline{} }, filter, opts...)
}

func (a *PipelineRepository) FindList(filter mongodb.Logical, opts ...*options.FindOptions) ([]*domain.Pipeline, error) {
	coll := a.Db.Collection(PipelineCollection)
	return handleFindList[domain.Pipeline](coll, filter, opts...)
}

func (a *PipelineRepository) Count(filter mongodb.Logical, opts ...*options.CountOptions) (int64, error) {
	coll := a.Db.Collection(PipelineCollection)
	return handleCount(coll, filter, opts...)
}
