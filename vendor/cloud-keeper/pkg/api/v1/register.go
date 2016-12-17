package v1

import (
	"gofreezer/pkg/api/unversioned"
	. "gofreezer/pkg/api/v1"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/runtime/schema"
	versionedwatch "gofreezer/pkg/watch/versioned"
)

// GroupName is the group name use in this package
const GroupName = ""

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes, addConversionFuncs)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Adds the list of known types to api.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&NodeUser{},
		&NodeUserList{},
		&Node{},
		&NodeList{},
		&APIServer{},
		&APIServerList{},
		&UserService{},
		&UserServiceList{},
		&UserServiceBindingNodes{},
		&Login{},
		&LoginList{},
		&AccServer{},
		&AccServerList{},
		&Account{},
		&AccountList{},
		&AccountInfo{},
		&AccExec{},
		&AccSSHKey{},
		//&AccSSHKeyList{},
		&User{},
		&UserList{},
		&ActiveAPINode{},
		&ActiveAPINodeList{},
		&UserToken{},
		&UserTokenList{},
		&UserPublicFile{},
		&UserPublicFileList{},

		&BatchUsers{},
		&BatchShadowss{},

		&DeleteOptions{},
		&ExportOptions{},
		&ListOptions{},
	)

	// Add common types
	scheme.AddKnownTypes(SchemeGroupVersion, &unversioned.Status{})

	// Add the watch version that applies
	versionedwatch.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
