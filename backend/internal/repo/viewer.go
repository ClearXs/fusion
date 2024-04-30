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

type ViewerRepository struct {
	Cfg *config.Config
	Db  *mongo.Database
}

var ViewerRepositorySet = wire.NewSet(wire.Struct(new(ViewerRepository), "*"))

func (a *ViewerRepository) Save(insert *domain.Viewer, opts ...*options.InsertOneOptions) (uint64, error) {
	coll := a.Db.Collection(ViewerCollection)
	return handleSave[domain.Viewer](coll, insert, opts...)
}

func (a *ViewerRepository) SaveMany(inserts []interface{}, opts ...*options.InsertManyOptions) ([]uint64, error) {
	coll := a.Db.Collection(ViewerCollection)
	return handleSaveMany(coll, inserts, opts...)
}

func (a *ViewerRepository) Update(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(ViewerCollection)
	return handleUpdate(coll, filter, update, opts...)
}

func (a *ViewerRepository) UpdateMany(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(ViewerCollection)
	return handleUpdateMany(coll, filter, update, opts...)
}

func (a *ViewerRepository) Remove(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(ViewerCollection)
	return handleRemove(coll, filter, opts...)
}

func (a *ViewerRepository) RemoveMany(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(ViewerCollection)
	return handleRemoveMany(coll, filter, opts...)
}

func (a *ViewerRepository) FindOne(filter mongodb.Logical, opts ...*options.FindOneOptions) (*domain.Viewer, error) {
	coll := a.Db.Collection(ViewerCollection)
	return handleFindOne[domain.Viewer](coll, func() *domain.Viewer { return &domain.Viewer{} }, filter, opts...)
}

func (a *ViewerRepository) FindList(filter mongodb.Logical, opts ...*options.FindOptions) ([]*domain.Viewer, error) {
	coll := a.Db.Collection(ViewerCollection)
	return handleFindList[domain.Viewer](coll, filter, opts...)
}

func (a *ViewerRepository) Count(filter mongodb.Logical, opts ...*options.CountOptions) (int64, error) {
	coll := a.Db.Collection(ViewerCollection)
	return handleCount(coll, filter, opts...)
}
