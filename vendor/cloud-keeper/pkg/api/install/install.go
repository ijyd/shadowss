// Package install installs the v1 monolithic api, making it available as an
// option to all of the API encoding/decoding machinery.
package install

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/api/v1"

	"github.com/golang/glog"

	"gofreezer/pkg/api/meta"
	"gofreezer/pkg/api/unversioned"
)

// const importPrefix = "gofreezer/examples/etcd/app/api"

var accessor = meta.NewAccessor()

// availableVersions lists all known external versions for this group from most preferred to least preferred
var availableVersions = []unversioned.GroupVersion{v1.SchemeGroupVersion}

func init() {
	//registered.RegisterVersions(availableVersions)
	externalVersions := []unversioned.GroupVersion{}
	for _, v := range availableVersions {
		//if registered.IsAllowedVersion(v) {
		externalVersions = append(externalVersions, v)
		//}
	}
	if len(externalVersions) == 0 {
		glog.V(4).Infof("No version is registered for group %v", api.GroupName)
		return
	}

	if err := enableVersions(externalVersions); err != nil {
		glog.V(4).Infof("%v", err)
		return
	}
}

// TODO: enableVersions should be centralized rather than spread in each API
// group.
// We can combine registered.RegisterVersions, registered.EnableVersions and
// registered.RegisterGroup once we have moved enableVersions there.
func enableVersions(externalVersions []unversioned.GroupVersion) error {
	addVersionsToScheme(externalVersions...)
	return nil
}

func addVersionsToScheme(externalVersions ...unversioned.GroupVersion) {
	// add the internal version to Scheme
	if err := api.AddToScheme(); err != nil {
		// Programmer error, detect immediately
		panic(err)
	}
	// add the enabled external versions to Scheme
	for _, v := range externalVersions {
		// if !registered.IsEnabledVersion(v) {
		// 	glog.Errorf("Version %s is not enabled, so it will not be added to the Scheme.", v)
		// 	continue
		// }
		switch v {
		case v1.SchemeGroupVersion:
			if err := v1.AddToScheme(api.Scheme); err != nil {
				// Programmer error, detect immediately
				panic(err)
			}
		}
	}
}
