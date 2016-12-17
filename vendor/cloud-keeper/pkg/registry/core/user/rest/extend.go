package rest

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/user"
	"fmt"
	"time"

	"github.com/golang/glog"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/errors"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/runtime"
)

func NewExtendREST(user user.Registry) (*BindingNodeREST, *PropertiesREST, *ActivationREST) {
	return &BindingNodeREST{user}, &PropertiesREST{user: user}, &ActivationREST{user}
}

type BindingNodeREST struct {
	user user.Registry
}

func (*BindingNodeREST) New() runtime.Object {
	return &api.UserService{}
}

func (*BindingNodeREST) NewList() runtime.Object {
	return &api.UserServiceList{}
}

func (r *BindingNodeREST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	if name, ok := options.FieldSelector.RequiresExactMatch("metadata.name"); ok {
		user, err := r.user.GetUser(ctx, name)
		if err != nil {
			return nil, err
		}

		bindingnodes := &api.UserServiceBindingNodes{}
		bindingnodes.Name = user.Name
		bindingnodes.Annotations = make(map[string]string)
		for k, v := range user.Annotations {
			bindingnodes.Annotations[k] = v
		}
		bindingnodes.Spec.NodeCnt = user.Spec.UserService.NodeCnt
		bindingnodes.Spec.Status = user.Spec.UserService.Status

		bindingnodes.Spec.NodeUserReference = make(map[string]api.NodeReferences)
		for nodeName, v := range user.Spec.UserService.Nodes {
			nodeRefer := api.NodeReferences{}

			userRefer := &nodeRefer.User
			userRefer.Name = user.Name
			userRefer.DownloadTraffic = v.User.DownloadTraffic
			userRefer.UploadTraffic = v.User.UploadTraffic
			userRefer.ID = v.User.ID
			userRefer.Name = v.User.Name
			userRefer.EnableOTA = v.User.EnableOTA
			userRefer.Method = v.User.Method
			userRefer.Port = v.User.Port
			userRefer.Password = v.User.Password

			nodeRefer.Host = v.Host

			bindingnodes.Spec.NodeUserReference[nodeName] = nodeRefer
		}

		return bindingnodes, err
	} else {
		return nil, errors.NewBadRequest("need a 'metadata.name' filed selector")
	}

}

func (r *BindingNodeREST) Create(ctx freezerapi.Context, obj runtime.Object) (runtime.Object, error) {
	userservice := obj.(*api.UserService)

	user, err := r.user.GetUser(ctx, userservice.Name)
	if err != nil {
		return nil, err
	}

	err = r.user.AddNodeToUser(ctx, user, userservice, true, true)
	if err != nil {
		return nil, err
	}

	return userservice, err
}

func (r *BindingNodeREST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	obj, err := r.user.GetUser(ctx, name)
	if err != nil {
		return nil, false, err
	}

	newobj, err := objInfo.UpdatedObject(ctx, nil)
	if err != nil {
		return nil, false, err
	}
	//only update nodes
	newUserService := newobj.(*api.UserService)

	if newUserService.Spec.Delete {
		nodeName := newUserService.Spec.NodeName
		err = r.user.DelNodeFromUser(ctx, obj, nodeName, true, true)
		if err != nil {
			return nil, false, err
		}
		return newUserService, true, nil
	} else {
		return nil, false, errors.NewInternalError(fmt.Errorf("not support update for userservice,only delete"))
	}

}

type PropertiesREST struct {
	user user.Registry
}

func (*PropertiesREST) New() runtime.Object {
	return &api.User{}
}

func (r *PropertiesREST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	obj, err := r.user.GetUser(ctx, name)
	if err != nil {
		return nil, false, err
	}
	newobj, err := objInfo.UpdatedObject(ctx, obj)
	if err != nil {
		return nil, false, err
	}

	//only update annotations
	user := newobj.(*api.User)
	obj.Annotations = make(map[string]string)
	for k, v := range user.Annotations {
		if len(v) > 0 {
			obj.Annotations[k] = v
		}
	}
	userAnnotationRefreshTime := "refreshTime"
	obj.Annotations[userAnnotationRefreshTime] = time.Now().String()

	glog.V(5).Infof("refresh user properties, will force delete user %+v from old node \r\n", *obj)
	err = r.user.DelAllNodeFromUser(ctx, obj, false, true)
	if err != nil {
		return nil, false, err
	}

	glog.V(5).Infof("will force reinitnodeuser with %+v\r\n", *obj)
	err = r.user.ReInitNodeToUser(ctx, obj)
	if err != nil {
		return nil, false, err
	}

	glog.V(5).Infof("got new node for update user %+v\r\n", *obj)

	_, err = r.user.UpdateUser(ctx, obj)
	if err != nil {
		return nil, false, err
	}

	return obj, true, err
}

type ActivationREST struct {
	user user.Registry
}

func (*ActivationREST) New() runtime.Object {
	return &api.User{}
}

func (r *ActivationREST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	obj, err := r.user.GetUser(ctx, name)
	if err != nil {
		return nil, false, err
	}
	newobj, err := objInfo.UpdatedObject(ctx, obj)
	if err != nil {
		return nil, false, err
	}

	newUser := newobj.(*api.User)

	switch obj.Spec.DetailInfo.UserType {
	case api.UserTypeDesktopRouter:
		//user not activation, allow user to activation
		if !obj.Spec.DetailInfo.Activation {
			//modify user expire_time
			expire := time.Now().Add(time.Duration(time.Hour * 24 * 365))
			obj.Spec.DetailInfo.ExpireTime = unversioned.NewTime(expire)
			obj.Spec.DetailInfo.Activation = newUser.Spec.DetailInfo.Activation
		}
	default:
		return nil, false, fmt.Errorf("not support user type %v", obj.Spec.DetailInfo.UserType)
	}

	_, err = r.user.UpdateUser(ctx, obj)
	if err != nil {
		return nil, false, err
	}

	return obj, true, err
}
