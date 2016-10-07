package storage

import "golang.org/x/net/context"

type RetrieveFilter interface {
	// Field() []string
	//Condition retrieve condition
	//resultFieldSets: what fields needed by result
	Field() []string
	//query:specific a plain sql  like as : ("name = ? AND age >= ?", "jinzhu", "22") for mysql
	//queryArgs: query args place here
	Condition() (query interface{}, args []interface{})
	//sort: sort by field
	Sort() interface{} //sort by field in interface{}
	//limit: what number of recorde will be return
	Limit() interface{} //Specify the number of records to be retrieved
	//skip: with above Condition ,offset number of record
	Skip() interface{} //skip the number of records
}

//Interface implement a storeage backend
type Interface interface {

	// Create adds a new object at a key unless it already exists. 'ttl' is time-to-live
	// in seconds (0 means forever). If no error is returned and out is not nil, out will be
	// set to the read value from database.
	Create(ctx context.Context, key string, obj, out interface{}) error

	// // Delete removes the specified key and returns the value that existed at that spot.
	// // If key didn't exist, it will return NotFound storage error.
	//Delete(ctx context.Context, key string, out runtime.Object, preconditions *Preconditions) error
	//Delete(ctx context.Context, key string, obj interface{}) error
	Delete(ctx context.Context, filter RetrieveFilter, out interface{}) error

	//filter support query arg
	GetToList(ctx context.Context, filter RetrieveFilter, result interface{}) error

	//filter support query arg
	GetCount(ctx context.Context, filter RetrieveFilter, result *uint64) error

	//keyField is index resource
	//updateFields will be only update that fileds in obj if that is null update all
	//obj is update object
	GuaranteedUpdate(ctx context.Context, keyField string, updateFields []string, obj interface{}) error
}
