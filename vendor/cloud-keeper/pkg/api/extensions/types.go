package extensions

import (
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/api/unversioned"
)

type Endpoint struct {
	URL string `json:"url,omitempty"`
}

type VPSKeeperSpec struct {
	Endpoint []Endpoint `json:"endpoint,omitempty"`
}

type VPSKeeper struct {
	unversioned.TypeMeta `json:",inline"`
	prototype.ObjectMeta `json:"metadata,omitempty"`

	Spec VPSKeeperSpec `json:"spec,omitempty"`
}

type VPSKeeperList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []VPSKeeper `json:"items"`
}
