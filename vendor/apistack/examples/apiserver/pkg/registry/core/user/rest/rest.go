package rest

import (
	"apistack/examples/apiserver/pkg/api"
	"apistack/examples/apiserver/pkg/registry/core/user/dynamodb"
	"apistack/examples/apiserver/pkg/registry/core/user/etcd"
	"apistack/examples/apiserver/pkg/registry/core/user/mysql"
	"fmt"

	"github.com/golang/glog"

	"gofreezer/pkg/api/errors"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/runtime"

	freezerapi "gofreezer/pkg/api"
)

type UserREST struct {
	etcd     *etcd.REST
	mysql    *mysql.REST
	dynamodb *dynamodb.REST
}

var mysqlbackend = false
var etcdbackend = false
var dynamodbbackend = true

func NewREST(etcdHandler *etcd.REST, mysqlHandler *mysql.REST, dynamodb *dynamodb.REST) *UserREST {
	return &UserREST{
		etcd:     etcdHandler,
		mysql:    mysqlHandler,
		dynamodb: dynamodb,
	}
}

func (*UserREST) New() runtime.Object {
	return &api.User{}
}

func (*UserREST) NewList() runtime.Object {
	return &api.UserList{}
}

func (r *UserREST) Get(ctx freezerapi.Context, name string) (runtime.Object, error) {
	if mysqlbackend {
		obj, err := r.mysql.Get(ctx, name)
		if err != nil {
			glog.Errorf("got from mysql error %v\r\n", err)
			return nil, err
		}
		return obj, err
	}

	if etcdbackend {
		etcdobj, err := r.etcd.Get(ctx, name)
		if err != nil {
			glog.Errorf("got from etcd error %v\r\n", err)
			return nil, err
		}
		return etcdobj, err
	}

	if dynamodbbackend {
		obj, err := r.dynamodb.Get(ctx, name)
		if err != nil {
			glog.Errorf("got from mysql error %v\r\n", err)
			return nil, err
		}
		return obj, err
	}
	return nil, errors.NewInternalError(fmt.Errorf("not enable any backend"))
}

func (r *UserREST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	glog.V(5).Infof("list options %+v", *options)
	if mysqlbackend {
		obj, err := r.mysql.List(ctx, options)
		if err != nil {
			glog.Errorf("got from mysql error %v\r\n", err)
			return nil, err
		}
		return obj, err
	}

	if etcdbackend {
		etcdobj, err := r.etcd.List(ctx, options)
		if err != nil {
			glog.Errorf("got from etcd error %v\r\n", err)
			return nil, err
		}
		return etcdobj, err
	}

	if dynamodbbackend {
		obj, err := r.dynamodb.List(ctx, options)
		if err != nil {
			glog.Errorf("got from dynamodb error %v\r\n", err)
			return nil, err
		}
		return obj, err
	}

	return nil, errors.NewInternalError(fmt.Errorf("not enable any backend"))
}

func (r *UserREST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	if mysqlbackend {
		obj, falg, err := r.mysql.Update(ctx, name, objInfo)
		if err != nil {
			glog.Errorf("got from mysql error %v\r\n", err)
			return nil, falg, err
		}
		return obj, falg, err
	}

	if etcdbackend {
		etcdobj, falg, err := r.etcd.Update(ctx, name, objInfo)
		if err != nil {
			glog.Errorf("got from etcd error %v\r\n", err)
			return nil, falg, err
		}
		return etcdobj, falg, err
	}

	if dynamodbbackend {
		glog.Infof("**************update dynamo db \r\n")
		obj, flag, err := r.dynamodb.Update(ctx, name, objInfo)
		if err != nil {
			glog.Errorf("got from dynamodb error %v\r\n", err)
			return nil, flag, err
		}
		return obj, flag, err
	}

	return nil, false, errors.NewInternalError(fmt.Errorf("not enable any backend"))
}

func (r *UserREST) Delete(ctx freezerapi.Context, name string, options *freezerapi.DeleteOptions) (runtime.Object, error) {
	glog.V(5).Infof("delete %v resoure", name)
	if mysqlbackend {
		obj, err := r.mysql.Delete(ctx, name, options)
		if err != nil {
			glog.Errorf("got from mysql error %v\r\n", err)
			return nil, err
		}
		return obj, err
	}

	if etcdbackend {
		etcdobj, err := r.etcd.Delete(ctx, name, options)
		if err != nil {
			glog.Errorf("got from etcd error %v\r\n", err)
			return nil, err
		}
		return etcdobj, err
	}

	if dynamodbbackend {
		obj, err := r.dynamodb.Delete(ctx, name, options)
		if err != nil {
			glog.Errorf("got from dynamodb error %v\r\n", err)
			return nil, err
		}
		return obj, err
	}

	return nil, errors.NewInternalError(fmt.Errorf("not enable any backend"))
}

func (r *UserREST) Create(ctx freezerapi.Context, obj runtime.Object) (runtime.Object, error) {
	if mysqlbackend {
		resObj, err := r.mysql.Create(ctx, obj)
		if err != nil {
			glog.Errorf("got from mysql error %v\r\n", err)
			return nil, err
		}
		return resObj, err
	}

	if etcdbackend {
		resObj, err := r.etcd.Create(ctx, obj)
		if err != nil {
			glog.Errorf("got from etcd error %v\r\n", err)
			return nil, err
		}
		return resObj, err
	}

	if dynamodbbackend {
		resObj, err := r.dynamodb.Create(ctx, obj)
		if err != nil {
			glog.Errorf("got from dynamodb error %v\r\n", err)
			return nil, err
		}
		return resObj, err
	}

	return nil, errors.NewInternalError(fmt.Errorf("not enable any backend"))
}
