package api

import (
	"reflect"

	"gofreezer/pkg/api/resource"
	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/conversion"
	"gofreezer/pkg/fields"
	"gofreezer/pkg/labels"
	"gofreezer/pkg/types"
	"gofreezer/pkg/util/sets"

	"github.com/davecgh/go-spew/spew"
)

// Conversion error conveniently packages up errors in conversions.
type ConversionError struct {
	In, Out interface{}
	Message string
}

// Return a helpful string about the error
func (c *ConversionError) Error() string {
	return spew.Sprintf(
		"Conversion error: %s. (in: %v(%+v) out: %v)",
		c.Message, reflect.TypeOf(c.In), c.In, reflect.TypeOf(c.Out),
	)
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

var standardFinalizers = sets.NewString(
	string(FinalizerKubernetes),
	FinalizerOrphan,
)

func IsStandardFinalizerName(str string) bool {
	return standardFinalizers.Has(str)
}

// NewDeleteOptions returns a DeleteOptions indicating the resource should
// be deleted within the specified grace period. Use zero to indicate
// immediate deletion. If you would prefer to use the default grace period,
// use &api.DeleteOptions{} directly.
func NewDeleteOptions(grace int64) *DeleteOptions {
	return &DeleteOptions{GracePeriodSeconds: &grace}
}

// NewPreconditionDeleteOptions returns a DeleteOptions with a UID precondition set.
func NewPreconditionDeleteOptions(uid string) *DeleteOptions {
	u := types.UID(uid)
	p := Preconditions{UID: &u}
	return &DeleteOptions{Preconditions: &p}
}

// NewUIDPreconditions returns a Preconditions with UID set.
func NewUIDPreconditions(uid string) *Preconditions {
	u := types.UID(uid)
	return &Preconditions{UID: &u}
}
