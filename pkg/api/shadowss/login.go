package shadowssapi

import (
	"encoding/json"

	"shadowsocks-go/pkg/api"

	restful "github.com/emicklei/go-restful"

	"github.com/golang/glog"
)

const (
	maxLoginCacheSize = 32 * 4
)

var cache = cache.NewCache(maxLoginCacheSize)

func PostRouter(request *restful.Request, response *restful.Response) {
	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	statusCode := 200

	output := api.NewSuccess().Encode()

	login := new(api.Login)
	err := request.ReadEntity(login)
	if err != nil {
		glog.Errorf("invalid request body:%v", err)
		newErr := errors.NewBadRequestError("request body invalid")
		internalErr, ok := newErr.(*errors.StatusError)
		if ok {
			output = internalErr.ErrStatus.Encode()
		} else {
			glog.Errorln("status type error")
		}
		statusCode = 400
	} else {
		glog.Infof("Got Post Routers:%+v\n", router)
		routerMap[router.Name] = *router
		output, err = json.Marshal(router)
		statusCode = 200

	}

	w.WriteHeader(statusCode)
	w.Write(output)
}
