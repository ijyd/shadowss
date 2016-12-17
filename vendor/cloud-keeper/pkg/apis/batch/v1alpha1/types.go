package v1alpha1

import (
	"cloud-keeper/pkg/api/v1"
	"gofreezer/pkg/api/unversioned"
	freezerv1 "gofreezer/pkg/api/v1"
)

// +genclient=true

type BatchAccServerSpec struct {
	ServerList []v1.AccServer `json:"serverList,omitempty"`
}

type BatchAccServer struct {
	unversioned.TypeMeta `json:",inline"`
	freezerv1.ObjectMeta `json:"metadata,omitempty"`

	Spec BatchAccServerSpec `json:"spec,omitempty"`
}

type BatchResumeUsersSpec struct {
	SchedulingTime unversioned.Time `json:"schedulingTime,omitempty"`
}

type BatchResumeUsers struct {
	unversioned.TypeMeta `json:",inline"`
	freezerv1.ObjectMeta `json:"metadata,omitempty"`

	Spec BatchResumeUsersSpec `json:"spec,omitempty"`
}
