package storage

import (
	"gofreezer/pkg/conversion"
	"gofreezer/pkg/runtime"
	"strings"
)

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
