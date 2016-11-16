package batch

import (
	apiext "cloud-keeper/pkg/api"
	api "gofreezer/pkg/api"
	"gofreezer/pkg/api/unversioned"
)

// +genclient=true

type BatchAccServerSpec struct {
	ServerList []apiext.AccServer `json:"serverList,omitempty"`
}

type BatchAccServer struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata,omitempty"`

	Spec BatchAccServerSpec `json:"spec,omitempty"`
}
