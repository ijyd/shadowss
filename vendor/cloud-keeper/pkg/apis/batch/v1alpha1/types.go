package v1alpha1

import (
	"cloud-keeper/pkg/api/v1"
	"gofreezer/pkg/api/unversioned"
)

// +genclient=true

type BatchAccServerSpec struct {
	ServerList []v1.AccServer `json:"serverList,omitempty"`
}

type BatchAccServer struct {
	unversioned.TypeMeta `json:",inline"`
	v1.ObjectMeta        `json:"metadata,omitempty"`

	Spec BatchAccServerSpec `json:"spec,omitempty"`
}
