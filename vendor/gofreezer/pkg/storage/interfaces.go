package storage

import (
	"gofreezer/pkg/fields"
	"gofreezer/pkg/labels"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/types"
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

// FilterFunc takes an API object and returns true if the object satisfies some requirements.
// TODO: We will remove this type and use SelectionPredicate everywhere.
type FilterFunc func(obj runtime.Object) bool

// Everything accepts all objects.
var Everything = SelectionPredicate{
	Label: labels.Everything(),
	Field: fields.Everything(),
}

// Pass an UpdateFunc to Interface.GuaranteedUpdate to make an update
// that is guaranteed to succeed.
// See the comment for GuaranteedUpdate for more details.
type UpdateFunc func(input runtime.Object, res ResponseMeta) (output runtime.Object, ttl *uint64, err error)

// Preconditions must be fulfilled before an operation (update, delete, etc.) is carried out.
type Preconditions struct {
	// Specifies the target UID.
	// +optional
	UID *types.UID `json:"uid,omitempty"`
}

// NewUIDPreconditions returns a Preconditions with UID set.
func NewUIDPreconditions(uid string) *Preconditions {
	u := types.UID(uid)
	return &Preconditions{UID: &u}
}

//Interface  wrap storage driver
type Interface interface {
	//storage type
	Type() string
}
