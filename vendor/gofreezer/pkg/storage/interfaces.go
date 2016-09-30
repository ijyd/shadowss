package storage

import (
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/types"
	"gofreezer/pkg/watch"

	"golang.org/x/net/context"
)

// Versioner abstracts setting and retrieving metadata fields from database response
// onto the object ot list.
type Versioner interface {
	// UpdateObject sets storage metadata into an API object. Returns an error if the object
	// cannot be updated correctly. May return nil if the requested object does not need metadata
	// from database.
	UpdateObject(obj runtime.Object, resourceVersion uint64) error
	// UpdateList sets the resource version into an API list object. Returns an error if the object
	// cannot be updated correctly. May return nil if the requested object does not need metadata
	// from database.
	UpdateList(obj runtime.Object, resourceVersion uint64) error
	// ObjectResourceVersion returns the resource version (for persistence) of the specified object.
	// Should return an error if the specified object does not have a persistable version.
	ObjectResourceVersion(obj runtime.Object) (uint64, error)
}

// ResponseMeta contains information about the database metadata that is associated with
// an object. It abstracts the actual underlying objects to prevent coupling with concrete
// database and to improve testability.
type ResponseMeta struct {
	// TTL is the time to live of the node that contained the returned object. It may be
	// zero or negative in some cases (objects may be expired after the requested
	// expiration time due to server lag).
	TTL int64
	// The resource version of the node that contained the returned object.
	ResourceVersion uint64
}

// MatchValue defines a pair (<index name>, <value for that index>).
type MatchValue struct {
	IndexName string
	Value     string
}

// TriggerPublisherFunc is a function that takes an object, and returns a list of pairs
// (<index name>, <index value for the given object>) for all indexes known
// to that function.
type TriggerPublisherFunc func(obj runtime.Object) []MatchValue

type RetrieveCondition struct {
	ResultFieldSets []string
	Query           interface{}
	QueryArgs       []interface{}
	Sort            interface{}
	Limit           interface{}
	Skip            interface{}
}

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

// Filter is interface that is used to pass filtering mechanism.
type Filter interface {
	// Filter is a predicate which takes an API object and returns true
	// if and only if the object should remain in the set.
	Filter(obj runtime.Object) bool
	// For any triggers known to the Filter, if Filter() can return only
	// (a subset of) objects for which indexing function returns <value>,
	// (<index name>, <value> pair would be returned.
	//
	// This is optimization to avoid computing Filter() function (which are
	// usually relatively expensive) in case we are sure they will return
	// false anyway.
	Trigger() []MatchValue

	//Retrieve() RetrieveFilter
}

// Everything is a Filter which accepts all objects.
var Everything Filter = everything{}

// everything is implementation of Everything.
type everything struct {
}

func (e everything) Filter(runtime.Object) bool {
	return true
}

func (e everything) Trigger() []MatchValue {
	return nil
}

// func (e everything) Retrieve() RetrieveFilter {
// 	return nil
// }

// Pass an UpdateFunc to Interface.GuaranteedUpdate to make an update
// that is guaranteed to succeed.
// See the comment for GuaranteedUpdate for more details.
type UpdateFunc func(input runtime.Object, res ResponseMeta) (output runtime.Object, ttl *uint64, err error)

// Preconditions must be fulfilled before an operation (update, delete, etc.) is carried out.
type Preconditions struct {
	// Specifies the target UID.
	UID *types.UID `json:"uid,omitempty"`
}

// NewUIDPreconditions returns a Preconditions with UID set.
func NewUIDPreconditions(uid string) *Preconditions {
	u := types.UID(uid)
	return &Preconditions{UID: &u}
}

// Interface offers a common interface for object marshaling/unmarshaling operations and
// hides all the storage-related operations behind it.
type Interface interface {
	// Returns Versioner associated with this interface.
	Versioner() Versioner

	// Create adds a new object at a key unless it already exists. 'ttl' is time-to-live
	// in seconds (0 means forever). If no error is returned and out is not nil, out will be
	// set to the read value from database.
	Create(ctx context.Context, key string, obj, out runtime.Object, ttl uint64) error

	// Delete removes the specified key and returns the value that existed at that spot.
	// If key didn't exist, it will return NotFound storage error.
	Delete(ctx context.Context, key string, out runtime.Object, preconditions *Preconditions) error

	// Watch begins watching the specified key. Events are decoded into API objects,
	// and any items passing 'filter' are sent down to returned watch.Interface.
	// resourceVersion may be used to specify what version to begin watching,
	// which should be the current resourceVersion, and no longer rv+1
	// (e.g. reconnecting without missing any updates).
	Watch(ctx context.Context, key string, resourceVersion string, filter Filter) (watch.Interface, error)

	// WatchList begins watching the specified key's items. Items are decoded into API
	// objects and any item passing 'filter' are sent down to returned watch.Interface.
	// resourceVersion may be used to specify what version to begin watching,
	// which should be the current resourceVersion, and no longer rv+1
	// (e.g. reconnecting without missing any updates).
	WatchList(ctx context.Context, key string, resourceVersion string, filter Filter) (watch.Interface, error)

	// Get unmarshals json found at key into objPtr. On a not found error, will either
	// return a zero object of the requested type, or an error, depending on ignoreNotFound.
	// Treats empty responses and nil response nodes exactly like a not found error.
	Get(ctx context.Context, key string, objPtr runtime.Object, ignoreNotFound bool) error

	// GetToList unmarshals json found at key and opaque it into *List api object
	// (an object that satisfies the runtime.IsList definition).
	GetToList(ctx context.Context, key string, filter Filter, listObj runtime.Object) error

	// List unmarshalls jsons found at directory defined by key and opaque them
	// into *List api object (an object that satisfies runtime.IsList definition).
	// The returned contents may be delayed, but it is guaranteed that they will
	// be have at least 'resourceVersion'.
	List(ctx context.Context, key string, resourceVersion string, filter Filter, listObj runtime.Object) error

	// GuaranteedUpdate keeps calling 'tryUpdate()' to update key 'key' (of type 'ptrToType')
	// retrying the update until success if there is index conflict.
	// Note that object passed to tryUpdate may change across invocations of tryUpdate() if
	// other writers are simultaneously updating it, so tryUpdate() needs to take into account
	// the current contents of the object when deciding how the update object should look.
	// If the key doesn't exist, it will return NotFound storage error if ignoreNotFound=false
	// or zero value in 'ptrToType' parameter otherwise.
	// If the object to update has the same value as previous, it won't do any update
	// but will return the object in 'ptrToType' parameter.
	//
	// Example:
	//
	// s := /* implementation of Interface */
	// err := s.GuaranteedUpdate(
	//     "myKey", &MyType{}, true,
	//     func(input runtime.Object, res ResponseMeta) (runtime.Object, *uint64, error) {
	//       // Before each incovation of the user defined function, "input" is reset to
	//       // current contents for "myKey" in database.
	//       curr := input.(*MyType)  // Guaranteed to succeed.
	//
	//       // Make the modification
	//       curr.Counter++
	//
	//       // Return the modified object - return an error to stop iterating. Return
	//       // a uint64 to alter the TTL on the object, or nil to keep it the same value.
	//       return cur, nil, nil
	//    }
	// })
	GuaranteedUpdate(ctx context.Context, key string, ptrToType runtime.Object, ignoreNotFound bool, precondtions *Preconditions, tryUpdate UpdateFunc) error
}

//Interface implement a storeage backend
// type Interface interface {
//
// 	// Create adds a new object at a key unless it already exists. 'ttl' is time-to-live
// 	// in seconds (0 means forever). If no error is returned and out is not nil, out will be
// 	// set to the read value from database.
// 	Create(ctx context.Context, key string, obj, out interface{}) error
//
// 	// // Delete removes the specified key and returns the value that existed at that spot.
// 	// // If key didn't exist, it will return NotFound storage error.
// 	//Delete(ctx context.Context, key string, out runtime.Object, preconditions *Preconditions) error
// 	//Delete(ctx context.Context, key string, obj interface{}) error
// 	Delete(ctx context.Context, filter RetrieveFilter, out interface{}) error
//
// 	//filter support query arg
// 	GetToList(ctx context.Context, filter RetrieveFilter, result interface{}) error
//
// 	//filter support query arg
// 	GetCount(ctx context.Context, filter RetrieveFilter, result *uint64) error
//
// 	//keyField is index resource
// 	//updateFields will be only update that fileds in obj if that is null update all
// 	//obj is update object
// 	GuaranteedUpdate(ctx context.Context, keyField string, updateFields []string, obj interface{}) error
// }
