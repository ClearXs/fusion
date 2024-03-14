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

type VisitRepository struct {
	Cfg *config.Config
	Db  *mongo.Database
}

var VisitRepositorySet = wire.NewSet(wire.Struct(new(VisitRepository), "*"))

func (a *VisitRepository) Save(insert *domain.Visit, opts ...*options.InsertOneOptions) (int64, error) {
	coll := a.Db.Collection(VisitCollection)
	return handleSave[domain.Visit](coll, insert, opts...)
}

func (a *VisitRepository) SaveMany(inserts []interface{}, opts ...*options.InsertManyOptions) ([]int64, error) {
	coll := a.Db.Collection(VisitCollection)
	return handleSaveMany(coll, inserts, opts...)
}

func (a *VisitRepository) Update(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(VisitCollection)
	return handleUpdate(coll, filter, update, opts...)
}

func (a *VisitRepository) UpdateMany(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(VisitCollection)
	return handleUpdateMany(coll, filter, update, opts...)
}

func (a *VisitRepository) Remove(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(VisitCollection)
	return handleRemove(coll, filter, opts...)
}

func (a *VisitRepository) RemoveMany(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(VisitCollection)
	return handleRemoveMany(coll, filter, opts...)
}

func (a *VisitRepository) FindOne(filter mongodb.Logical, opts ...*options.FindOneOptions) (*domain.Visit, error) {
	coll := a.Db.Collection(VisitCollection)
	return handleFindOne[domain.Visit](coll, func() *domain.Visit { return &domain.Visit{} }, filter, opts...)
}

func (a *VisitRepository) FindList(filter mongodb.Logical, opts ...*options.FindOptions) ([]*domain.Visit, error) {
	coll := a.Db.Collection(VisitCollection)
	return handleFindList[domain.Visit](coll, filter, opts...)
}

func (a *VisitRepository) Count(filter mongodb.Logical, opts ...*options.CountOptions) (int64, error) {
	coll := a.Db.Collection(VisitCollection)
	return handleCount(coll, filter, opts...)
}
