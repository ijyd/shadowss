package userctl

import (
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/storage"

	"cloud-keeper/pkg/api"

	"github.com/golang/glog"
)

type UserNodeFilter struct {
	nodeName string
}

func NewUserNodeFilter(nodeName string) *UserNodeFilter {
	return &UserNodeFilter{
		nodeName: nodeName,
	}
}

func (e UserNodeFilter) Filter(obj runtime.Object) bool {
	kind := obj.GetObjectKind()

	if "UserService" == kind.GroupVersionKind().Kind {
		glog.V(5).Infof("ignore  user in a node\r\n")
		return false
	}

	user := obj.(*api.UserService)
	_, ok := user.Spec.NodeUserReference[e.nodeName]
	return ok
}

func (e UserNodeFilter) Trigger() []storage.MatchValue {
	return nil
}
