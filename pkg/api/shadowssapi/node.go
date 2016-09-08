package shadowssapi

import (
	"encoding/json"

	"shadowsocks-go/pkg/api"
	apierr "shadowsocks-go/pkg/api/errors"
	"shadowsocks-go/pkg/backend/db"

	restful "github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

func getNodeInfo(user *db.User) ([]byte, int) {
	statusCode := 500
	var output []byte
	var err error

	var nodeList api.NodeList
	nodes, err := Storage.GetNodesByUID(user.ID)
	if err != nil {
		if err.Error() == "not found" {
			nodeList = api.NodeList{
				TypeMeta: api.TypeMeta{
					Kind:       "NodeList",
					APIVersion: "v1",
				},
				ListMeta: api.ListMeta{
					SelfLink: "/api/v1/nodes",
				},
			}
		} else {
			glog.Errorf("Get nodes by id %v failure %v \r\n", user.ID, err)
			newErr := apierr.NewInternalError("marshal nodes resource failure")
			internalErr, _ := newErr.(*apierr.StatusError)

			output = internalErr.ErrStatus.Encode()
			return output, statusCode
		}
	} else {
		var items []api.Node
		for _, v := range nodes {
			item := api.Node{
				TypeMeta: api.TypeMeta{
					Kind:       "Node",
					APIVersion: "v1",
				},
				Spec: api.NodeSpec{
					Server: api.NodeServer{
						Host:   v.Host,
						Status: (v.Status) == string("Enable"),
					},
					Account: api.NodeAccout{
						ID:     user.ID,
						Port:   user.Port,
						Method: user.Method,
					},
				},
			}
			items = append(items, item)
		}

		nodeList = api.NodeList{
			TypeMeta: api.TypeMeta{
				Kind:       "NodeList",
				APIVersion: "v1",
			},
			ListMeta: api.ListMeta{
				SelfLink: "/api/v1/nodes",
			},
			Items: items,
		}
	}

	output, err = json.Marshal(nodeList)
	if err != nil {
		glog.Errorln("Marshal router err", err)
		newErr := apierr.NewInternalError("marshal router resource failure")
		internalErr, _ := newErr.(*apierr.StatusError)

		output = internalErr.ErrStatus.Encode()

	} else {
		statusCode = 200
	}
	return output, statusCode
}

//GetRouters ... get router list
func GetNodes(request *restful.Request, response *restful.Response) {

	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	encoded := request.Request.Header.Get("Authorization")
	statusCode := 200
	var output []byte

	user, err := CheckToken(encoded)
	if err != nil || user == nil {
		glog.Errorln("Unauth request %v", err)
		newErr := apierr.NewUnauthorized("invalid token")
		internalErr, ok := newErr.(*apierr.StatusError)
		if ok {
			output = internalErr.ErrStatus.Encode()
		} else {
			glog.Errorln("status type error")
		}
		statusCode = 401
	} else {
		output, statusCode = getNodeInfo(user)
	}

	w.WriteHeader(statusCode)
	w.Write(output)
}
