/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package etcds

import (
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/storage"
	"gofreezer/pkg/watch"

	"golang.org/x/net/context"
)

// Pass an UpdateFunc to Interface.GuaranteedUpdate to make an update
// that is guaranteed to succeed.
// See the comment for GuaranteedUpdate for more details.
type UpdateFunc func(input runtime.Object, res storage.ResponseMeta) (output runtime.Object, ttl *uint64, err error)

// Interface offers a common interface for object marshaling/unmarshaling operations and
// hides all the storage-related operations behind it.
type Interface interface {
	// Returns Versioner associated with this interface.
	Versioner() storage.Versioner

	// Create adds a new object at a key unless it already exists. 'ttl' is time-to-live
	// in seconds (0 means forever). If no error is returned and out is not nil, out will be
	// set to the read value from database.
	Create(ctx context.Context, key string, obj, out runtime.Object, ttl uint64) error

	// Delete removes the specified key and returns the value that existed at that spot.
	// If key didn't exist, it will return NotFound storage error.
	Delete(ctx context.Context, key string, out runtime.Object, preconditions *storage.Preconditions) error

	// Watch begins watching the specified key. Events are decoded into API objects,
	// and any items selected by 'p' are sent down to returned watch.Interface.
	// resourceVersion may be used to specify what version to begin watching,
	// which should be the current resourceVersion, and no longer rv+1
	// (e.g. reconnecting without missing any updates).
	Watch(ctx context.Context, key string, resourceVersion string, p storage.SelectionPredicate) (watch.Interface, error)

	// WatchList begins watching the specified key's items. Items are decoded into API
	// objects and any item selected by 'p' are sent down to returned watch.Interface.
	// resourceVersion may be used to specify what version to begin watching,
	// which should be the current resourceVersion, and no longer rv+1
	// (e.g. reconnecting without missing any updates).
	WatchList(ctx context.Context, key string, resourceVersion string, p storage.SelectionPredicate) (watch.Interface, error)

	// Get unmarshals json found at key into objPtr. On a not found error, will either
	// return a zero object of the requested type, or an error, depending on ignoreNotFound.
	// Treats empty responses and nil response nodes exactly like a not found error.
	Get(ctx context.Context, key string, objPtr runtime.Object, ignoreNotFound bool) error

	// GetToList unmarshals json found at key and opaque it into *List api object
	// (an object that satisfies the runtime.IsList definition).
	GetToList(ctx context.Context, key string, resourceVersion string, p storage.SelectionPredicate, listObj runtime.Object) error

	// List unmarshalls jsons found at directory defined by key and opaque them
	// into *List api object (an object that satisfies runtime.IsList definition).
	// The returned contents may be delayed, but it is guaranteed that they will
	// be have at least 'resourceVersion'.
	List(ctx context.Context, key string, resourceVersion string, p storage.SelectionPredicate, listObj runtime.Object) error

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
	GuaranteedUpdate(
		ctx context.Context, key string, ptrToType runtime.Object, ignoreNotFound bool,
		precondtions *storage.Preconditions, tryUpdate UpdateFunc, suggestion ...runtime.Object) error
}
