package api

import (
	"gofreezer/pkg/api/resource"
	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/conversion"
	"gofreezer/pkg/fields"
	"gofreezer/pkg/labels"
	"gofreezer/pkg/util/sets"
)

var standardFinalizers = sets.NewString(
	string(FinalizerKubernetes),
	FinalizerOrphan,
)

// NewDeleteOptions returns a DeleteOptions indicating the resource should
// be deleted within the specified grace period. Use zero to indicate
// immediate deletion. If you would prefer to use the default grace period,
// use &api.DeleteOptions{} directly.
func NewDeleteOptions(grace int64) *DeleteOptions {
	return &DeleteOptions{GracePeriodSeconds: &grace}
}

func IsStandardFinalizerName(str string) bool {
	return standardFinalizers.Has(str)
}

// Semantic can do semantic deep equality checks for api objects.
// Example: api.Semantic.DeepEqual(aPod, aPodWithNonNilButEmptyMaps) == true
var Semantic = conversion.EqualitiesOrDie(
	func(a, b resource.Quantity) bool {
		// Ignore formatting, only care that numeric value stayed the same.
		// TODO: if we decide it's important, it should be safe to start comparing the format.
		//
		// Uninitialized quantities are equivalent to 0 quantities.
		return a.Cmp(b) == 0
	},
	func(a, b unversioned.Time) bool {
		return a.UTC() == b.UTC()
	},
	func(a, b labels.Selector) bool {
		return a.String() == b.String()
	},
	func(a, b fields.Selector) bool {
		return a.String() == b.String()
	},
)
