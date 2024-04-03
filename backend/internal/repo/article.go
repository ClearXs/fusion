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

type ArticleRepository struct {
	Cfg *config.Config
	Db  *mongo.Database
}

var ArticleRepositorySet = wire.NewSet(wire.Struct(new(ArticleRepository), "*"))

func (a *ArticleRepository) Save(insert *domain.Article, opts ...*options.InsertOneOptions) (string, error) {
	coll := a.Db.Collection(ArticleCollection)
	return handleSave[domain.Article](coll, insert, opts...)
}

func (a *ArticleRepository) SaveMany(inserts []interface{}, opts ...*options.InsertManyOptions) ([]string, error) {
	coll := a.Db.Collection(ArticleCollection)
	return handleSaveMany(coll, inserts, opts...)
}

func (a *ArticleRepository) Update(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(ArticleCollection)
	return handleUpdate(coll, filter, update, opts...)
}

func (a *ArticleRepository) UpdateMany(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	coll := a.Db.Collection(ArticleCollection)
	return handleUpdateMany(coll, filter, update, opts...)
}

func (a *ArticleRepository) Remove(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(ArticleCollection)
	return handleRemove(coll, filter, opts...)
}

func (a *ArticleRepository) RemoveMany(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	coll := a.Db.Collection(ArticleCollection)
	return handleRemoveMany(coll, filter, opts...)
}

func (a *ArticleRepository) FindOne(filter mongodb.Logical, opts ...*options.FindOneOptions) (*domain.Article, error) {
	coll := a.Db.Collection(ArticleCollection)
	return handleFindOne[domain.Article](coll, func() *domain.Article { return &domain.Article{} }, filter, opts...)
}

func (a *ArticleRepository) FindList(filter mongodb.Logical, opts ...*options.FindOptions) ([]*domain.Article, error) {
	coll := a.Db.Collection(ArticleCollection)
	return handleFindList[domain.Article](coll, filter, opts...)
}

func (a *ArticleRepository) Count(filter mongodb.Logical, opts ...*options.CountOptions) (int64, error) {
	coll := a.Db.Collection(ArticleCollection)
	return handleCount(coll, filter, opts...)
}
