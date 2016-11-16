package helper

import (
	"fmt"
	"reflect"
	"strings"

	"gofreezer/pkg/api/meta"
	"gofreezer/pkg/conversion"
	"gofreezer/pkg/runtime"

	"github.com/golang/glog"
)

func CloneRuntimeObj(objPtr runtime.Object) (runtime.Object, error) {
	v, err := conversion.EnforcePtr(objPtr)
	if err != nil {
		return nil, err
	}

	newObj := reflect.New(v.Type())
	storeObj := newObj.Interface().(runtime.Object)
	return storeObj, nil
}

func GetObjKind(objPtr runtime.Object) string {
	v, err := conversion.EnforcePtr(objPtr)
	if err != nil {
		return string("")
	}

	kind := v.Type().String()
	if i := strings.IndexAny(kind, "."); i >= 0 {
		kind = kind[i+1:]
	}
	return kind
}

func GetListItemObj(listObj runtime.Object) (listPtr interface{}, itemObj runtime.Object, err error) {
	listPtr, err = meta.GetItemsPtr(listObj)
	if err != nil {
		return
	}

	items, err := conversion.EnforcePtr(listPtr)
	if err != nil {
		return
	}
	if items.Kind() != reflect.Slice {
		err = fmt.Errorf("object(%v) not a slice", items.Kind())
		return
	}

	itemObj = reflect.New(items.Type().Elem()).Interface().(runtime.Object)
	glog.V(5).Infof("get list object %+v", itemObj)

	return
}
