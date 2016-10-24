package vps

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/controller"
	"encoding/base64"
	"encoding/json"
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/api/unversioned"
	"golib/pkg/util/crypto"

	apierr "cloud-keeper/pkg/api/errors"
	. "cloud-keeper/pkg/api/vps/common"

	restful "github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

var rsaPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

var rsaPublicKey = `-----BEGIN CERTIFICATE-----
MIIFrTCCA5WgAwIBAgIJAKfuafURcL6oMA0GCSqGSIb3DQEBCwUAMG0xCzAJBgNV
BAYTAmNuMQswCQYDVQQIDAJjbjELMAkGA1UEBwwCY24xCzAJBgNVBAoMAmNuMQsw
CQYDVQQLDAJjbjEMMAoGA1UEAwwDYm94MRwwGgYJKoZIhvcNAQkBFg1ib3hAZ21h
aWwuY29tMB4XDTE2MTAxMTAyMDMzNVoXDTQ0MDIyNzAyMDMzNVowbTELMAkGA1UE
BhMCY24xCzAJBgNVBAgMAmNuMQswCQYDVQQHDAJjbjELMAkGA1UECgwCY24xCzAJ
BgNVBAsMAmNuMQwwCgYDVQQDDANib3gxHDAaBgkqhkiG9w0BCQEWDWJveEBnbWFp
bC5jb20wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDhUhU22Bn/Ktbx
ZPN0vKY57cUS5oIoIdMoh66YNm+0WsL0HBl6+bx13GHAAUHI7i1XqUl3G+GuhNde
6umx7DELwnMIDpHLy8m5tlqma99oD2RYZWg6f1ya2XzNpG1CElKDQYqRkzEDyOlF
AXWqsKkREN84QSLC9pV0YZb3TA9lu/YT30Qsfbwd67aEY++30njg24jK1YiPCG82
2t3dF985n+4dC3jFH+ZJoSFYiExrczgSr+llKrIBpf67gq6BPo0JgSkyRbEp2eNY
dwRODmYCfP6NTvKsGQV81MQ6b7dcWCjIEUjp8clKNeoQpIAgcDuOWhsgz6ojAcW+
MeX3SnH9880R9JVL/Ic2tF5whf4DAZm38m0PvDvdsMxH44foSipWY5CQcHI3rXuN
hbm/LH/JdLxXOaL+k1Rw8Dz53u2h6B6OR4MR9ZqVxYvwl43VUvw5zF5HZo94gNE3
cfnwgcd4tChVZlKDKadMuMA5K6FCVQXlG/Lb1XGv5FABNVVNKER09/WxwJvjljXt
JAfj9P+qtKCHS2+V+f6DMOSJAJLvcxn0LCM2nDUlpnBX8cu93QMcTMfDjNSDLoEn
HCszZDouOUMbbbVPJINUrC9CYI7PccWP5ZeaOZ+q97Camxxn1Sx5g0WpT47m/y87
mDytkXpK7Pc02J4IhM26jG5y8zqfRwIDAQABo1AwTjAdBgNVHQ4EFgQUO+U9RFt0
qKMx5c4k0bh3ED5ktbAwHwYDVR0jBBgwFoAUO+U9RFt0qKMx5c4k0bh3ED5ktbAw
DAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAgEAw5Dc9sbm8Hqehlurm5Md
/q6RMij9Vjxv65L8JJ0O/gQuOgIicT7vmVLzwTDT9LeOuXKUSOu/R9vkPoSgXODr
ay486WZe2rtkE2Eew49mCPXOB1qoSiLkaySXoWjQzus2C++xVsolr70TVZhPw5V/
dhxOZhWFGsa4oi1SkVwliO2f6lbEx49VH/avlNVv5weYjPSTO+c11iOqEZszlzUT
QbVTpXBBJRV1gIhdx/ysfP4//s7MJRpuNhBvwJtAzPpe9rkLHNQme0A0gsmqWs70
8W/61UnOIFd0OR9KLleObyet73rMGXY7zw7bVNyJWGdGIWSg/ONWkajweqa4wGjV
/fXlLcINpwJkQdLP8E96KKwSQbdCk3d1qoiQwPhbfP6tfXm6lU38aYL8x8Ho93MQ
8GgjY37x+r8pxul7OjTgL4Q6Re+CRneSv99sUbgxlqYJryG9+zGOlMtfhWSqqT/E
tKDbrVEwi+KUuJ0GQBUWRnpIfxy7LcTFHezlwZTIKp04ifl+ysEt/g+q4AsFYsp8
kic9aLeZpnirwRSPYE30Z83wbdKdAa9Qq4gLlT6e+amW3pm4LbhbPlOcO99ax15y
V/C1TTnEnZZAYv4UqcX9SEBK6N8+7G8pSTJuuB12dVl0BTLq3fGmJOtGs69fd1oW
PbmnvB+0obn9amKJm63MmiU=
-----END CERTIFICATE-----`

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
	var activeAPINodeList []api.ActiveAPINodeSpec
	nodeInfo := controller.GetAvailableNodeAPINode(limit)
	glog.V(5).Infof("Get api node %+v\r\n", nodeInfo)
	for _, v := range nodeInfo {
		apiNode := api.ActiveAPINodeSpec{
			Host:     v.Host,
			Port:     48888,
			Method:   "aes-256-cfb",
			Password: "48c8591290877f737202ad20c06780e9",
		}

		item := api.ActiveAPINode{
			TypeMeta: unversioned.TypeMeta{
				Kind:       "APIServer",
				APIVersion: "v1",
			},
			ObjectMeta: prototype.ObjectMeta{
				Name: v.Name,
			},
			Spec: apiNode,
		}
		nodes = append(nodes, item)
		activeAPINodeList = append(activeAPINodeList, apiNode)
	}

	apiInfo, err := json.Marshal(activeAPINodeList)
	if err != nil {
		newErr := apierr.NewInternalError(err.Error())
		output = EncodeError(newErr)
		statusCode = 500
	}

	rsaUtil := crypto.NewRSAHelper(nil, []byte(rsaPublicKey))
	encData, err := rsaUtil.RsaEncrypt(apiInfo)
	if err != nil {
		newErr := apierr.NewInternalError(err.Error())
		output = EncodeError(newErr)
		statusCode = 500
	}
	encDataBase64 := base64.StdEncoding.EncodeToString(encData)

	nodeList := api.ActiveAPINodeList{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "ActiveAPINodeList",
			APIVersion: "v1",
		},
		ListMeta: unversioned.ListMeta{
			SelfLink: "/api/v1/activeapinode",
		},
		Items:       nodes,
		EncryptData: encDataBase64,
	}

	output, err = json.Marshal(nodeList)
	if err != nil {
		newErr := apierr.NewInternalError(err.Error())
		output = EncodeError(newErr)
		statusCode = 500
	}

}
