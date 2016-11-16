package mongodb

import (
	"fmt"
	"gofreezer/pkg/runtime"
	storage "gofreezer/pkg/storage"
	storagehelper "gofreezer/pkg/storage/helper"
	"gofreezer/pkg/storage/mongodbs/client"
	"reflect"
	"time"

	mgo "gopkg.in/mgo.v2"
)

type Collection struct {
	name        string
	database    string
	keyIndex    []string
	expireIndex []string
}

const (
	//if set expire index. if index expire then will be after this period to remove doc
	expirePeriod = time.Duration(0) * time.Second
	//this index is our convention  for runtime object
	uidIndex    = "uid"
	keyIndex    = "key"
	expireIndex = "ttl"
)

func GetCollection(dbName string, sess *mgo.Session, obj runtime.Object) (*Collection, error) {
	collection := storagehelper.GetObjKind(obj)
	if len(collection) == 0 {
		return nil, storage.NewInternalError(fmt.Sprintf("object(%v) not have kind", reflect.TypeOf(obj)))
	}

	c := &Collection{
		name:        collection,
		database:    dbName,
		keyIndex:    []string{uidIndex, keyIndex},
		expireIndex: []string{expireIndex},
	}

	//ensure index
	err := c.CreateIndex(c.GetRequestMeta(sess))
	if err != nil {
		return nil, storage.NewInternalError(err.Error())
	}

	return c, nil
}

//CreateIndex by runtime object
func (c *Collection) CreateIndex(meta *client.RequestMeta) error {
	err := client.MongoEnsureIndex(meta, c.keyIndex)
	if err != nil {
		return err
	}

	err = client.MongoEnsureIndexWithExpire(meta, c.expireIndex, expirePeriod)
	if err != nil {
		return err
	}
	return nil
}

func (c *Collection) GetRequestMeta(sess *mgo.Session) *client.RequestMeta {
	return &client.RequestMeta{
		DBName:     c.database,
		Collection: c.name,
		Sess:       sess,
	}
}
