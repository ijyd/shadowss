package api

import (
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/api/unversioned"
)

type LoginSpec struct {
	AuthName string `json:"authname,omitempty"`
	Auth     string `json:"auth,ommitempty"`
	Token    string `json:"token,omitempty"`
}

type Login struct {
	unversioned.TypeMeta `json:",inline"`
	prototype.ObjectMeta `json:"metadata,omitempty"`

	Spec LoginSpec `json:"spec,omitempty"`
}

type LoginList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []Login `json:"items"`
}
