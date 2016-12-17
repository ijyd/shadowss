package v1beta1

import (
	"apistack/pkg/apis/abac/v1beta1"
	"gofreezer/pkg/api/unversioned"
)

// +genclient=true

// PolicyList contains a list ABAC policy rule
type Policy struct {
	v1beta1.Policy
}

// PolicyList contains a list ABAC policy rule
type PolicyList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []v1beta1.Policy `json:"item,omitempty"`
}
