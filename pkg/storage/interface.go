package storage

import "golang.org/x/net/context"

type Filter interface {
	Field() []string
	Condition() (query interface{}, args []interface{})
}

// Everything is a Filter which accepts all objects.
var Everything Filter = everything{}

// everything is implementation of Everything.
type everything struct {
}

func (e everything) Field() []string {
	return []string{}
}

func (e everything) Condition() (query interface{}, args []interface{}) {
	return nil, nil
}

//Interface implement a storeage backend
type Interface interface {

	// Create adds a new object at a key unless it already exists. 'ttl' is time-to-live
	// in seconds (0 means forever). If no error is returned and out is not nil, out will be
	// set to the read value from database.
	Create(ctx context.Context, key string, obj, out interface{}) error

	// // Delete removes the specified key and returns the value that existed at that spot.
	// // If key didn't exist, it will return NotFound storage error.
	// Delete(ctx context.Context, key string, out runtime.Object, preconditions *Preconditions) error

	GetToList(ctx context.Context, filter Filter, result interface{}) error
	//keyField is index resource
	//updateFields will be only update that fileds in obj if that is null update all
	//obj is update object
	GuaranteedUpdate(ctx context.Context, keyField string, updateFields []string, obj interface{}) error
}
