package nodectl

import (
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/storage"

	"github.com/golang/glog"
)

type NodeFilter struct {
}

func (e NodeFilter) Filter(obj runtime.Object) bool {
	kind := obj.GetObjectKind()

	if "NodeUser" == kind.GroupVersionKind().Kind {
		glog.V(5).Infof("ignore node user in a node\r\n")
		return false
	}

	return true
}

func (e NodeFilter) Trigger() []storage.MatchValue {
	return nil
}
