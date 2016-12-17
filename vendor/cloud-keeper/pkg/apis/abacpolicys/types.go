package abacpolicys

import (
	"apistack/pkg/apis/abac"
	"gofreezer/pkg/api/unversioned"
)

// +genclient=true

// PolicyList contains a list ABAC policy rule
type Policy struct {
	abac.Policy
}

// PolicyList contains a list ABAC policy rule
type PolicyList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []abac.Policy `json:"item,omitempty"`
}
