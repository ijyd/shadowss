package nodectl

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/backend"
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
	NodeAnnotationUserCnt    = "userCount"
	NodeAnnotationRefreshCnt = "Refresh"
	NodeAnnotationVersion    = "version"
)

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
			Status:        1,
			Location:      nodehelper.Location,
			AccServerID:   nodehelper.AccsrvID,
			AccServerName: nodehelper.AccsrvName,
			EnableOTA:     1,
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

//UpdateNodeLease update node leass
func UpdateNodeAndLease(helper *etcdhelper.EtcdHelper, item *api.Node, ttl uint64) (runtime.Object, error) {

	ctx := prototype.NewContext()
	outItem := new(api.Node)
	err := helper.StorageCodec.Storage.GuaranteedUpdate(ctx, PrefixNode+"/"+item.Name, outItem, false, nil, func(existing runtime.Object, res storage.ResponseMeta) (runtime.Object, *uint64, error) {
		return item, &ttl, nil
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
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

	var count uint64
	if ok {
		if cnt, err := strconv.ParseInt(usercnt, 10, 32); err == nil {
			if del {
				if cnt == 0 {
					cnt = 0
				} else {
					cnt = cnt - 1
				}
			} else {
				cnt = cnt + 1
			}
			count = uint64(cnt)
		} else {
			return nil, err
		}
	} else {
		if del {
			count = 0
		} else {
			count = 1
		}
	}

	usercnt = strconv.FormatUint(uint64(count), 10)

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
