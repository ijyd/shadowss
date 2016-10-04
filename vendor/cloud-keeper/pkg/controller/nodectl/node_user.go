package nodectl

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/etcdhelper"
	"fmt"
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/storage"

	"github.com/golang/glog"
)

func BuildNodeUserPrefix(nodeName, userName string) string {
	//key := PrefixNode + "/" + user.Spec.NodeName + PrefixNodeUser + "/" + user.Name
	if len(nodeName) == 0 {
		return string("")
	}
	key := PrefixNode + "/" + nodeName

	if len(userName) != 0 {
		key = key + PrefixNodeUser + "/" + userName
	}

	return key
}

func AddNodeUserHelper(helper *etcdhelper.EtcdHelper, nodeName string, userRefer api.UserReferences, phase api.NodeUserPhase, ttl uint64) error {
	if nodeName == "" {
		return fmt.Errorf("invalid node name")
	}

	if userRefer.Name == "" {
		return fmt.Errorf("invalid user name")
	}

	nodeUser := &api.NodeUser{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "NodeUser",
			APIVersion: "v1",
		},
		ObjectMeta: prototype.ObjectMeta{
			Name: userRefer.Name,
		},
		Spec: api.NodeUserSpec{
			User:     userRefer,
			NodeName: nodeName,
			Phase:    phase,
		},
	}

	glog.V(5).Infof("app node user %+v refer:%v\r\n", *nodeUser, userRefer)

	return AddNodeUser(helper, nodeUser, ttl)
}

//AddNodeUsers update node user
//like as /api/node/{nodename}/nodeuser/{username}
func AddNodeUser(helper *etcdhelper.EtcdHelper, user *api.NodeUser, ttl uint64) error {
	ctx := prototype.NewContext()
	outItem := new(api.NodeUser)

	key := BuildNodeUserPrefix(user.Spec.NodeName, user.Name)

	err := helper.StorageCodec.Storage.Create(ctx, key, user, outItem, ttl)
	if err != nil {
		return err
	}

	return nil
}

//AddNodeUsers update node user
//like as /api/node/{nodename}/nodeuser/{username}
func DelNodeUsers(nodeName string, helper *etcdhelper.EtcdHelper, name string) error {
	ctx := prototype.NewContext()
	outItem := new(api.NodeUser)
	//it is a strik we use node name for key

	key := BuildNodeUserPrefix(nodeName, name)
	err := helper.StorageCodec.Storage.Delete(ctx, key, outItem, nil)
	if err != nil {
		glog.Errorf("Create node config err %v items %v\r\n", err, outItem)
		return err
	}

	return nil
}

func GetNodeUser(helper *etcdhelper.EtcdHelper, nodeName string, name string) (runtime.Object, error) {
	ctx := prototype.NewContext()
	outItem := new(api.NodeUser)

	key := BuildNodeUserPrefix(nodeName, name)

	err := helper.StorageCodec.Storage.Get(ctx, key, outItem, true)
	if err != nil {
		return nil, err
	}

	return outItem, nil
}

func GetNodeAllUsers(helper *etcdhelper.EtcdHelper, nodeName string) (runtime.Object, error) {

	ctx := prototype.NewContext()
	outItem := new(api.NodeUserList)

	options := &prototype.ListOptions{ResourceVersion: "0"}
	prefix := BuildNodeUserPrefix(nodeName, string(""))
	err := helper.StorageCodec.Storage.List(ctx, prefix, options.ResourceVersion, storage.Everything, outItem)

	return outItem, err
}

func UpdateNodeUsersRefer(helper *etcdhelper.EtcdHelper, spec api.NodeUserSpec) error {

	return updateNodeUsers(helper, spec)
}

//updateNodeUsers update node user not export
func updateNodeUsers(helper *etcdhelper.EtcdHelper, spec api.NodeUserSpec) error {

	nodeName := spec.NodeName

	key := BuildNodeUserPrefix(nodeName, spec.User.Name)

	oldObj, err := GetNodeUser(helper, nodeName, spec.User.Name)
	if err != nil {
		return err
	}

	nodeUser := oldObj.(*api.NodeUser)

	ctx := prototype.NewContext()
	outItem := new(api.NodeUser)

	err = helper.StorageCodec.Storage.GuaranteedUpdate(ctx, key, outItem, false, nil, func(existing runtime.Object, res storage.ResponseMeta) (runtime.Object, *uint64, error) {

		nodeUser.Spec = spec
		return nodeUser, nil, nil
	})

	if err != nil {
		return err
	}

	return nil
}
