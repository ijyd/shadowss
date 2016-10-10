package controller

import (
	"fmt"
	"strconv"
	"sync"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/backend"
	"cloud-keeper/pkg/controller/nodectl"
	"cloud-keeper/pkg/etcdhelper"

	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/storage"
	"gofreezer/pkg/watch"

	"github.com/golang/glog"
)

type dynamicNodeInfo struct {
	Name      string
	Host      string
	ISP       string
	UserSpace string
}

type NodeSchedule struct {
	helper             *etcdhelper.EtcdHelper
	be                 *backend.Backend
	availableNode      map[string]dynamicNodeInfo
	availableNodeMutex sync.RWMutex
}

func NewNodeSchedule(helper *etcdhelper.EtcdHelper, be *backend.Backend) *NodeSchedule {
	return &NodeSchedule{
		helper:        helper,
		be:            be,
		availableNode: make(map[string]dynamicNodeInfo),
	}
}

//retrieveNode retrieve available Node  when system startup
func (ns *NodeSchedule) retrieveNode() {
	objList, err := nodectl.GetAllNodes(ns.helper)
	if err != nil {
		glog.Fatalf("get all node failure %v\r\n", err)
	}

	nodeList := objList.(*api.NodeList)

	for _, v := range nodeList.Items {
		ns.checkAvailableNode(v)
	}
}

//getAvailableNodeAPINode retrieve available API Node
func (ns *NodeSchedule) getAvailableNodeAPINode(limit int) []dynamicNodeInfo {
	var info []dynamicNodeInfo
	var loop int
	ns.availableNodeMutex.RLock()
	for _, val := range ns.availableNode {
		glog.V(5).Infof("traverse node %+v loop(%v) limit(%v) userSpace(%v:%v)\r\n", val, loop, limit, val.UserSpace, api.NodeUserSpaceAPI)
		if loop < limit && val.UserSpace == api.NodeUserSpaceAPI {
			loop++
			glog.V(5).Infof("got availableNode %+v \r\n", val)
			info = append(info, val)
		}
	}
	ns.availableNodeMutex.RUnlock()

	return info
}

//getAvailableNodeByISP retrieve available Node  with isp info
func (ns *NodeSchedule) getAvailableNodeByISP(userSpace string, isp string, limit int) []dynamicNodeInfo {
	var info []dynamicNodeInfo
	var loop int
	ns.availableNodeMutex.RLock()
	for _, val := range ns.availableNode {
		glog.V(5).Infof("traverse node %+v loop(%v) limit(%v) userSpace(%v:%v) isp(%v:%v)\r\n", val, loop, limit, val.UserSpace, userSpace, val.ISP, isp)
		if loop < limit && val.UserSpace == userSpace && val.ISP == isp {
			loop++
			glog.V(5).Infof("got availableNode %+v \r\n", val)
			info = append(info, val)
		}
	}
	ns.availableNodeMutex.RUnlock()

	return info
}

//checkAvailableNode check available Node  when have node event
func (ns *NodeSchedule) checkAvailableNode(node api.Node) {

	nodeName := node.Name
	cnt, ok := node.Annotations[nodectl.NodeAnnotationUserCnt]

	if node.Spec.Server.Status != 1 {
		goto delNode
	}

	//check user number on this node
	if ok {
		if cnt, err := strconv.ParseUint(cnt, 10, 32); err == nil && cnt < 80 {
			nodeInfo := dynamicNodeInfo{
				Name:      node.Name,
				Host:      node.Spec.Server.Host,
				ISP:       node.Labels[api.NodeLablesChinaISP],
				UserSpace: node.Labels[api.NodeLablesUserSpace],
			}
			ns.availableNodeMutex.Lock()
			ns.availableNode[nodeInfo.Name] = nodeInfo
			glog.V(5).Infof("get node %+v cnt %v\r\n", ns.availableNode, cnt)
			ns.availableNodeMutex.Unlock()
		} else {
			goto delNode
		}
	} else {
		goto delNode
	}

	return

delNode:
	ns.availableNodeMutex.RLock()
	_, ok = ns.availableNode[nodeName]
	ns.availableNodeMutex.RUnlock()
	if ok {
		ns.availableNodeMutex.Lock()
		delete(ns.availableNode, nodeName)
		ns.availableNodeMutex.Unlock()
	}
}

func (ns *NodeSchedule) NewNodeEvent(node *api.Node) {
	if node.Name == "" {
		glog.Errorf("invalid node name\r\n")
		return
	}

	var mysql bool
	_, err := nodectl.GetNodeFromDB(ns.be, node.Name)
	if err != nil && err.Error() == "not found" {
		mysql = true
	} else if err != nil {
		glog.Errorf("get node from db error %v\r\n", err)
		return
	} else {
		mysql = false
	}

	if mysql {
		_, err := nodectl.AddNode(ns.be, nil, node, 0, false, true)
		if err != nil {
			glog.Errorf("add node to db error %v\r\n", err)
			return
		}
	}
	ns.checkAvailableNode(*node)

	go ns.SyncUserServiceToNodeUser(*node)
}

func (ns *NodeSchedule) UpdateNodeEvent(node *api.Node) {

	glog.V(5).Infof("udpate node %+v\r\n", *node)

	err := ns.UpdateNodeTraffic(node)
	if err != nil {
		glog.Errorf("update node traffic error %v\r\n", err)
	}
	ns.checkAvailableNode(*node)

}

func (ns *NodeSchedule) DelNodeEvent(node *api.Node) {
	err := ns.UpdateNodeStatus(node, false)
	if err != nil {
		glog.Errorf("delete event... disable node %+v error %v\r\n", node, err)
	}
	ns.checkAvailableNode(*node)
}

func (ns *NodeSchedule) configureDynamicNode(info []dynamicNodeInfo, userInfo api.UserInfo) error {

	node2UserRefer := make(map[string]api.UserReferences)
	nodeRefer := make(map[string]api.NodeReferences)
	userName := userInfo.Name
	userRefer := api.UserReferences{
		ID:        userInfo.ID,
		Name:      userName,
		Port:      0,
		Method:    string("aes-256-cfb"),
		Password:  userInfo.Passwd,
		EnableOTA: true,
	}
	for _, v := range info {
		//create a user for node
		node2UserRefer[v.Name] = userRefer

		node := api.NodeReferences{
			User: userRefer,
			Host: v.Host,
		}

		nodeRefer[v.Name] = node
	}

	//direct update user service spec
	err := ns.UpdateNewUserServiceSpec(nodeRefer, userName)
	if err != nil {
		glog.Errorf("update user service %v with %v error %v\r\n", userName, nodeRefer, err)
		return err
	}

	err = ns.BindUserToNode(node2UserRefer)
	if err != nil {
		glog.Errorf("bind user  %v to node %v user err %v", userName, node2UserRefer, err)
		return err
	}

	return nil
}

//AllocDefaultNode search default node for user
func (ns *NodeSchedule) AllocDefaultNode(name string) error {

	dbUserInfo, err := ns.be.GetUserByName(name)
	if err != nil || name != dbUserInfo.Name {
		return fmt.Errorf("invalid user %v error %v", dbUserInfo, err)
	}
	userInfo := *dbUserInfo

	nodeList := ns.findAPINode(name)
	if len(nodeList) == 0 {
		return fmt.Errorf("not have enough node")
	}

	//direct update user service spec
	err = ns.configureDynamicNode(nodeList, userInfo)
	if err != nil {
		glog.Errorf("update user service %v with %v error %v\r\n", userInfo, nodeList, err)
		return err
	}

	return nil
}

func (ns *NodeSchedule) AllocNodeByUserProperties(name string, properties map[string]string) error {

	dbUserInfo, err := ns.be.GetUserByName(name)
	if err != nil || name != dbUserInfo.Name {
		return fmt.Errorf("invalid user %v error %v", dbUserInfo, err)
	}
	userInfo := *dbUserInfo

	nodeList := ns.findNodeByUserProperties(name, properties)
	if len(nodeList) == 0 {
		return fmt.Errorf("not have enough node")
	}

	//direct update user service spec
	err = ns.configureDynamicNode(nodeList, userInfo)
	if err != nil {
		glog.Errorf("update user service %v with %v error %v\r\n", userInfo, nodeList, err)
		return err
	}

	return nil
}

func (ns *NodeSchedule) findAPINode(userName string) []dynamicNodeInfo {

	nodeLimit := 5
	return ns.getAvailableNodeAPINode(nodeLimit)

	// var info []dynamicNodeInfo
	// objList, err := nodectl.GetAllNodes(ns.helper)
	// if err != nil {
	// 	glog.Errorf("get all node failure %v\r\n", err)
	// 	return nil
	// }
	//
	// nodeList := objList.(*api.NodeList)
	// glog.V(5).Infof("check all node %v \r\n", nodeList)
	//
	// for _, v := range nodeList.Items {
	// 	userSpace, ok := v.Labels[api.NodeLablesUserSpace]
	// 	if ok && userSpace == api.NodeUserSpaceAPI {
	// 		node := dynamicNodeInfo{
	// 			name: v.Name,
	// 			host: v.Spec.Server.Host,
	// 		}
	// 		info = append(info, node)
	// 	}
	// }
	//
	// return info
}

func (ns *NodeSchedule) findNodeByUserProperties(userName string, properties map[string]string) []dynamicNodeInfo {

	userISP, _ := properties[api.NodeLablesChinaISP]
	userSpace, _ := properties[api.NodeLablesUserSpace]
	nodeLimit := 4

	glog.V(5).Infof("check user space(%v) isp(%v) with node\r\n", userSpace, userISP)

	return ns.getAvailableNodeByISP(userSpace, userISP, nodeLimit)

	// var info []dynamicNodeInfo
	//
	// objList, err := nodectl.GetAllNodes(ns.helper)
	// if err != nil {
	// 	glog.Errorf("get all node failure %v\r\n", err)
	// 	return nil
	// }
	//
	// nodeList := objList.(*api.NodeList)
	//availableNodeCnt := 0

	// for _, v := range nodeList.Items {
	// 	if availableNodeCnt > 4 {
	// 		//have enough node
	// 		break
	// 	}
	// 	//check this node is default node
	// 	nodeUserSpace, ok := v.Labels[api.NodeLablesUserSpace]
	// 	if ok && userSpace == nodeUserSpace {
	// 		//check user number on this node
	// 		cnt, ok := v.Annotations[nodectl.NodeAnnotationUserCnt]
	// 		if ok {
	// 			if cnt, err := strconv.ParseUint(cnt, 10, 32); err == nil && cnt < 80 {
	// 				cnISP, ok := v.Labels[api.NodeLablesChinaISP]
	// 				if ok && cnISP == userISP {
	// 					node := dynamicNodeInfo{
	// 						name: v.Name,
	// 						host: v.Spec.Server.Host,
	// 					}
	// 					availableNodeCnt++
	// 					info = append(info, node)
	// 				}
	// 			}
	// 		}
	// 	}
	// }
}

func (ns *NodeSchedule) GetActiveNodeByName(name string) (*api.Node, error) {
	obj, err := nodectl.GetNode(ns.helper, name)
	if err != nil {
		return nil, err
	}

	node := obj.(*api.Node)

	return node, nil
}

func (ns *NodeSchedule) GetNodeByName(name string) (*api.NodeServer, error) {
	nodeSrv, err := ns.be.GetNodeByName(name)
	if err != nil {
		return nil, err
	}

	return nodeSrv, nil
}

func (ns *NodeSchedule) UpdateNodeTraffic(node *api.Node) error {
	nodeSrv, err := ns.GetNodeByName(node.Name)
	if err != nil {
		glog.Errorf("update node traffic err %v\r\n", err)
		return err
	}

	upload := nodeSrv.Upload + node.Spec.Server.Upload
	download := nodeSrv.Download + node.Spec.Server.Download
	totalUpload := nodeSrv.TotalUploadTraffic + node.Spec.Server.Upload
	totalDownload := nodeSrv.TotalDownloadTraffic + node.Spec.Server.Download
	err = ns.be.UpdateNodeTraffic(nodeSrv.ID, totalUpload, totalDownload, upload, download)
	if err != nil {
		glog.Errorf("delete %+v error %v\r\n", node, err)
		return err
	}

	return nil
}

func (ns *NodeSchedule) UpdateNodeStatus(node *api.Node, status bool) error {
	nodeSrv, err := ns.GetNodeByName(node.Name)
	if err != nil {
		glog.Errorf("update node traffic err %v\r\n", err)
		return err
	}

	err = ns.be.UpdateNodeStatus(nodeSrv.ID, status)
	if err != nil {
		glog.Errorf("delete %+v error %v\r\n", node, err)
		return err
	}

	return nil
}

func manageNode(helper *etcdhelper.EtcdHelper) {
	watchKey := nodectl.PrefixNode
	ctx := prototype.NewContext()
	resourceVer := string("0")

	glog.V(5).Infof("watch at %v with resource %v", watchKey, resourceVer)
	watcher, err := helper.StorageCodec.Storage.WatchList(ctx, watchKey, resourceVer, storage.Everything)

	if err != nil {
		glog.Fatalf("Unexpected error: %v", err)
	}
	defer watcher.Stop()

	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				glog.Errorf("Unexpected channel close")
				return
			}

			switch event.Type {
			case watch.Added:
				glog.V(5).Infof("Got Add  got: %#v", event.Object)
				gotObject, ok := event.Object.(*api.Node)
				if ok {
					AutoSchedule.NewNodeEvent(gotObject)
				} else {
					gotObject, ok := event.Object.(*api.NodeUser)
					if ok {
						AutoSchedule.NewNodeUserEvent(gotObject)
					}
				}
			case watch.Modified:
				glog.V(5).Infof("Got modify  got: %#v", event.Object)
				gotObject, ok := event.Object.(*api.Node)
				if ok {
					AutoSchedule.UpdateNodeEvent(gotObject)
				} else {
					gotObject, ok := event.Object.(*api.NodeUser)
					if ok {
						AutoSchedule.UpdateNodeUserEvent(gotObject)
					}
				}
			case watch.Deleted:
				glog.V(5).Infof("Got Deleted  got: %#v", event.Object)
				gotObject, ok := event.Object.(*api.Node)
				if ok {
					AutoSchedule.DelNodeEvent(gotObject)
				} else {
					gotObject, ok := event.Object.(*api.NodeUser)
					if ok {
						AutoSchedule.DelNodeUserEvent(gotObject)
					}
				}
			case watch.Error:
				glog.V(5).Infof("Got Error  got: %#v", event.Object)
				return
			default:
				glog.Errorf("UnExpected: %#v, got: %#v", event.Type, event.Object)
			}

		}
	}
}
