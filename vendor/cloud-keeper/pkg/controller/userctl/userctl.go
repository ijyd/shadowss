package userctl

import (
	"cloud-keeper/pkg/api"
	"fmt"

	"cloud-keeper/pkg/etcdhelper"
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/storage"

	"github.com/golang/glog"
)

const (
	PrefixUserService = "/" + "UserService"
)

func AddUserServiceHelper(helper *etcdhelper.EtcdHelper, username string, nodeUserRefer map[string]api.NodeReferences) error {

	spec := api.UserServiceSpec{
		NodeUserReference: nodeUserRefer,
		NodeCnt:           uint(len(nodeUserRefer)),
	}

	srv := &api.UserService{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "UserService",
			APIVersion: "v1",
		},
		ObjectMeta: prototype.ObjectMeta{
			Name: username,
		},
		Spec: spec,
	}

	return AddUserService(helper, srv)
}

//AddUserService key
//like as /api/UserService/{username}
//username is a User resource name in db
func AddUserService(helper *etcdhelper.EtcdHelper, userSrv *api.UserService) error {
	ctx := prototype.NewContext()
	outItem := new(api.UserService)

	key := PrefixUserService + "/" + userSrv.Name
	err := helper.StorageCodec.Storage.Create(ctx, key, userSrv, outItem, 0)
	if err != nil {
		glog.Errorf("Create user service  err %v items %v\r\n", err, outItem)
		return err
	}

	return nil
}

func DelUserService(helper *etcdhelper.EtcdHelper, name string) error {
	ctx := prototype.NewContext()
	outItem := new(api.UserService)

	key := PrefixUserService + "/" + name
	err := helper.StorageCodec.Storage.Delete(ctx, key, outItem, nil)
	if err != nil {
		glog.Errorf("Create node config err %v items %v\r\n", err, outItem)
		return err
	}

	return nil
}

func GetUserService(helper *etcdhelper.EtcdHelper, name string) (runtime.Object, error) {
	ctx := prototype.NewContext()
	outItem := new(api.UserService)

	key := PrefixUserService + "/" + name

	err := helper.StorageCodec.Storage.Get(ctx, key, outItem, true)
	if err != nil {
		return nil, err
	}

	return outItem, nil
}

//UpdateNodeUsers update node user
func UpdateUserAnnotations(helper *etcdhelper.EtcdHelper, name string, annotation map[string]string) (runtime.Object, error) {

	obj, err := GetUserService(helper, name)
	if err != nil {
		glog.Errorf("Get resource error %v\r\n", err)
		return nil, err
	}

	oldobj := obj.(*api.UserService)
	if oldobj.Name == "" {
		errStr := "Get resource error not found this user"
		glog.Errorf("%s\r\n", errStr)
		return nil, fmt.Errorf("%s", errStr)
	}

	ctx := prototype.NewContext()
	outItem := new(api.UserService)
	newres := oldobj

	key := PrefixUserService + "/" + oldobj.Name
	err = helper.StorageCodec.Storage.GuaranteedUpdate(ctx, key, outItem, false, nil, func(existing runtime.Object, res storage.ResponseMeta) (runtime.Object, *uint64, error) {
		glog.Infof("existing obj %+v", existing)

		newres.Annotations = annotation

		return newres, nil, nil
	})
	if err != nil {
		glog.Errorf("update failure %v\r\n", err)
		return nil, err
	}

	return newres, nil
}

//UpdateNodeUsers update node user
func UpdateUserService(helper *etcdhelper.EtcdHelper, userSrv *api.UserService) error {

	obj, err := GetUserService(helper, userSrv.Name)
	if err != nil {
		glog.Errorf("Get resource error %v\r\n", err)
		return err
	}

	oldobj := obj.(*api.UserService)
	if oldobj.Name == "" {
		errStr := "Get resource error not found this user"
		glog.Errorf("%s\r\n", errStr)
		return fmt.Errorf("%s", errStr)
	}

	ctx := prototype.NewContext()
	outItem := new(api.UserService)

	key := PrefixUserService + "/" + oldobj.Name
	err = helper.StorageCodec.Storage.GuaranteedUpdate(ctx, key, outItem, false, nil, func(existing runtime.Object, res storage.ResponseMeta) (runtime.Object, *uint64, error) {
		glog.Infof("existing obj %+v", existing)

		newres := userSrv
		return newres, nil, nil
	})
	if err != nil {
		glog.Errorf("update failure %v\r\n", err)
		return err
	}

	return nil
}
