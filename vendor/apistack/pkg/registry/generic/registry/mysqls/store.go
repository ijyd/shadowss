package mysqls

import (
	"fmt"
	"reflect"
	"sync"

	"apistack/pkg/registry/generic/registry"
	api "gofreezer/pkg/api"
	apierr "gofreezer/pkg/api/errors"
	storeerr "gofreezer/pkg/api/errors/storage"
	"gofreezer/pkg/api/meta"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/fields"
	"gofreezer/pkg/labels"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/storage"
	"gofreezer/pkg/storage/mysqls"
	utilruntime "gofreezer/pkg/util/runtime"
	"gofreezer/pkg/watch"

	"github.com/golang/glog"
)

// TODO: make the default exposed methods exactly match a generic RESTStorage
type Store struct {
	registry.Store
	// Used for all storage access functions
	Storage mysqls.Interface
}

// NewStore implements Store
func NewStore(store registry.Store) *Store {
	return &Store{
		Store:   store,
		Storage: store.Storage.(mysqls.Interface),
	}
}

const OptimisticLockErrorMsg = "the object has been modified; please apply your changes to the latest version and try again"

// New implements RESTStorage
func (e *Store) New() runtime.Object {
	return e.NewFunc()
}

// NewList implements RESTLister
func (e *Store) NewList() runtime.Object {
	return e.NewListFunc()
}

// List returns a list of items matching labels and field
func (e *Store) List(ctx api.Context, options *api.ListOptions) (runtime.Object, error) {
	label := labels.Everything()
	if options != nil && options.LabelSelector != nil {
		label = options.LabelSelector
	}
	field := fields.Everything()
	if options != nil && options.FieldSelector != nil {
		field = options.FieldSelector
	}
	out, err := e.ListPredicate(ctx, e.PredicateFunc(label, field), options)
	if err != nil {
		return nil, err
	}
	if e.Decorator != nil {
		if err := e.Decorator(out); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// ListPredicate returns a list of all the items matching m.
func (e *Store) ListPredicate(ctx api.Context, p storage.SelectionPredicate, options *api.ListOptions) (runtime.Object, error) {
	if options == nil {
		options = &api.ListOptions{ResourceVersion: "0"}
	}
	list := e.NewListFunc()
	if name, ok := p.MatchesSingle(); ok {
		if key, err := e.KeyFunc(ctx, name); err == nil {
			err := e.Storage.GetToList(ctx, key, p, list)
			return list, storeerr.InterpretListError(err, e.QualifiedResource)
		}
		// if we cannot extract a key based on the current context, the optimization is skipped
	}

	err := e.Storage.GetToList(ctx, e.KeyRootFunc(ctx), p, list)
	return list, storeerr.InterpretListError(err, e.QualifiedResource)
}

// Create inserts a new item according to the unique key from the object.
func (e *Store) Create(ctx api.Context, obj runtime.Object) (runtime.Object, error) {
	if err := rest.BeforeCreate(e.CreateStrategy, ctx, obj); err != nil {
		return nil, err
	}
	name, err := e.ObjectNameFunc(obj)
	if err != nil {
		return nil, err
	}
	key, err := e.KeyFunc(ctx, name)
	if err != nil {
		return nil, err
	}
	out := e.NewFunc()
	if err := e.Storage.Create(ctx, key, obj, out); err != nil {
		err = storeerr.InterpretCreateError(err, e.QualifiedResource, name)
		err = rest.CheckGeneratedNameError(e.CreateStrategy, err, obj)
		if !apierr.IsAlreadyExists(err) {
			return nil, err
		}
		return nil, err
	}
	glog.V(5).Infof("Get create obj %+v", out)
	if e.AfterCreate != nil {
		if err := e.AfterCreate(out); err != nil {
			return nil, err
		}
	}
	if e.Decorator != nil {
		if err := e.Decorator(obj); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// shouldDelete checks if a Update is removing all the object's finalizers. If so,
// it further checks if the object's DeletionGracePeriodSeconds is 0. If so, it
// returns true.
func (e *Store) shouldDelete(ctx api.Context, key string, obj, existing runtime.Object) bool {
	if !e.EnableGarbageCollection {
		return false
	}
	newMeta, err := api.ObjectMetaFor(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return false
	}
	oldMeta, err := api.ObjectMetaFor(existing)
	if err != nil {
		utilruntime.HandleError(err)
		return false
	}
	return len(newMeta.Finalizers) == 0 && oldMeta.DeletionGracePeriodSeconds != nil && *oldMeta.DeletionGracePeriodSeconds == 0
}

func (e *Store) deleteForEmptyFinalizers(ctx api.Context, name, key string, obj runtime.Object, preconditions *storage.Preconditions) (runtime.Object, bool, error) {
	out := e.NewFunc()
	glog.V(6).Infof("going to delete %s from regitry, triggered by update", name)
	if err := e.Storage.Delete(ctx, key, out, nil); err != nil {
		// Deletion is racy, i.e., there could be multiple update
		// requests to remove all finalizers from the object, so we
		// ignore the NotFound error.
		if storage.IsNotFound(err) {
			_, err := e.finalizeDelete(obj, true)
			// clients are expecting an updated object if a PUT succeeded,
			// but finalizeDelete returns a unversioned.Status, so return
			// the object in the request instead.
			return obj, false, err
		}
		return nil, false, storeerr.InterpretDeleteError(err, e.QualifiedResource, name)
	}
	_, err := e.finalizeDelete(out, true)
	// clients are expecting an updated object if a PUT succeeded, but
	// finalizeDelete returns a unversioned.Status, so return the object in
	// the request instead.
	return obj, false, err
}

// Update performs an atomic update and set of the object. Returns the result of the update
// or an error. If the registry allows create-on-update, the create flow will be executed.
// A bool is returned along with the object and any errors, to indicate object creation.
func (e *Store) Update(ctx api.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	key, err := e.KeyFunc(ctx, name)
	if err != nil {
		return nil, false, err
	}

	var (
		creatingObj runtime.Object
		creating    = false
	)

	storagePreconditions := &storage.Preconditions{}
	if preconditions := objInfo.Preconditions(); preconditions != nil {
		storagePreconditions.UID = preconditions.UID
	}

	out := e.NewFunc()
	// deleteObj is only used in case a deletion is carried out
	var deleteObj runtime.Object
	err = e.Storage.GuaranteedUpdate(ctx, key, out, true, nil, func(existing runtime.Object) (output runtime.Object, updateField []string, err error) {
		// Given the existing object, get the new object
		obj, err := objInfo.UpdatedObject(ctx, existing)
		glog.V(9).Infof("updated object(%v) err:%v", existing, err)
		if err != nil {
			return nil, nil, err
		}

		version, err := e.Storage.Versioner().ObjectResourceVersion(existing)
		glog.V(9).Infof("updated object(%v) err:%v", version, err)
		if err != nil {
			return nil, nil, err
		}
		if version == 0 {
			glog.V(9).Infof("Will create obj instead of update err")
			if !e.UpdateStrategy.AllowCreateOnUpdate() {
				return nil, nil, apierr.NewNotFound(e.QualifiedResource, name)
			}
			creating = true
			creatingObj = obj
			if err := rest.BeforeCreate(e.CreateStrategy, ctx, obj); err != nil {
				glog.V(9).Infof("Will BeforeCreate obj instead of update err")
				return nil, nil, err
			}
			return obj, nil, nil
		}

		creating = false
		creatingObj = nil
		glog.V(9).Infof("Will update 2")
		if err := rest.BeforeUpdate(e.UpdateStrategy, ctx, obj, existing); err != nil {
			return nil, nil, err
		}
		glog.V(9).Infof("Will update")
		delete := e.shouldDelete(ctx, key, obj, existing)
		if delete {
			deleteObj = obj
			return nil, nil, errEmptiedFinalizers
		}
		glog.V(9).Infof("Will update 4")
		return obj, nil, nil
	})

	if err != nil {
		// delete the object
		if err == errEmptiedFinalizers {
			return e.deleteForEmptyFinalizers(ctx, name, key, deleteObj, storagePreconditions)
		}
		if creating {
			err = storeerr.InterpretCreateError(err, e.QualifiedResource, name)
			err = rest.CheckGeneratedNameError(e.CreateStrategy, err, creatingObj)
		} else {
			err = storeerr.InterpretUpdateError(err, e.QualifiedResource, name)
		}
		return nil, false, err
	}
	if creating {
		if e.AfterCreate != nil {
			if err := e.AfterCreate(out); err != nil {
				return nil, false, err
			}
		}
	} else {
		if e.AfterUpdate != nil {
			if err := e.AfterUpdate(out); err != nil {
				return nil, false, err
			}
		}
	}
	if e.Decorator != nil {
		if err := e.Decorator(out); err != nil {
			return nil, false, err
		}
	}
	return out, creating, nil
}

// Get retrieves the item from storage.
func (e *Store) Get(ctx api.Context, name string) (runtime.Object, error) {
	obj := e.NewFunc()
	key, err := e.KeyFunc(ctx, name)
	if err != nil {
		return nil, err
	}
	if err := e.Storage.Get(ctx, key, obj, false); err != nil {
		return nil, storeerr.InterpretGetError(err, e.QualifiedResource, name)
	}
	if e.Decorator != nil {
		if err := e.Decorator(obj); err != nil {
			return nil, err
		}
	}
	return obj, nil
}

var (
	errAlreadyDeleting   = fmt.Errorf("abort delete")
	errDeleteNow         = fmt.Errorf("delete now")
	errEmptiedFinalizers = fmt.Errorf("emptied finalizers")
	errNotSupportMethod  = fmt.Errorf("not support watch method")
)

// shouldUpdateFinalizers returns if we need to update the finalizers of the
// object, and the desired list of finalizers.
// When deciding whether to add the OrphanDependent finalizer, factors in the
// order of highest to lowest priority are: options.OrphanDependents, existing
// finalizers of the object, e.DeleteStrategy.DefaultGarbageCollectionPolicy.
func shouldUpdateFinalizers(e *Store, accessor meta.Object, options *api.DeleteOptions) (shouldUpdate bool, newFinalizers []string) {
	shouldOrphan := false
	// Get default orphan policy from this REST object type
	if gcStrategy, ok := e.DeleteStrategy.(rest.GarbageCollectionDeleteStrategy); ok {
		if gcStrategy.DefaultGarbageCollectionPolicy() == rest.OrphanDependents {
			shouldOrphan = true
		}
	}
	// If a finalizer is set in the object, it overrides the default
	hasOrphanFinalizer := false
	finalizers := accessor.GetFinalizers()
	for _, f := range finalizers {
		if f == api.FinalizerOrphan {
			shouldOrphan = true
			hasOrphanFinalizer = true
			break
		}
		// TODO: update this when we add a finalizer indicating a preference for the other behavior
	}
	// If an explicit policy was set at deletion time, that overrides both
	if options != nil && options.OrphanDependents != nil {
		shouldOrphan = *options.OrphanDependents
	}
	if shouldOrphan && !hasOrphanFinalizer {
		finalizers = append(finalizers, api.FinalizerOrphan)
		return true, finalizers
	}
	if !shouldOrphan && hasOrphanFinalizer {
		var newFinalizers []string
		for _, f := range finalizers {
			if f == api.FinalizerOrphan {
				continue
			}
			newFinalizers = append(newFinalizers, f)
		}
		return true, newFinalizers
	}
	return false, finalizers
}

// Delete removes the item from storage.
func (e *Store) Delete(ctx api.Context, name string, options *api.DeleteOptions) (runtime.Object, error) {
	key, err := e.KeyFunc(ctx, name)
	if err != nil {
		return nil, err
	}

	obj := e.NewFunc()
	if err := e.Storage.Get(ctx, key, obj, false); err != nil {
		return nil, storeerr.InterpretDeleteError(err, e.QualifiedResource, name)
	}
	// support older consumers of delete by treating "nil" as delete immediately
	if options == nil {
		options = api.NewDeleteOptions(0)
	}
	var preconditions storage.Preconditions
	if options.Preconditions != nil {
		preconditions.UID = options.Preconditions.UID
	}
	graceful, pendingGraceful, err := rest.BeforeDelete(e.DeleteStrategy, ctx, obj, options)
	if err != nil {
		return nil, err
	}
	// this means finalizers cannot be updated via DeleteOptions if a deletion is already pending
	if pendingGraceful {
		return e.finalizeDelete(obj, false)
	}
	// check if obj has pending finalizers
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, apierr.NewInternalError(err)
	}
	pendingFinalizers := len(accessor.GetFinalizers()) != 0
	var ignoreNotFound bool
	var deleteImmediately bool = true
	var lastExisting, out runtime.Object
	if !e.EnableGarbageCollection {
		// TODO: remove the check on graceful, because we support no-op updates now.
		if graceful {
			//err, ignoreNotFound, deleteImmediately, out, lastExisting = e.updateForGracefulDeletion(ctx, name, key, options, preconditions, obj)
		}
	} else {
		shouldUpdateFinalizers, _ := shouldUpdateFinalizers(e, accessor, options)
		// TODO: remove the check, because we support no-op updates now.
		if graceful || pendingFinalizers || shouldUpdateFinalizers {
			//err, ignoreNotFound, deleteImmediately, out, lastExisting = e.updateForGracefulDeletionAndFinalizers(ctx, name, key, options, preconditions, obj)
		}
	}
	// !deleteImmediately covers all cases where err != nil. We keep both to be future-proof.
	if !deleteImmediately || err != nil {
		return out, err
	}

	// delete immediately, or no graceful deletion supported
	glog.V(6).Infof("going to delete %s from regitry: ", name)
	out = e.NewFunc()
	if err := e.Storage.Delete(ctx, key, out, nil); err != nil {
		// Please refer to the place where we set ignoreNotFound for the reason
		// why we ignore the NotFound error .
		if storage.IsNotFound(err) && ignoreNotFound && lastExisting != nil {
			// The lastExisting object may not be the last state of the object
			// before its deletion, but it's the best approximation.
			return e.finalizeDelete(lastExisting, true)
		}
		return nil, storeerr.InterpretDeleteError(err, e.QualifiedResource, name)
	}
	return e.finalizeDelete(out, true)
}

// DeleteCollection remove all items returned by List with a given ListOptions from storage.
//
// DeleteCollection is currently NOT atomic. It can happen that only subset of objects
// will be deleted from storage, and then an error will be returned.
// In case of success, the list of deleted objects will be returned.
//
// TODO: Currently, there is no easy way to remove 'directory' entry from storage (if we
// are removing all objects of a given type) with the current API (it's technically
// possibly with storage API, but watch is not delivered correctly then).
// It will be possible to fix it with v3 etcd API.
func (e *Store) DeleteCollection(ctx api.Context, options *api.DeleteOptions, listOptions *api.ListOptions) (runtime.Object, error) {
	listObj, err := e.List(ctx, listOptions)
	if err != nil {
		return nil, err
	}
	items, err := meta.ExtractList(listObj)
	if err != nil {
		return nil, err
	}
	// Spawn a number of goroutines, so that we can issue requests to storage
	// in parallel to speed up deletion.
	// TODO: Make this proportional to the number of items to delete, up to
	// DeleteCollectionWorkers (it doesn't make much sense to spawn 16
	// workers to delete 10 items).
	workersNumber := e.DeleteCollectionWorkers
	if workersNumber < 1 {
		workersNumber = 1
	}
	wg := sync.WaitGroup{}
	toProcess := make(chan int, 2*workersNumber)
	errs := make(chan error, workersNumber+1)

	go func() {
		defer utilruntime.HandleCrash(func(panicReason interface{}) {
			errs <- fmt.Errorf("DeleteCollection distributor panicked: %v", panicReason)
		})
		for i := 0; i < len(items); i++ {
			toProcess <- i
		}
		close(toProcess)
	}()

	wg.Add(workersNumber)
	for i := 0; i < workersNumber; i++ {
		go func() {
			// panics don't cross goroutine boundaries
			defer utilruntime.HandleCrash(func(panicReason interface{}) {
				errs <- fmt.Errorf("DeleteCollection goroutine panicked: %v", panicReason)
			})
			defer wg.Done()

			for {
				index, ok := <-toProcess
				if !ok {
					return
				}
				accessor, err := meta.Accessor(items[index])
				if err != nil {
					errs <- err
					return
				}
				if _, err := e.Delete(ctx, accessor.GetName(), options); err != nil && !apierr.IsNotFound(err) {
					glog.V(4).Infof("Delete %s in DeleteCollection failed: %v", accessor.GetName(), err)
					errs <- err
					return
				}
			}
		}()
	}
	wg.Wait()
	select {
	case err := <-errs:
		return nil, err
	default:
		return listObj, nil
	}
}

func (e *Store) finalizeDelete(obj runtime.Object, runHooks bool) (runtime.Object, error) {
	if runHooks && e.AfterDelete != nil {
		if err := e.AfterDelete(obj); err != nil {
			return nil, err
		}
	}
	if e.ReturnDeletedObject {
		if e.Decorator != nil {
			if err := e.Decorator(obj); err != nil {
				return nil, err
			}
		}
		return obj, nil
	}
	return &unversioned.Status{Status: unversioned.StatusSuccess}, nil
}

// Watch makes a matcher for the given label and field, and calls
// WatchPredicate. If possible, you should customize PredicateFunc to produre a
// matcher that matches by key. SelectionPredicate does this for you
// automatically.
// for mysql we does't support resource wath feature
func (e *Store) Watch(ctx api.Context, options *api.ListOptions) (watch.Interface, error) {
	return nil, errNotSupportMethod
}

func exportObjectMeta(accessor meta.Object, exact bool) {
	accessor.SetUID("")
	if !exact {
		accessor.SetNamespace("")
	}
	accessor.SetCreationTimestamp(unversioned.Time{})
	accessor.SetDeletionTimestamp(nil)
	accessor.SetResourceVersion("")
	accessor.SetSelfLink("")
	if len(accessor.GetGenerateName()) > 0 && !exact {
		accessor.SetName("")
	}
}

// Implements the rest.Exporter interface
func (e *Store) Export(ctx api.Context, name string, opts unversioned.ExportOptions) (runtime.Object, error) {
	obj, err := e.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if accessor, err := meta.Accessor(obj); err == nil {
		exportObjectMeta(accessor, opts.Exact)
	} else {
		glog.V(4).Infof("Object of type %v does not have ObjectMeta: %v", reflect.TypeOf(obj), err)
	}

	if e.ExportStrategy != nil {
		if err = e.ExportStrategy.Export(ctx, obj, opts.Exact); err != nil {
			return nil, err
		}
	} else {
		e.CreateStrategy.PrepareForCreate(ctx, obj)
	}
	return obj, nil
}
