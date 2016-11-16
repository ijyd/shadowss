package rest

import (
	apierr "cloud-keeper/cmd/prometheushook/app/errors"
	"encoding/json"
	"golib/pkg/util/exec"
	"io/ioutil"
	"os"

	restful "github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

type Alert struct {
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	StartsAt    string            `json:"startsAt,omitempty"`
	EndsAt      string            `json:"endsAt,omitempty"`
}

type AlertMgrHook struct {
	Vsersion string  `json:"version,omitempty"`
	Status   string  `json:"status,omitempty"`
	Alerts   []Alert `json:"alerts,omitempty"`
}

func execAlertScript() error {
	execCom := exec.New()
	cmd := execCom.Command("./alert.sh")
	out, err := cmd.CombinedOutput()
	if err == exec.ErrExecutableNotFound {
		glog.Errorf("Expected error ErrExecutableNotFound but got %v", err)
		return err
	}
	glog.V(5).Infof("alert.sh result %s", string(out))
	return nil
}

func PostAlert(request *restful.Request, response *restful.Response) {
	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")

	statusCode := 200
	output := apierr.NewSuccess().Encode()

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	alert := new(AlertMgrHook)
	err := request.ReadEntity(alert)

	if err != nil {
		newErr := apierr.NewBadRequestError("request body invalid")
		output = apierr.EncodeError(newErr)
		statusCode = 400
		return
	}

	data, err := json.Marshal(alert)
	if err != nil {
		glog.Errorf("marshal result error:%v\r\n", err)
	}
	glog.V(5).Infof("Got alert data:%+v\r\n", *alert)
	ioutil.WriteFile("./alert.json", data, os.FileMode(0664))

	err = execAlertScript()

	output = apierr.NewSuccess().Encode()
	statusCode = 200

	return
}
