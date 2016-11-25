package registry

import (
	"fmt"
	"strings"

	api "gofreezer/pkg/api"
	kubeerr "gofreezer/pkg/api/errors"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/api/validation/path"
	"gofreezer/pkg/fields"
	"gofreezer/pkg/labels"
	"gofreezer/pkg/pagination"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/storage"
)

// Store implements generic.Registry.
// It's intended to be embeddable, so that you can implement any
// non-generic functions if needed.
// You must supply a value for every field below before use; these are
// left public as it's meant to be overridable if need be.
// This object is intended to be copyable so that it can be used in
// different ways but share the same underlying behavior.
//
// The intended use of this type is embedding within a Kind specific
// RESTStorage implementation. This type provides CRUD semantics on
// a Kubelike resource, handling details like conflict detection with
// ResourceVersion and semantics. The RESTCreateStrategy and
// RESTUpdateStrategy are generic across all backends, and encapsulate
// logic specific to the API.
//
// TODO: make the default exposed methods exactly match a generic RESTStorage
type Store struct {
	// Called to make a new object, should return e.g., &api.Pod{}
	NewFunc func() runtime.Object

	// Called to make a new listing object, should return e.g., &api.PodList{}
	NewListFunc func() runtime.Object

	// Used for error reporting
	QualifiedResource unversioned.GroupResource

	// Used for listing/watching; should not include trailing "/"
	KeyRootFunc func(ctx api.Context) string

	// Called for Create/Update/Get/Delete. Note that 'namespace' can be
	// gotten from ctx.
	KeyFunc func(ctx api.Context, name string) (string, error)

	// Called to get the name of an object
	ObjectNameFunc func(obj runtime.Object) (string, error)

	// Return the TTL objects should be persisted with. Update is true if this
	// is an operation against an existing object. Existing is the current TTL
	// or the default for this operation.
	TTLFunc func(obj runtime.Object, existing uint64, update bool) (uint64, error)

	// Returns a matcher corresponding to the provided labels and fields.
	PredicateFunc func(label labels.Selector, field fields.Selector, page pagination.Pager) storage.SelectionPredicate

	// Called to cleanup storage clients.
	DestroyFunc func()

	// EnableGarbageCollection affects the handling of Update and Delete requests. It
	// must be synced with the corresponding flag in kube-controller-manager.
	EnableGarbageCollection bool

	// DeleteCollectionWorkers is the maximum number of workers in a single
	// DeleteCollection call.
	DeleteCollectionWorkers int

	// Decorator is called as exit hook on object returned from the underlying storage.
	// The returned object could be individual object (e.g. Pod) or the list type (e.g. PodList).
	// Decorator is intended for integrations that are above storage and
	// should only be used for specific cases where storage of the value is
	// not appropriate, since they cannot be watched.
	Decorator rest.ObjectFunc
	// Allows extended behavior during creation, required
	CreateStrategy rest.RESTCreateStrategy
	// On create of an object, attempt to run a further operation.
	AfterCreate rest.ObjectFunc
	// Allows extended behavior during updates, required
	UpdateStrategy rest.RESTUpdateStrategy
	// On update of an object, attempt to run a further operation.
	AfterUpdate rest.ObjectFunc
	// Allows extended behavior during updates, optional
	DeleteStrategy rest.RESTDeleteStrategy
	// On deletion of an object, attempt to run a further operation.
	AfterDelete rest.ObjectFunc
	// If true, return the object that was deleted. Otherwise, return a generic
	// success status response.
	ReturnDeletedObject bool
	// Allows extended behavior during export, optional
	ExportStrategy rest.RESTExportStrategy

	// Used for all storage access functions
	Storage storage.Interface
}

// NamespaceKeyRootFunc is the default function for constructing storage paths to resource directories enforcing namespace rules.
func NamespaceKeyRootFunc(ctx api.Context, prefix string) string {
	key := prefix
	ns, ok := api.NamespaceFrom(ctx)
	if ok && len(ns) > 0 {
		key = key + "/" + ns
	}
	return key
}

// NamespaceKeyFunc is the default function for constructing storage paths to a resource relative to prefix enforcing namespace rules.
// If no namespace is on context, it errors.
func NamespaceKeyFunc(ctx api.Context, prefix string, name string) (string, error) {
	key := NamespaceKeyRootFunc(ctx, prefix)
	ns, ok := api.NamespaceFrom(ctx)
	if !ok || len(ns) == 0 {
		return "", kubeerr.NewBadRequest("Namespace parameter required.")
	}
	if len(name) == 0 {
		return "", kubeerr.NewBadRequest("Name parameter required.")
	}
	if msgs := path.IsValidPathSegmentName(name); len(msgs) != 0 {
		return "", kubeerr.NewBadRequest(fmt.Sprintf("Name parameter invalid: %q: %s", name, strings.Join(msgs, ";")))
	}
	key = key + "/" + name
	return key, nil
}

// NoNamespaceKeyFunc is the default function for constructing storage paths to a resource relative to prefix without a namespace
func NoNamespaceKeyFunc(ctx api.Context, prefix string, name string) (string, error) {
	if len(name) == 0 {
		return "", kubeerr.NewBadRequest("Name parameter required.")
	}
	if msgs := path.IsValidPathSegmentName(name); len(msgs) != 0 {
		return "", kubeerr.NewBadRequest(fmt.Sprintf("Name parameter invalid: %q: %s", name, strings.Join(msgs, ";")))
	}
	key := prefix + "/" + name
	return key, nil
}
