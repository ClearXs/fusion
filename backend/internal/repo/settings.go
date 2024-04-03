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

type SettingsRepository struct {
	Cfg *config.Config
	Db  *mongo.Database
}

var SettingsRepositorySet = wire.NewSet(wire.Struct(new(SettingsRepository), "*"))

func (a *SettingsRepository) Save(insert *domain.Setting, opts ...*options.InsertOneOptions) (string, error) {
	coll := a.Db.Collection(SettingsCollection)
	return handleSave[domain.Setting](coll, insert, opts...)
}

func (a *SettingsRepository) SaveMany(inserts []interface{}, opts ...*options.InsertManyOptions) ([]string, error) {
	coll := a.Db.Collection(SettingsCollection)
	return handleSaveMany(coll, inserts, opts...)
}

func (a *SettingsRepository) Update(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(SettingsCollection)
	return handleUpdate(coll, filter, update, opts...)
}

func (a *SettingsRepository) UpdateMany(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(SettingsCollection)
	return handleUpdateMany(coll, filter, update, opts...)
}

func (a *SettingsRepository) Remove(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(SettingsCollection)
	return handleRemove(coll, filter, opts...)
}

func (a *SettingsRepository) RemoveMany(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(SettingsCollection)
	return handleRemoveMany(coll, filter, opts...)
}

func (a *SettingsRepository) FindOne(filter mongodb.Logical, opts ...*options.FindOneOptions) (*domain.Setting, error) {
	coll := a.Db.Collection(SettingsCollection)
	return handleFindOne[domain.Setting](coll, func() *domain.Setting { return &domain.Setting{} }, filter, opts...)
}

func (a *SettingsRepository) FindList(filter mongodb.Logical, opts ...*options.FindOptions) ([]*domain.Setting, error) {
	coll := a.Db.Collection(SettingsCollection)
	return handleFindList[domain.Setting](coll, filter, opts...)
}

func (a *SettingsRepository) Count(filter mongodb.Logical, opts ...*options.CountOptions) (int64, error) {
	coll := a.Db.Collection(SettingsCollection)
	return handleCount(coll, filter, opts...)
}
