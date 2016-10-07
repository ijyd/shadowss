package apiserverctl

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/backend"
	"cloud-keeper/pkg/etcdhelper"
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/storage"
	"golib/pkg/util/network"
	"time"

	"github.com/golang/glog"
)

const (
	PrefixAPIServer = "/" + "APIServer"
)

func GetAPIServerName() (string, error) {
	return network.ExternalMAC()
}

func AddLocalAPIServer(db *backend.Backend, helper *etcdhelper.EtcdHelper, localhost string, host []string, port int, ttl uint64, etcd bool, mysql bool) (runtime.Object, error) {

	APIName, err := GetAPIServerName()
	if err != nil {
		return nil, err
	}

	spec := api.APIServerSpec{
		Server: api.APIServerInfor{
			Name:       APIName,
			Host:       localhost,
			Status:     true,
			Port:       int64(port),
			CreateTime: time.Now(),
		},
		HostList: host,
	}

	apisrv := &api.APIServer{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "APIServer",
			APIVersion: "v1",
		},
		ObjectMeta: prototype.ObjectMeta{
			Name: spec.Server.Name,
		},
		Spec: spec,
	}

	//del db api server first
	DelAPIServer(db, nil, APIName, false, true)

	return AddAPIServer(db, helper, apisrv, ttl, etcd, mysql)
}

func AddAPIServer(db *backend.Backend, helper *etcdhelper.EtcdHelper, srv *api.APIServer, ttl uint64, etcd bool, mysql bool) (runtime.Object, error) {
	if mysql {
		err := db.CreateAPIServer(srv.Spec.Server)
		if err != nil {
			return nil, err
		}
	}

	if etcd {
		ctx := prototype.NewContext()
		outItem := new(api.APIServer)
		err := helper.StorageCodec.Storage.Create(ctx, PrefixAPIServer+"/"+srv.Name, srv, outItem, ttl)
		glog.V(5).Infof("Add apiserver %v err %v\r\n", srv, err)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func CheckLocalAPIServer(helper *etcdhelper.EtcdHelper) bool {
	APIName, err := GetAPIServerName()
	if err != nil {
		glog.Errorf("Got external mac error %v \r\n", err)
		return false
	}

	objs, err := GetAPIServers(helper)
	glog.Infof("got obj %v err %v\r\n", objs, err)

	obj, err := GetAPIServer(helper, PrefixAPIServer+"/"+APIName)
	srv := obj.(*api.APIServer)
	if err == nil && len(srv.Name) > 0 {
		return true
	} else {
		return false
	}
}

func GetAPIServer(helper *etcdhelper.EtcdHelper, key string) (runtime.Object, error) {
	ctx := prototype.NewContext()
	outItem := new(api.APIServer)
	err := helper.StorageCodec.Storage.Get(ctx, key, outItem, true)
	return outItem, err
}

func GetAPIServers(helper *etcdhelper.EtcdHelper) (runtime.Object, error) {
	ctx := prototype.NewContext()
	outItem := new(api.APIServerList)

	options := &prototype.ListOptions{ResourceVersion: "0"}
	err := helper.StorageCodec.Storage.List(ctx, PrefixAPIServer, options.ResourceVersion, storage.Everything, outItem)
	return outItem, err
}

func DelAPIServer(db *backend.Backend, helper *etcdhelper.EtcdHelper, name string, etcd bool, mysql bool) error {

	if etcd {
		ctx := prototype.NewContext()
		outItem := new(api.APIServer)
		//it is a strik we use node name for key

		key := PrefixAPIServer + "/" + name
		err := helper.StorageCodec.Storage.Delete(ctx, key, outItem, nil)
		if err != nil {
			glog.Errorf("Create node config err %v items %v\r\n", err, outItem)
			return err
		}
	}

	if mysql {
		err := db.DeleteAPIServerByName(name)
		if err != nil {
			return err
		}
	}

	return nil
}
