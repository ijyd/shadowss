// Package install installs the experimental API group, making it available as
// an option to all of the API encoding/decoding machinery.
package install

import (
	"apistack/pkg/apimachinery/announced"
	"cloud-keeper/pkg/apis/abacpolicys"
	"cloud-keeper/pkg/apis/abacpolicys/v1beta1"

	"gofreezer/pkg/util/sets"
)

func init() {
	if err := announced.NewGroupMetaFactory(
		&announced.GroupMetaFactoryArgs{
			GroupName:                  abacpolicys.GroupName,
			VersionPreferenceOrder:     []string{v1beta1.SchemeGroupVersion.Version},
			ImportPrefix:               "cloud-keeper/pkg/apis/abacpolicys",
			AddInternalObjectsToScheme: abacpolicys.AddToScheme,
			RootScopedKinds: sets.NewString(
				"Policy",
			),
			IgnoredKinds: sets.NewString(
				"ListOptions",
				"Status",
			),
		},
		announced.VersionToSchemeFunc{
			v1beta1.SchemeGroupVersion.Version: v1beta1.AddToScheme,
		},
	).Announce().RegisterAndEnable(); err != nil {
		panic(err)
	}
}
