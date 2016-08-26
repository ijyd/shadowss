package storage

import "golang.org/x/net/context"

// type Filter struct {
// 	ResultFields []string
// 	//plain sql like as : ("name = ? AND age >= ?", "jinzhu", "22")
// 	//Query := string("name = ? AND age >= ?")
// 	Query interface{}
// 	//QueryArgs := make([]interface{}, 2)
// 	//QueryArgs[0] = string("jinzhu")
// 	//QueryArgs[1] = string("22")
// 	QueryArgs []interface{}
// }

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
	GetToList(ctx context.Context, filter Filter, result interface{}) error
	//keyField is index resource
	//updateFields will be only update that fileds in obj if that is null update all
	//obj is update object
	GuaranteedUpdate(ctx context.Context, keyField string, updateFields []string, obj interface{}) error
}
