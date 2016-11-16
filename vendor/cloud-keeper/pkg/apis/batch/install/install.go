// Package install installs the experimental API group, making it available as
// an option to all of the API encoding/decoding machinery.
package install

import (
	"apistack/pkg/apimachinery/announced"
	"cloud-keeper/pkg/apis/batch"
	"cloud-keeper/pkg/apis/batch/v1alpha1"

	"gofreezer/pkg/util/sets"
)

func init() {
	if err := announced.NewGroupMetaFactory(
		&announced.GroupMetaFactoryArgs{
			GroupName:                  batch.GroupName,
			VersionPreferenceOrder:     []string{v1alpha1.SchemeGroupVersion.Version},
			ImportPrefix:               "cloud-keeper/pkg/apis/batch",
			AddInternalObjectsToScheme: batch.AddToScheme,
			RootScopedKinds: sets.NewString(
				"BatchAccServer",
			),
			IgnoredKinds: sets.NewString(
				"ListOptions",
				"DeleteOptions",
				"Status",
			),
		},
		announced.VersionToSchemeFunc{
			v1alpha1.SchemeGroupVersion.Version: v1alpha1.AddToScheme,
		},
	).Announce().RegisterAndEnable(); err != nil {
		panic(err)
	}
}
