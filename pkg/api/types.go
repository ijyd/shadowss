package api

import (
	"time"

	"github.com/seanchann/goutil/uuid"
)

//TypeMeta import from k8s
type TypeMeta struct {
	Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`

	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`
}

//ObjectMeta import from k8s
type ObjectMeta struct {
	Name         string `json:"name,omitempty"`
	GenerateName string `json:"generateName,omitempty"`

	Namespace string `json:"namespace,omitempty"`

	SelfLink string `json:"selfLink,omitempty"`

	UID uuid.UID `json:"uid,omitempty"`

	ResourceVersion string `json:"resourceVersion,omitempty"`

	Generation int64 `json:"generation,omitempty"`

	Labels map[string]string `json:"labels,omitempty"`

	Annotations map[string]string `json:"annotations,omitempty"`
}

//ListMeta import from k8s
type ListMeta struct {
	SelfLink string `json:"selfLink,omitempty" protobuf:"bytes,1,opt,name=selfLink"`

	ResourceVersion string `json:"resourceVersion,omitempty" protobuf:"bytes,2,opt,name=resourceVersion"`
}

const (
	//NodeUserPhaseAdd dynamic add user by api server
	NodeUserPhaseAdd = "add"
	//NodeUserPhaseDelete dynamic delete user by api server
	NodeUserPhaseDelete = "del"
	//NodeUserPhaseUpdate dynamic update user from node
	NodeUserPhaseUpdate = "update"
)

//NodeUserPhase current user action
type NodeUserPhase string

//UserReferences contains all data to start a shadowsocket service for user
type UserReferences struct {
	ID              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Port            int64  `json:"port,omitempty"`
	Method          string `json:"method,omitempty"`
	Password        string `json:"password,omitempty"`
	EnableOTA       bool   `json:"enableOTA, omitempty"`
	UploadTraffic   int64  `json:"uploadTraffic,omitempty"`   //upload traffic for per user
	DownloadTraffic int64  `json:"downloadTraffic,omitempty"` //download traffic for per user
}

//NodeUserSpec configure a user into node
type NodeUserSpec struct {
	User     UserReferences `json:"user,omitempty"`
	NodeName string         `json:"nodeName,omitempty"`
	Phase    NodeUserPhase  `json:"phase,omitempty"`
}

//NodeUser manager user in node
type NodeUser struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata,omitempty"`

	Spec NodeUserSpec `json:"spec,omitempty"`
}

const (
	//CNISPCMCC Network service provider cmcc
	CNISPCMCC = "cnISPCMCC"
	//CNISPUNICOM Network service provider unicom
	CNISPUNICOM = "cnISPUnicom"
	//CNISPASPCTCC Network service provider ctcc
	CNISPASPCTCC = "cnISPCTCC"
	//CNISPOther Network service provider other
	CNISPOther = "cnISPOther"
)

const (
	//NodeUserSpaceDefault normal node for user
	NodeUserSpaceDefault = "default"
	//NodeUserSpaceAPI api node for user
	NodeUserSpaceAPI = "api"
	//NodeUserSpaceDev development node for user
	NodeUserSpaceDev = "develop"
)

const (
	//NodeLablesChinaISP china isp type
	NodeLablesChinaISP = "cnISP"
	//NodeLablesUserSpace user space
	NodeLablesUserSpace = "userSpace"
	//NodeLablesVPSLocation vps location
	NodeLablesVPSLocation = "vpsLocation"
	//NodeLablesVPSID vps id
	NodeLablesVPSID = "vpsID"
	//NodeLablesVPSOP vps operator
	NodeLablesVPSOP = "vpsOperator"
	//NodeLablesVPSName vps name
	NodeLablesVPSName = "vpsName"
	//NodeLablesVPSIP vps ip addr
	NodeLablesVPSIP = "vpsIP"
)

const (
	//NodeAnnotationUserCnt  user count in shadowss node
	NodeAnnotationUserCnt = "userCount"
	//NodeAnnotationRefreshCnt node refresh count
	NodeAnnotationRefreshCnt = "Refresh"
	//NodeAnnotationVersion node version
	NodeAnnotationVersion = "version"
	//NodeAnnotationRefreshTime refresh time
	NodeAnnotationRefreshTime = "refreshTime"
)

//NodeServer Describe the current node information
type NodeServer struct {
	ID                   int64  `json:"-"`
	Name                 string `json:"name,omitempty"`
	EnableOTA            bool   `json:"enableOTA"`
	Host                 string `json:"host,omitempty"`
	Method               string `json:"method"`
	Status               int64  `json:"status,omitempty"`
	Location             string `json:"location,omitempty"`
	AccServerID          int64  `json:"accServerID,omitempty"`
	AccServerName        string `json:"accServerName,omitempty"`
	Description          string `json:"description,omitempty"`
	TrafficLimit         int64  `json:"trafficLimit,omitempty"`
	Upload               int64  `json:"upload,omitempty"`
	Download             int64  `json:"download,omitempty"`
	TrafficRate          int64  `json:"trafficRate,omitempty"`
	TotalUploadTraffic   int64  `json:"totalUploadTraffic,omitempty"`
	TotalDownloadTraffic int64  `json:"totalDownloadTraffic,omitempty"`
	CustomMethod         int    `json:"customMethod,omitempty"`
	RawObj               []byte `json:"-"`
}

//NodeSpec shadowss node server information
type NodeSpec struct {
	Server NodeServer              `json:"server,omitempty" freezer:"table:vps_new_node"`
	Users  map[string]NodeUserSpec `json:"users,omitempty"`
}

//Node manage shadowss node
type Node struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata,omitempty"`

	Spec NodeSpec `json:"spec,omitempty"`
}

//APIServerInfor api server information. if status is false, this server disable
type APIServerInfor struct {
	ID         int64     `json:"id, omitempty" column:"id"`
	Name       string    `json:"name, omitempty" column:"name"`
	Host       string    `json:"host, omitempty" column:"host"`
	Port       int64     `json:"port, omitempty" column:"port"`
	Status     bool      `json:"status, omitempty" column:"status"`
	CreateTime time.Time `json:"creationTime,omitempty" column:"created_time" gorm:"column:created_time"`
}

//APIServerSpec api server spec
type APIServerSpec struct {
	Server   APIServerInfor `json:"server, omitempty"`
	HostList []string       `json:"hosts, omitempty"`
}

//APIServer api server information
type APIServer struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata,omitempty"`

	Spec APIServerSpec `json:"spec,omitempty"`
}

//APIServerList list apiserver
type APIServerList struct {
	TypeMeta `json:",inline"`
	ListMeta `json:"metadata,omitempty"`

	Items []APIServer `json:"items"`
}

const (
	//UserFakeAnnotationLastActiveTime user last active time
	UserFakeAnnotationLastActiveTime = "lastActiveTime"
)
