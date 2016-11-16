// Package install installs the experimental API group, making it available as
// an option to all of the API encoding/decoding machinery.
package install

import (
	"apistack/examples/apiserver/pkg/apis/testgroup"
	"apistack/examples/apiserver/pkg/apis/testgroup/v1"
	"apistack/pkg/apimachinery/announced"

	"gofreezer/pkg/util/sets"
)

func init() {
	if err := announced.NewGroupMetaFactory(
		&announced.GroupMetaFactoryArgs{
			GroupName:                  testgroup.GroupName,
			VersionPreferenceOrder:     []string{v1.SchemeGroupVersion.Version},
			ImportPrefix:               "apistack/examples/apiserver/pkg/apis/testgroup",
			AddInternalObjectsToScheme: testgroup.AddToScheme,
			RootScopedKinds: sets.NewString(
				"TestType",
				"TestTypeList",
			),
			IgnoredKinds: sets.NewString(
				"ListOptions",
				"DeleteOptions",
				"Status",
			),
		},
		announced.VersionToSchemeFunc{
			v1.SchemeGroupVersion.Version: v1.AddToScheme,
		},
	).Announce().RegisterAndEnable(); err != nil {
		panic(err)
	}
}
