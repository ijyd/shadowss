package testgroup

import (
	api "gofreezer/pkg/api"
	"gofreezer/pkg/api/unversioned"
)

// +genclient=true

type TestType struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata,omitempty"`
	Status               TestTypeStatus `json:"status,omitempty"`
}

type TestTypeList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []TestType `json:"items"`
}

type TestTypeStatus struct {
	Blah string
}
