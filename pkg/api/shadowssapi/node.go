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

	nodes, err := Storage.GetNodesByUID(user.ID)
	if err != nil {
		glog.Errorf("Get nodes by id %v failure %v \r\n", user.ID, err)
		newErr := apierr.NewInternalError("marshal router resource failure")
		internalErr, _ := newErr.(*apierr.StatusError)

		output = internalErr.ErrStatus.Encode()
		return output, statusCode
	}

	var nodeSrv []api.NodeServer
	for _, v := range nodes {
		server := api.NodeServer{
			Host:   v.Host,
			Status: (v.Status) == string("Enable"),
		}
		nodeSrv = append(nodeSrv, server)
	}

	ss := api.Node{
		TypeMeta: api.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		Spec: api.NodeSpec{
			Server: nodeSrv,
			Account: api.NodeAccout{
				ID:     user.ID,
				Port:   user.Port,
				Method: user.Method,
			},
		},
	}

	output, err = json.Marshal(ss)
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
