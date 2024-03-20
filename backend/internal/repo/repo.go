package repo

import (
	"cc.allio/fusion/pkg/mongodb"
	"cc.allio/fusion/pkg/util"
	"context"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	VisitCollection      = "visits"
	ArticleCollection    = "articles"
	MetaCollection       = "metas"
	SettingsCollection   = "settings"
	TokenCollection      = "tokens"
	UserCollection       = "users"
	ViewerCollection     = "viewers"
	CategoryCollection   = "categories"
	DraftCollection      = "drafts"
	StaticCollection     = "statics"
	CustomPageCollection = "custompages"
	PipelineCollection   = "pipelines"
	FileCollection       = "files"
)

type Repository struct {
	ArticleRepository    *ArticleRepository
	CategoryRepository   *CategoryRepository
	DraftRepository      *DraftRepository
	MetaRepository       *MetaRepository
	SettingsRepository   *SettingsRepository
	StaticRepository     *StaticRepository
	TokenRepository      *TokenRepository
	UserRepository       *UserRepository
	ViewerRepository     *ViewerRepository
	VisitRepository      *VisitRepository
	CustomPageRepository *CustomPageRepository
	PipelineRepository   *PipelineRepository
	FileRepository       *FileRepository
}

var RepositorySet = wire.NewSet(
	ArticleRepositorySet,
	CategoryRepositorySet,
	DraftRepositorySet,
	MetaRepositorySet,
	SettingsRepositorySet,
	StaticRepositorySet,
	TokenRepositorySet,
	UserRepositorySet,
	ViewerRepositorySet,
	VisitRepositorySet,
	CustomPageRepositorySet,
	PipelineRepositorySet,
	FileRepositorySet,
	wire.Struct(new(Repository), "*"),
)

type DomainRepository[T interface{}] interface {
	Save(insert *T, opts ...*options.InsertOneOptions) (int64, error)
	SaveMany(inserts []interface{}, opts ...*options.InsertManyOptions) ([]int64, error)
	Update(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error)
	UpdateMany(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error)
	Remove(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error)
	RemoveMany(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error)
	Count(filter mongodb.Logical, opts ...*options.CountOptions) (int64, error)
	FindOne(filter mongodb.Logical, opts ...*options.FindOneOptions) (*T, error)
	FindList(filter mongodb.Logical, opts ...*options.FindOptions) ([]*T, error)
}

// handleSave handle domain object on save
func handleSave[T interface{}](coll *mongo.Collection, insert *T, opts ...*options.InsertOneOptions) (int64, error) {
	result, err := coll.InsertOne(context.TODO(), insert, opts...)
	if err != nil {
		return int64(-1), err
	}
	idString := string(result.InsertedID.([]byte))
	id := util.ToStringInt(idString, -1)
	return int64(id), nil
}

func handleSaveMany(coll *mongo.Collection, insert []interface{}, opts ...*options.InsertManyOptions) ([]int64, error) {
	result, err := coll.InsertMany(context.TODO(), insert, opts...)
	if err != nil {
		return make([]int64, 0), err
	}
	ids := make([]int64, 0)
	for _, insertedId := range result.InsertedIDs {
		idString := string(insertedId.([]byte))
		id := util.ToStringInt(idString, -1)
		ids = append(ids, int64(id))
	}
	return ids, nil
}

// handleUpdate handle domain object on update
func handleUpdate(coll *mongo.Collection, filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	result, err := coll.UpdateOne(context.TODO(), filter, update, opts...)
	if err != nil {
		return false, err
	}
	return result.ModifiedCount > 0, nil
}

// handleUpdateMany handle update many domain entity
func handleUpdateMany(coll *mongo.Collection, filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	result, err := coll.UpdateMany(context.TODO(), filter, update, opts...)
	if err != nil {
		return false, err
	}
	return result.ModifiedCount > 0, nil
}

// handleRemove handle delete domain object
func handleRemove(coll *mongo.Collection, filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	result, err := coll.DeleteOne(context.TODO(), filter, opts...)
	if err != nil {
		return false, err
	}
	return result.DeletedCount > 0, nil
}

// handleRemoveMany handle delete domain object
func handleRemoveMany(coll *mongo.Collection, filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	result, err := coll.DeleteMany(context.TODO(), filter, opts...)
	if err != nil {
		return false, err
	}
	return result.DeletedCount > 0, nil
}

// handleCount handle domain object count of document
func handleCount(coll *mongo.Collection, filter mongodb.Logical, opts ...*options.CountOptions) (int64, error) {
	return coll.CountDocuments(context.TODO(), filter.GetBson(), opts...)
}

// handleFindOne handle find one
// implementation support domain entity by 'new' function
func handleFindOne[T interface{}](coll *mongo.Collection, new func() *T, filter mongodb.Logical, opts ...*options.FindOneOptions) (*T, error) {
	result := coll.FindOne(context.TODO(), filter, opts...)
	if err := result.Err(); err != nil {
		return nil, err
	}
	domainObject := new()
	if err := result.Decode(domainObject); err != nil {
		return nil, err
	}
	return domainObject, nil
}

// handleFindList handle find list
func handleFindList[T interface{}](coll *mongo.Collection, filter mongodb.Logical, opts ...*options.FindOptions) ([]*T, error) {
	cursor, err := coll.Find(context.TODO(), filter.GetBson(), opts...)
	if err != nil {
		return nil, err
	}
	var domains = make([]*T, 0)
	if err = cursor.Decode(domains); err != nil {
		return nil, err
	}
	return domains, nil
}
