package rest

import (
	"cloud-keeper/pkg/api"
	"fmt"
	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/errors"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/labels"
	"strconv"

	"github.com/golang/glog"
)

const (
	MaxNodeForUser = 4
)

func (r *UserREST) checkUserTraffic(user *api.User) bool {

	upload := user.Spec.DetailInfo.UploadTraffic
	download := user.Spec.DetailInfo.DownloadTraffic

	traffic := upload + download
	limit := user.Spec.DetailInfo.TrafficLimit

	if traffic > limit {
		return false
	}

	return true
}

//checkuser check user traffic and status
func (r *UserREST) CheckUser(user *api.User) bool {
	if !user.Spec.DetailInfo.Status {
		return false
	}

	if !r.checkUserTraffic(user) {
		return false
	}

	return true
}

func (r *UserREST) updateNodeUser(ctx freezerapi.Context, userRefer *api.UserReferences, nodeName string, delete bool) error {
	nodeUser := r.NewNodeUser(userRefer, nodeName)
	if delete {
		nodeUser.Spec.Phase = api.NodeUserPhaseDelete
	}
	_, err := r.nodeuser.UpdateNodeUser(ctx, nodeUser)
	if err != nil {
		return fmt.Errorf("update the user(%v) in node(%v) error:%v", userRefer.Name, nodeName, err)
	}

	return nil
}

func (r *UserREST) NewNodeUser(user *api.UserReferences, nodeName string) *api.NodeUser {
	nodeUser := &api.NodeUser{
		Spec: api.NodeUserSpec{
			User:     *user,
			NodeName: nodeName,
			Phase:    api.NodeUserPhase(api.NodeUserPhaseAdd),
		},
	}
	nodeUser.Name = user.Name
	return nodeUser
}

func (r *UserREST) DumpNodeUserToNode(ctx freezerapi.Context, nodename string) ([]*api.NodeUser, error) {
	glog.V(5).Infof("dump user into node(%v)", nodename)

	userlist, err := r.FilterUserWithNodeName(ctx, nodename, nil)
	if err != nil {
		glog.Errorf("filter userlist with node(%v) error:%v\r\n", nodename, err)
		return nil, err
	}
	glog.V(5).Infof("dump %v user into node(%v)", len(userlist.Items), nodename)
	var nodeUserList []*api.NodeUser
	for _, user := range userlist.Items {
		nodeRefer, ok := user.Spec.UserService.Nodes[nodename]
		if ok {
			nodeUser := r.NewNodeUser(&nodeRefer.User, nodename)
			nodeUserList = append(nodeUserList, nodeUser)
			//need to call nodeuser dump
			_, err = r.nodeuser.UpdateNodeUser(ctx, nodeUser)
			if err != nil {
				glog.Warningf("dump nodeuser(%) to node(%v) error:%v", nodeUser.Name, nodename, err)
			}
		}
	}

	return nodeUserList, nil
}

func (r *UserREST) DelNodeUser(ctx freezerapi.Context, user *api.User, nodeName string, update bool, syncToNode bool) error {
	delUserService := &user.Spec.UserService
	userName := user.Name
	userRefer, ok := delUserService.Nodes[nodeName]
	if ok {
		delete(delUserService.Nodes, nodeName)
	} else {
		return errors.NewBadRequest(fmt.Sprintf("not found node(%v) in user(%v)", nodeName, userName))
	}

	var err error
	if syncToNode {
		err = r.updateNodeUser(ctx, &userRefer.User, nodeName, true)
		if err != nil {
			glog.Errorf("del user(%s) in node(%s)  error:%v", userName, nodeName, err)
			return err
		}
	}

	if update {
		_, _, err = r.Update(ctx, user.Name, rest.DefaultUpdatedObjectInfo(user, api.Scheme))
	}

	return err
}

func (r *UserREST) DelAllNodeUser(ctx freezerapi.Context, user *api.User, update bool, syncToNode bool) error {
	delUserService := &user.Spec.UserService
	userName := user.Name

	for nodeName, _ := range delUserService.Nodes {
		err := r.DelNodeUser(ctx, user, nodeName, update, syncToNode)
		if err != nil {
			glog.Warningf("delete user(%v) from node(%v) failure:%vr\r\n", userName, nodeName)
		}
	}

	return nil
}

func (r *UserREST) AddNodeUser(ctx freezerapi.Context, updateUser *api.User,
	userservice *api.UserService, update bool, syncToNode bool) error {
	if !r.CheckUser(updateUser) {
		return fmt.Errorf("user disable or traffic exceed")
	}

	nodeName := userservice.Spec.NodeName

	userRefer := userservice.Spec.UserRefer

	if syncToNode {
		err := r.updateNodeUser(ctx, &userRefer, nodeName, false)
		if err != nil {
			return err
		}
	}

	nodeRefer := api.NodeReferences{
		User: userRefer,
		Host: userservice.Spec.Host,
	}
	updateUser.Spec.UserService.Nodes[nodeName] = nodeRefer

	var err error
	if update {
		_, _, err = r.Update(ctx, updateUser.Name, rest.DefaultUpdatedObjectInfo(updateUser, api.Scheme))
	}

	return err
}

func (r *UserREST) InitNodeUser(ctx freezerapi.Context, user *api.User) error {
	if !r.CheckUser(user) {
		return fmt.Errorf("user disable or traffic exceed")
	}

	//for user default info
	userRefer := api.UserReferences{
		ID:        user.Spec.DetailInfo.ID,
		Name:      user.Name,
		Port:      0,
		Method:    string("aes-256-cfb"),
		Password:  user.Spec.DetailInfo.Passwd,
		EnableOTA: true,
	}

	userspace, hasUserSpace := user.Annotations[api.NodeLablesUserSpace]
	isp, hasISP := user.Annotations[api.NodeLablesChinaISP]

	options := &freezerapi.ListOptions{}
	var nodelist *api.NodeList
	var err error
	if hasUserSpace && hasISP {
		ls := labels.Set(map[string]string{
			api.NodeLablesUserSpace: userspace,
			api.NodeLablesChinaISP:  isp,
		})
		options.LabelSelector = labels.SelectorFromSet(ls)
		nodelist, err = r.node.ListNodes(ctx, options)
		glog.V(5).Infof("use options(%+v) Got proxy Nodes(%+v)\r\n", *options, *nodelist)
		if err != nil {
			return err
		}
	} else {
		nodelist, err = r.node.GetAPINodes(ctx, nil)
		if err != nil {
			return errors.NewInternalError(fmt.Errorf("create user(%v) with annotation(%v) find api api node error %+v\r\n",
				user.Name, user.Annotations, err))
		}
		glog.V(5).Infof("Got API Nodes(%+v)\r\n", *nodelist)
	}

	user.Spec.UserService.Nodes = make(map[string]api.NodeReferences)
	var i int
	for _, v := range nodelist.Items {
		userInNodeCnt, ok := v.Labels[api.NodeAnnotationUserCnt]
		if ok {
			if cnt, err := strconv.ParseUint(userInNodeCnt, 10, 32); err == nil && cnt > 80 {
				continue
			}
		}

		err = r.updateNodeUser(ctx, &userRefer, v.Name, false)
		if err != nil {
			glog.Warningf("update node user :%v\r\n", err)
		}
		user.Spec.UserService.Nodes[v.Name] = api.NodeReferences{
			User: userRefer,
			Host: v.Spec.Server.Host,
		}

		if i++; i >= MaxNodeForUser {
			break
		}
	}
	user.Spec.UserService.NodeCnt = uint(len(nodelist.Items))
	return nil
}
