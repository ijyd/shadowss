package mysqls

import (
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/storage"

	"golang.org/x/net/context"
)

// Pass an UpdateFunc to Interface.GuaranteedUpdate to make an update
// that is guaranteed to succeed.
// See the comment for GuaranteedUpdate for more details.
type UpdateFunc func(input runtime.Object) (output runtime.Object, updateField []string, err error)

//Interface implement a storeage backend
type Interface interface {
	// Returns Versioner associated with this interface.
	Versioner() storage.Versioner

	Create(ctx context.Context, key string, obj, out runtime.Object) error
	//cannot get item by key, so need selection predicate
	Delete(ctx context.Context, key string, out runtime.Object, preconditions *storage.Preconditions) error

	Get(ctx context.Context, key string, objPtr runtime.Object, ignoreNotFound bool) error
	GetToList(ctx context.Context, key string, p storage.SelectionPredicate, listObj runtime.Object) error

	//cannot get item by key, so need selection predicate
	GuaranteedUpdate(ctx context.Context, key string, ptrToType runtime.Object, ignoreNotFound bool,
		precondtions *storage.Preconditions, tryUpdate UpdateFunc, suggestion ...runtime.Object) error
}
