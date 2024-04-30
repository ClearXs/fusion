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

type DraftRepository struct {
	Cfg *config.Config
	Db  *mongo.Database
}

var DraftRepositorySet = wire.NewSet(wire.Struct(new(DraftRepository), "*"))

func (a *DraftRepository) Save(insert *domain.Draft, opts ...*options.InsertOneOptions) (uint64, error) {
	coll := a.Db.Collection(DraftCollection)
	return handleSave[domain.Draft](coll, insert, opts...)
}

func (a *DraftRepository) SaveMany(inserts []interface{}, opts ...*options.InsertManyOptions) ([]uint64, error) {
	coll := a.Db.Collection(DraftCollection)
	return handleSaveMany(coll, inserts, opts...)
}

func (a *DraftRepository) Update(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(DraftCollection)
	return handleUpdate(coll, filter, update, opts...)
}

func (a *DraftRepository) UpdateMany(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(DraftCollection)
	return handleUpdateMany(coll, filter, update, opts...)
}

func (a *DraftRepository) Remove(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(DraftCollection)
	return handleRemove(coll, filter, opts...)
}

func (a *DraftRepository) RemoveMany(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(DraftCollection)
	return handleRemoveMany(coll, filter, opts...)
}

func (a *DraftRepository) FindOne(filter mongodb.Logical, opts ...*options.FindOneOptions) (*domain.Draft, error) {
	coll := a.Db.Collection(DraftCollection)
	return handleFindOne[domain.Draft](coll, func() *domain.Draft { return &domain.Draft{} }, filter, opts...)
}

func (a *DraftRepository) FindList(filter mongodb.Logical, opts ...*options.FindOptions) ([]*domain.Draft, error) {
	coll := a.Db.Collection(DraftCollection)
	return handleFindList[domain.Draft](coll, filter, opts...)
}

func (a *DraftRepository) Count(filter mongodb.Logical, opts ...*options.CountOptions) (int64, error) {
	coll := a.Db.Collection(DraftCollection)
	return handleCount(coll, filter, opts...)
}
