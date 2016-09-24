package v1

import (
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/api/unversioned"
	"time"
)

type UserReferences struct {
	ID              int64  `json:"id,omitempty"`
	Port            int64  `json:"port,omitempty"`
	Method          string `json:"method,omitempty"`
	Password        string `json:"password,omitempty"`
	EnableOTA       bool   `json:"enableOTA, omitempty"`
	UploadTraffic   int64  `json:"uploadTraffic,omitempty"`   //upload traffic for per user
	DownloadTraffic int64  `json:"downloadTraffic,omitempty"` //download traffic for per user
}

type NodeUserSpec struct {
	User     UserReferences `json:"user,omitempty"`
	NodeName string         `json:"nodeName,omitempty"`
}

//put your user into your node with node name
//like as /api/node/{nodename}/nodeuser/{resourcename}
type NodeUser struct {
	unversioned.TypeMeta `json:",inline"`
	ObjectMeta           `json:"metadata,omitempty"`

	Spec NodeUserSpec `json:"spec,omitempty"`
}

type NodeUserList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []NodeUser `json:"items,omitempty"`
}

type NodeServer struct {
	ID            int64  `json:"id" column:"id" gorm:"column:id"`
	Name          string `json:"name,omitempty" column:"name" gorm:"column:name"`
	EnableOTA     bool   `json:"enableOTA" column:"enableota" gorm:"column:enableota"`
	Host          string `json:"host,omitempty" column:"server" gorm:"column:server"`
	Method        string `json:"method" column:"method" gorm:"column:method"`
	Status        bool   `json:"status,omitempty" column:"status" gorm:"column:status"`
	Location      string `json:"location,omitempty" column:"location" gorm:"column:location"`
	AccServerID   int64  `json:"accServerID,omitempty" column:"vps_server_id" gorm:"column:vps_server_id"`
	AccServerName string `json:"accServerName,omitempty" column:"vps_server_name" gorm:"column:vps_server_name"`
	Descryption   string `json:"descryption,omitempty" column:"descryption" gorm:"column:descryption"`
	TrafficLimit  int64  `json:"trafficLimit,omitempty" column:"traffic_limit" gorm:"column:traffic_limit"`
	Upload        int64  `json:"upload,omitempty" column:"upload" gorm:"column:upload"`
	Download      int64  `json:"download,omitempty" column:"download" gorm:"column:download"`
	TrafficRate   int64  `json:"trafficRate,omitempty" column:"traffic_rate" gorm:"column:traffic_rate"`
}

type NodeSpec struct {
	Server NodeServer `json:"server,omitempty"`
}

type Node struct {
	unversioned.TypeMeta `json:",inline"`
	prototype.ObjectMeta `json:"metadata,omitempty"`

	Spec NodeSpec `json:"spec,omitempty"`
}

type NodeList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []Node `json:"items"`
}

type APIServerInfor struct {
	ID         int64     `json:"id, omitempty" column:"id"`
	Name       string    `json:"name, omitempty" column:"name"`
	Host       string    `json:"host, omitempty" column:"host"`
	Port       int64     `json:"port, omitempty" column:"port"`
	Status     bool      `json:"status, omitempty" column:"status"`
	CreateTime time.Time `json:"creationTime,omitempty" column:"created_time" gorm:"column:created_time"`
}

type APIServerSpec struct {
	Server APIServerInfor `json:"server, omitempty"`
}

type APIServer struct {
	unversioned.TypeMeta `json:",inline"`
	ObjectMeta           `json:"metadata,omitempty"`

	Spec APIServerSpec `json:"spec,omitempty"`
}

type APIServerList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []APIServer `json:"items"`
}

type UserServiceSpec struct {
	NodeUserReference map[string]UserReferences `json:"nodeUserReference,omitempty"`
	NodeCnt           uint                      `json:"nodecnt,omitempty"`
}

type UserService struct {
	unversioned.TypeMeta `json:",inline"`
	ObjectMeta           `json:"metadata,omitempty"`

	Spec UserServiceSpec `json:"spec,omitempty"`
}

type UserServiceList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []UserService `json:"spec,omitempty"`
}
