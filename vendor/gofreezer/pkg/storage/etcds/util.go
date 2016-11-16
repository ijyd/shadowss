package etcds

import (
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/storage"
)

type SimpleUpdateFunc func(runtime.Object) (runtime.Object, error)

// SimpleUpdateFunc converts SimpleUpdateFunc into UpdateFunc
func SimpleUpdate(fn SimpleUpdateFunc) UpdateFunc {
	return func(input runtime.Object, _ storage.ResponseMeta) (runtime.Object, *uint64, error) {
		out, err := fn(input)
		return out, nil, err
	}
}
