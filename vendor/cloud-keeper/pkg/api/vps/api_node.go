package vps

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/controller"
	"encoding/json"
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/api/unversioned"

	apierr "cloud-keeper/pkg/api/errors"
	. "cloud-keeper/pkg/api/vps/common"

	restful "github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

//GetAPINodes ...
func GetAPINodes(request *restful.Request, response *restful.Response) {

	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	encoded := request.Request.Header.Get("Authorization")
	statusCode := 200
	var output []byte

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	user, err := CheckToken(encoded)
	if err != nil || user == nil {
		glog.Errorln("Unauth request ", err)
		newErr := apierr.NewUnauthorized("invalid token")
		output = EncodeError(newErr)
		statusCode = 401
		return
	}

	limit := 4
	var nodes []api.ActiveAPINode
	nodeInfo := controller.GetAvailableNodeAPINode(limit)
	glog.V(5).Infof("Get api node %+v\r\n", nodeInfo)
	for _, v := range nodeInfo {
		item := api.ActiveAPINode{
			TypeMeta: unversioned.TypeMeta{
				Kind:       "APIServer",
				APIVersion: "v1",
			},
			ObjectMeta: prototype.ObjectMeta{
				Name: v.Name,
			},
			Spec: api.ActiveAPINodeSpec{
				Host:     v.Host,
				Port:     48888,
				Method:   "aes-256-cfb",
				Password: "48c8591290877f737202ad20c06780e9",
			},
		}
		nodes = append(nodes, item)
	}

	nodeList := api.ActiveAPINodeList{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "ActiveAPINodeList",
			APIVersion: "v1",
		},
		ListMeta: unversioned.ListMeta{
			SelfLink: "/api/v1/activeapinode",
		},
		Items: nodes,
	}

	output, err = json.Marshal(nodeList)
	if err != nil {
		newErr := apierr.NewInternalError(err.Error())
		output = EncodeError(newErr)
		statusCode = 500
	}

}
