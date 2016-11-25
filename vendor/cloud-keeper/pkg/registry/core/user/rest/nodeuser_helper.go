package rest

import (
	"cloud-keeper/pkg/api"
	"fmt"
	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/errors"
	"gofreezer/pkg/labels"

	"github.com/golang/glog"
)

const (
	MaxNodeForUser = 4
)

func (r *UserREST) updateNodeUser(ctx freezerapi.Context, userRefer *api.UserReferences, nodeName string) error {
	nodeUser := r.NewNodeUser(userRefer, nodeName)
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

func (r *UserREST) DumpNodeUserToNode(ctx freezerapi.Context, nodename string) error {
	userlist, err := r.FilterUserWithNodeName(ctx, nodename, nil)
	if err != nil {
		return err
	}
	for _, user := range userlist.Items {
		nodeRefer, ok := user.Spec.UserService.Nodes[nodename]
		if ok {
			err = r.updateNodeUser(ctx, &nodeRefer.User, nodename)
			if err != nil {
				glog.Warningf("create node user:%+v failure:%v\r\n", err)
			}
		}
	}

	return err
}

func (r *UserREST) AddNodeUser(ctx freezerapi.Context, updateUser *api.User, userservice *api.UserService) error {

	nodeName := userservice.Spec.NodeName

	userRefer := userservice.Spec.UserRefer
	err := r.updateNodeUser(ctx, &userRefer, nodeName)
	if err != nil {
		return err
	}
	nodeRefer := api.NodeReferences{
		User: userRefer,
		Host: userservice.Spec.Host,
	}
	updateUser.Spec.UserService.Nodes[nodeName] = nodeRefer

	return nil
}

func (r *UserREST) InitNodeUser(ctx freezerapi.Context, user *api.User) error {
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
		if err != nil {
			return err
		}
	} else {
		nodelist, err = r.node.GetAPINodes(ctx, nil)
		if err != nil {
			return errors.NewInternalError(fmt.Errorf("create user(%v) with annotation(%v) find api api node error %+v\r\n",
				user.Name, user.Annotations, err))
		}
	}

	user.Spec.UserService.Nodes = make(map[string]api.NodeReferences)
	var i int
	for _, v := range nodelist.Items {
		err = r.updateNodeUser(ctx, &userRefer, v.Name)
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
