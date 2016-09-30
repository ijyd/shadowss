package api

import "gofreezer/pkg/api/unversioned"

var AccountInfoType = unversioned.TypeMeta{
	Kind:       "AccountInfo",
	APIVersion: "v1",
}

var AccServerType = unversioned.TypeMeta{
	Kind:       "AccServer",
	APIVersion: "v1",
}

var AccServerSSHKeyType = unversioned.TypeMeta{
	Kind:       "AccServerSSHKey",
	APIVersion: "v1",
}
