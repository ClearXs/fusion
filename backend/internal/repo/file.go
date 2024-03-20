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

type FileRepository struct {
	Cfg *config.Config
	Db  *mongo.Database
}

var FileRepositorySet = wire.NewSet(wire.Struct(new(FileRepository), "*"))

func (a *FileRepository) Save(insert *domain.File, opts ...*options.InsertOneOptions) (int64, error) {
	coll := a.Db.Collection(FileCollection)
	return handleSave[domain.File](coll, insert, opts...)
}

func (a *FileRepository) SaveMany(inserts []interface{}, opts ...*options.InsertManyOptions) ([]int64, error) {
	coll := a.Db.Collection(FileCollection)
	return handleSaveMany(coll, inserts, opts...)
}

func (a *FileRepository) Update(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(FileCollection)
	return handleUpdate(coll, filter, update, opts...)
}

func (a *FileRepository) UpdateMany(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(FileCollection)
	return handleUpdateMany(coll, filter, update, opts...)
}

func (a *FileRepository) Remove(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(FileCollection)
	return handleRemove(coll, filter, opts...)
}

func (a *FileRepository) RemoveMany(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(FileCollection)
	return handleRemoveMany(coll, filter, opts...)
}

func (a *FileRepository) FindOne(filter mongodb.Logical, opts ...*options.FindOneOptions) (*domain.File, error) {
	coll := a.Db.Collection(FileCollection)
	return handleFindOne[domain.File](coll, func() *domain.File { return &domain.File{} }, filter, opts...)
}

func (a *FileRepository) FindList(filter mongodb.Logical, opts ...*options.FindOptions) ([]*domain.File, error) {
	coll := a.Db.Collection(FileCollection)
	return handleFindList[domain.File](coll, filter, opts...)
}

func (a *FileRepository) Count(filter mongodb.Logical, opts ...*options.CountOptions) (int64, error) {
	coll := a.Db.Collection(FileCollection)
	return handleCount(coll, filter, opts...)
}
