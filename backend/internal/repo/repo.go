package repo

import (
	"cc.allio/fusion/pkg/mongodb"
	"context"
	"errors"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
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
	wire.Struct(new(Repository), "*"),
)

type DomainRepository[T interface{}] interface {
	Save(insert *T, opts ...*options.InsertOneOptions) (uint64, error)
	SaveMany(inserts []interface{}, opts ...*options.InsertManyOptions) ([]uint64, error)
	Update(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error)
	UpdateMany(filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error)
	Remove(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error)
	RemoveMany(filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error)
	Count(filter mongodb.Logical, opts ...*options.CountOptions) (int64, error)
	FindOne(filter mongodb.Logical, opts ...*options.FindOneOptions) (*T, error)
	FindList(filter mongodb.Logical, opts ...*options.FindOptions) ([]*T, error)
}

// handleSave handle domain object on save
func handleSave[T interface{}](coll *mongo.Collection, insert *T, opts ...*options.InsertOneOptions) (uint64, error) {
	id, err := setId(insert)
	if err != nil {
		return 0, err
	}
	_, err = coll.InsertOne(context.TODO(), insert, opts...)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// handleSaveMany handle many domain object on save
func handleSaveMany(coll *mongo.Collection, insert []interface{}, opts ...*options.InsertManyOptions) ([]uint64, error) {
	ids := make([]uint64, 0)
	for _, entity := range insert {
		id, err := setId(entity)
		if err != nil {
			return make([]uint64, 0), err
		}
		ids = append(ids, id)
	}
	return writeTransaction(coll.Database(), func(ctx mongo.SessionContext) ([]uint64, error) {
		result, err := coll.InsertMany(context.TODO(), insert, opts...)
		if err != nil {
			return make([]uint64, 0), err
		}
		if len(result.InsertedIDs) != len(insert) {
			err := ctx.AbortTransaction(context.TODO())
			return make([]uint64, 0), errors.Join(err, errors.New("Failed to save many. "))
		}
		err = ctx.CommitTransaction(context.TODO())
		if err != nil {
			return ids, err
		}
		return make([]uint64, 0), err
	})
}

// handleUpdate handle domain object on update
func handleUpdate(coll *mongo.Collection, filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	result, err := coll.UpdateOne(context.TODO(), filter.ToBson(), update, opts...)
	if err != nil {
		return false, err
	}
	return result.ModifiedCount > 0, nil
}

// handleUpdateMany handle update many domain entity
func handleUpdateMany(coll *mongo.Collection, filter mongodb.Logical, update bson.D, opts ...*options.UpdateOptions) (bool, error) {
	result, err := coll.UpdateMany(context.TODO(), filter.ToBson(), update, opts...)
	if err != nil {
		return false, err
	}
	return result.ModifiedCount > 0, nil
}

// handleRemove handle delete domain object
func handleRemove(coll *mongo.Collection, filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	result, err := coll.DeleteOne(context.TODO(), filter.ToBson(), opts...)
	if err != nil {
		return false, err
	}
	return result.DeletedCount > 0, nil
}

// handleRemoveMany handle delete domain object
func handleRemoveMany(coll *mongo.Collection, filter mongodb.Logical, opts ...*options.DeleteOptions) (bool, error) {
	result, err := coll.DeleteMany(context.TODO(), filter.ToBson(), opts...)
	if err != nil {
		return false, err
	}
	return result.DeletedCount > 0, nil
}

// handleCount handle domain object count of document
func handleCount(coll *mongo.Collection, filter mongodb.Logical, opts ...*options.CountOptions) (int64, error) {
	return coll.CountDocuments(context.TODO(), filter.ToBson(), opts...)
}

// handleFindOne handle find one
// implementation support domain entity by 'new' function
func handleFindOne[T interface{}](coll *mongo.Collection, new func() *T, filter mongodb.Logical, opts ...*options.FindOneOptions) (*T, error) {
	result := coll.FindOne(context.TODO(), filter.ToBson(), opts...)
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
	ctx := context.Background()
	cursor, err := coll.Find(ctx, filter.ToBson(), opts...)
	if err != nil {
		return nil, err
	}
	var domains = make([]*T, 0)
	if err = cursor.All(ctx, &domains); err != nil {
		return nil, err
	}
	return domains, nil
}

// writeTransaction wrap write concern writeTransaction open a write concern, and execute argument func f
func writeTransaction[T interface{}](db *mongo.Database, f func(ctx mongo.SessionContext) (T, error)) (T, error) {
	wc := db.WriteConcern()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	client := db.Client()
	session, err := client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(context.TODO())
	res, err := session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) { return f(ctx) }, txnOptions)
	result := res.(T)
	return result, err
}

// setId base on proxy set entity id
func setId(entity interface{}) (uint64, error) {
	v := reflect.ValueOf(entity)
	idValue := v.Elem().FieldByName("Id")
	nextId, err := mongodb.NextId()
	if err != nil {
		return 0, err
	}
	if idValue.CanSet() {
		idValue.SetUint(nextId)
		return nextId, nil
	} else {
		return 0, errors.New("unable set id")
	}
}
