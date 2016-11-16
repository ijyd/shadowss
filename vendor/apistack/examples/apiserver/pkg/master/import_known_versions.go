package master

// These imports are the API groups the API server will support.
import (
	"fmt"

	_ "apistack/examples/apiserver/pkg/api/install"
	_ "apistack/examples/apiserver/pkg/apis/testgroup/install"

	"apistack/pkg/apimachinery/registered"
)

func init() {
	if missingVersions := registered.ValidateEnvRequestedVersions(); len(missingVersions) != 0 {
		panic(fmt.Sprintf("KUBE_API_VERSIONS contains versions that are not installed: %q.", missingVersions))
	}
}
