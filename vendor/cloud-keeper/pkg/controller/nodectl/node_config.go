package nodectl

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/backend"
	"fmt"
	"strconv"

	"cloud-keeper/pkg/etcdhelper"
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/storage"

	"github.com/golang/glog"
)

const (
	PrefixNodeUser = "/" + "NodeUser"
	PrefixNode     = "/" + "Node"
)

const (
	NodeAnnotationUserCnt = "userCount"
)

func AddNodeUserHelper(helper *etcdhelper.EtcdHelper, nodeName string, userRefer api.UserReferences) error {
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
		},
	}

	glog.V(5).Infof("app node user %+v refer:%v\r\n", *nodeUser, userRefer)

	return AddNodeUser(helper, nodeUser)
}

//AddNodeUsers update node user
//like as /api/node/{nodename}/nodeuser/{username}
func AddNodeUser(helper *etcdhelper.EtcdHelper, user *api.NodeUser) error {
	ctx := prototype.NewContext()
	outItem := new(api.NodeUser)

	key := PrefixNode + "/" + user.Spec.NodeName + PrefixNodeUser + "/" + user.Name

	err := helper.StorageCodec.Storage.Create(ctx, key, user, outItem, 0)
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

	key := PrefixNode + "/" + nodeName + PrefixNodeUser + "/" + name
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

	key := PrefixNode + "/" + nodeName + PrefixNodeUser + "/" + name

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
	prefix := PrefixNode + "/" + nodeName + PrefixNodeUser
	err := helper.StorageCodec.Storage.List(ctx, prefix, options.ResourceVersion, storage.Everything, outItem)

	return outItem, err
}

func UpdateNodeUsersRefer(helper *etcdhelper.EtcdHelper, spec api.NodeUserSpec) error {

	return updateNodeUsers(helper, spec)
}

//updateNodeUsers update node user not export
func updateNodeUsers(helper *etcdhelper.EtcdHelper, spec api.NodeUserSpec) error {

	nodeName := spec.NodeName
	key := PrefixNode + "/" + nodeName + PrefixNodeUser + "/" + spec.User.Name

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

type NodeHelper struct {
	TTL         uint64
	Name        string
	Host        string
	Location    string
	AccsrvID    int64
	AccsrvName  string
	Annotations map[string]string
	Labels      map[string]string
}

func AddNodeToEtcdHelper(helper *etcdhelper.EtcdHelper, nodehelper *NodeHelper) (runtime.Object, error) {

	spec := api.NodeSpec{
		Server: api.NodeServer{
			Name:          nodehelper.Name,
			Host:          nodehelper.Host,
			Method:        "aes-256-cfb",
			Status:        true,
			Location:      nodehelper.Location,
			AccServerID:   nodehelper.AccsrvID,
			AccServerName: nodehelper.AccsrvName,
			EnableOTA:     true,
		},
	}

	node := api.Node{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: prototype.ObjectMeta{
			Name:        spec.Server.Name,
			Labels:      nodehelper.Labels,
			Annotations: nodehelper.Annotations,
		},
		Spec: spec,
	}
	glog.V(5).Infof("add node %v\r\v", node)

	return AddNode(nil, helper, &node, nodehelper.TTL, true, false)
}

func AddNode(db *backend.Backend, helper *etcdhelper.EtcdHelper, node *api.Node, ttl uint64, etcd bool, mysql bool) (runtime.Object, error) {

	if mysql {
		err := db.CreateNode(node.Spec.Server)
		if err != nil {
			return nil, err
		}
	}

	if etcd {
		ctx := prototype.NewContext()
		outItem := new(api.Node)
		err := helper.StorageCodec.Storage.Create(ctx, PrefixNode+"/"+node.Name, node, outItem, ttl)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func DelNode(db *backend.Backend, helper *etcdhelper.EtcdHelper, name string, etcd bool, mysql bool) error {

	if mysql {
		err := db.DeleteNode(name)
		if err != nil {
			return err
		}
	}

	if etcd {
		ctx := prototype.NewContext()
		outItem := new(api.Node)
		err := helper.StorageCodec.Storage.Delete(ctx, PrefixNode+"/"+name, outItem, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateNode(db *backend.Backend, helper *etcdhelper.EtcdHelper, item *api.Node, etcd bool, mysql bool) (runtime.Object, error) {

	if mysql {
		err := db.UpdateNode(item.Spec.Server)
		if err != nil {
			return nil, err
		}
	}

	if etcd {
		ctx := prototype.NewContext()
		outItem := new(api.Node)
		err := helper.StorageCodec.Storage.GuaranteedUpdate(ctx, PrefixNode+"/"+item.Name, outItem, false, nil, func(existing runtime.Object, res storage.ResponseMeta) (runtime.Object, *uint64, error) {
			return item, nil, nil
		})
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func UpdateNodeAnotationsUserCnt(helper *etcdhelper.EtcdHelper, name string, del bool) (runtime.Object, error) {

	oldObj, err := GetNode(helper, name)
	if err != nil {
		return nil, err
	}

	newObj := oldObj.(*api.Node)
	usercnt, ok := newObj.Annotations[NodeAnnotationUserCnt]

	if ok {
		if cnt, err := strconv.ParseInt(usercnt, 10, 32); err == nil {
			if del {
				cnt = cnt - 1
			} else {
				cnt = cnt + 1
			}
			usercnt = strconv.FormatUint(uint64(cnt), 10)
		} else {
			return nil, err
		}
	} else {
		if del {
			usercnt = strconv.FormatUint(uint64(0), 10)
		} else {
			usercnt = strconv.FormatUint(uint64(1), 10)
		}
	}

	ctx := prototype.NewContext()
	outItem := new(api.Node)
	err = helper.StorageCodec.Storage.GuaranteedUpdate(ctx, PrefixNode+"/"+name, outItem, false, nil, func(existing runtime.Object, res storage.ResponseMeta) (runtime.Object, *uint64, error) {

		newObj.Annotations[NodeAnnotationUserCnt] = usercnt
		return newObj, nil, nil
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetNode(helper *etcdhelper.EtcdHelper, name string) (runtime.Object, error) {
	ctx := prototype.NewContext()
	outItem := new(api.Node)

	key := PrefixNode + "/" + name

	err := helper.StorageCodec.Storage.Get(ctx, key, outItem, true)
	if err != nil {
		return nil, err
	}

	return outItem, nil
}

func GetAllNodes(helper *etcdhelper.EtcdHelper) (runtime.Object, error) {
	ctx := prototype.NewContext()
	outItem := new(api.NodeList)

	options := &prototype.ListOptions{ResourceVersion: "0"}
	prefix := PrefixNode + "/"

	filter := NodeFilter{}
	err := helper.StorageCodec.Storage.List(ctx, prefix, options.ResourceVersion, filter, outItem)

	return outItem, err

}

func GetNodeFromDB(db *backend.Backend, name string) (*api.NodeServer, error) {

	nodeserver, err := db.GetNodeByName(name)
	if err != nil {
		return nil, err
	}

	return nodeserver, nil
}
