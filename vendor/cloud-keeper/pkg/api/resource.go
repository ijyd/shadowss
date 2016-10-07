package api

import (
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/api/unversioned"
	"golib/pkg/util/timewrap"
	"time"
)

type LoginSpec struct {
	AuthName string `json:"authname,omitempty"`
	Auth     string `json:"auth,omitempty"`
	AuthID   string `json:"authID,omitempty"`
	Token    string `json:"token,omitempty"`
}

type Login struct {
	unversioned.TypeMeta `json:",inline"`
	prototype.ObjectMeta `json:"metadata,omitempty"`

	Spec LoginSpec `json:"spec,omitempty"`
}

type AccServer struct {
	unversioned.TypeMeta `json:",inline"`

	ID       string `json:"id"`
	Size     string `json:"size,omitempty"`
	Region   string `json:"region,omitempty"`
	Image    string `json:"image,omitempty"`
	SSHKeyID string `json:"sshKeyID,omitempty"`
	Name     string `json:"name,omitempty"`

	Information map[string]interface{} `json:"info,omitempty"`
}

const (
	CNISPCMCC    = "cnISPCMCC"
	CNISPUNICOM  = "cnISPUnicom"
	CNISPASPCTCC = "cnISPCTCC"
	CNISPOther   = "cnISPOther"
)

const (
	NodeUserSpaceDefault = "default"
	NodeUserSpaceAPI     = "api"
	NodeUserSpaceDev     = "develop"
)

const (
	NodeLablesChinaISP    = "cnISP"
	NodeLablesUserSpace   = "userSpace"
	NodeLablesVPSLocation = "vpsLocation"
	NodeLablesVPSOP       = "vpsOperator"
	NodeLablesVPSName     = "vpsName"
	NodeLablesVPSIP       = "vpsIP"
)

type AccServerDeploySS struct {
	HostList  []string          `json:"hostList,omitempty"`
	Attribute map[string]string `json:"attr,omitempty"`
}

type AccServerCommand struct {
	unversioned.TypeMeta `json:",inline"`

	SSHKey  string            `json:"sshKey,omitempty"`
	Command string            `json:"command,omitempty"`
	Deploy  AccServerDeploySS `json:"deploySS,omitempty"`
}

type AccServerSSHKey struct {
	unversioned.TypeMeta `json:",inline"`

	Key interface{} `json:"key,omitempty"`
}

type AccServerList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []AccServer `json:"items"`
}

type AccountInfo struct {
	unversioned.TypeMeta `json:",inline"`

	Information map[string]interface{} `json:"info,omitempty"`
}

const (
	OperatorVultr        = "Vultr"
	OperatorDigitalOcean = "DigitalOcean"
)

type OperatorType string

type AccountDetail struct {
	ID             int64         `json:"-" column:"id"`
	Name           string        `json:"name,omitempty" column:"name"`
	Operators      string        `json:"operators,omitempty" column:"operators"`
	Key            string        `json:"key,omitempty" column:"api_key" gorm:"column:api_key"`
	Descryption    string        `json:"descryption,omitempty" column:"descryption"`
	CreditCeilings float64       `json:"creditCeilings,omitempty" column:"credit_ceilings"`
	Lables         string        `json:"lables,omitempty" column:"lables"`
	CreateTime     timewrap.Time `json:"creationTime,omitempty"`
	ExpireTime     timewrap.Time `json:"expire,omitempty"`
	//a trick for database use
	ExpireDBTime time.Time `json:"-" column:"expire_time" gorm:"column:expire_time"`
	CreateDBTime time.Time `json:"-" column:"created_time" gorm:"column:created_time"`
}

// AccountSpec of Vultr account
type AccountSpec struct {
	AccDetail AccountDetail `json:"account,omitempty"`
}

type Account struct {
	unversioned.TypeMeta `json:",inline"`
	prototype.ObjectMeta `json:"metadata,omitempty"`

	Spec AccountSpec `json:"spec,omitempty"`
}

type AccountList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []Account `json:"items"`
}

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

const (
	NodeUserPhaseAdd    = "add"
	NodeUserPhaseDelete = "del"
	NodeUserPhaseUpdate = "update"
)

type NodeUserPhase string

type NodeUserSpec struct {
	User     UserReferences `json:"user,omitempty"`
	NodeName string         `json:"nodeName,omitempty"`
	Phase    NodeUserPhase  `json:"phase,omitempty"`
}

//put your user into your node with node name
//like as /api/node/{nodename}/nodeuser/{username}
type NodeUser struct {
	unversioned.TypeMeta `json:",inline"`
	prototype.ObjectMeta `json:"metadata,omitempty"`

	Spec NodeUserSpec `json:"spec,omitempty"`
}

type NodeUserList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []NodeUser `json:"items,omitempty"`
}

type NodeServer struct {
	ID                   int64  `json:"id" column:"id" gorm:"column:id"`
	Name                 string `json:"name,omitempty" column:"name" gorm:"column:name"`
	EnableOTA            int64  `json:"enableOTA" column:"enableota" gorm:"column:enableota"`
	Host                 string `json:"host,omitempty" column:"server" gorm:"column:server"`
	Method               string `json:"method" column:"method" gorm:"column:method"`
	Status               int64  `json:"status,omitempty" column:"status" gorm:"column:status"`
	Location             string `json:"location,omitempty" column:"location" gorm:"column:location"`
	AccServerID          int64  `json:"accServerID,omitempty" column:"vps_server_id" gorm:"column:vps_server_id"`
	AccServerName        string `json:"accServerName,omitempty" column:"vps_server_name" gorm:"column:vps_server_name"`
	Description          string `json:"description,omitempty" column:"description" gorm:"column:description"`
	TrafficLimit         int64  `json:"trafficLimit,omitempty" column:"traffic_limit" gorm:"column:traffic_limit"`
	Upload               int64  `json:"upload,omitempty" column:"upload" gorm:"column:upload"`
	Download             int64  `json:"download,omitempty" column:"download" gorm:"column:download"`
	TrafficRate          int64  `json:"trafficRate,omitempty" column:"traffic_rate" gorm:"column:traffic_rate"`
	TotalUploadTraffic   int64  `json:"totalUploadTraffic,omitempty" column:"total_upload" gorm:"column:total_upload"`
	TotalDownloadTraffic int64  `json:"totalDownloadTraffic,omitempty" column:"total_download" gorm:"column:total_download"`
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
	Server   APIServerInfor `json:"server, omitempty"`
	HostList []string       `json:"hosts, omitempty"`
}

type APIServer struct {
	unversioned.TypeMeta `json:",inline"`
	prototype.ObjectMeta `json:"metadata,omitempty"`

	Spec APIServerSpec `json:"spec,omitempty"`
}

type APIServerList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []APIServer `json:"items"`
}

type NodeReferences struct {
	Host string         `json:"host,omitempty"`
	User UserReferences `json:"user,omitempty"`
}

type UserServiceSpec struct {
	NodeUserReference map[string]NodeReferences `json:"nodeUserReference,omitempty"`
	NodeCnt           uint                      `json:"nodecnt,omitempty"`
	Status            bool                      `json:"status,omitempty"`
}

type UserService struct {
	unversioned.TypeMeta `json:",inline"`
	prototype.ObjectMeta `json:"metadata,omitempty"`

	Spec UserServiceSpec `json:"spec,omitempty"`
}

type UserServiceList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []UserService `json:"items"`
}

//User is a mysql users map
type UserInfo struct {
	ID                   int64     `json:"id,omitempty" column:"id"`
	Passwd               string    `json:"passwd,omitempty" column:"passwd"`
	Email                string    `json:"email,omitempty" column:"email"`
	EnableOTA            int64     `json:"enableOTA,omitempty" column:"enable_ota"`
	TrafficLimit         int64     `json:"trafficLimit,omitempty" column:"traffic_limit" gorm:"column:traffic_limit"` //traffic for per user
	UploadTraffic        int64     `json:"uploadTraffic,omitempty" column:"upload" gorm:"column:upload"`              //upload traffic for per user
	DownloadTraffic      int64     `json:"downloadTraffic,omitempty" column:"download" gorm:"column:download"`        //download traffic for per user
	Name                 string    `json:"name,omitempty" column:"user_name" gorm:"column:user_name"`
	ManagePasswd         string    `json:"managePasswd,omitempty" column:"manage_pass" gorm:"column:manage_pass"`
	ExpireTime           time.Time `json:"expireTime,omitempty" column:"expire_time" gorm:"column:expire_time"`
	EmailVerify          int16     `json:"emailVerify,omitempty" column:"is_email_verify" gorm:"column:is_email_verify"`
	RegIPAddr            string    `json:"regIPAddr,omitempty" column:"reg_ip" gorm:"column:reg_ip"`
	RegDBTime            time.Time `json:"regTime,omitempty" column:"reg_date" gorm:"column:reg_date"`
	Description          string    `json:"description,omitempty" column:"description" gorm:"column:description"`
	TrafficRate          float64   `json:"trafficRate,omitempty" column:"traffic_rate" gorm:"column:traffic_rate"`
	IsAdmin              int64     `json:"isAdmin,omitempty" column:"is_admin" gorm:"column:is_admin"`
	LastCheckInTime      time.Time `json:"-" column:"last_check_in_time" gorm:"column:last_check_in_time"`
	LastResetPwdTime     time.Time `json:"-" column:"last_reset_pass_time" gorm:"column:last_reset_pass_time"`
	TotalUploadTraffic   int64     `json:"totalUploadTraffic,omitempty" column:"total_upload" gorm:"column:total_upload"`
	TotalDownloadTraffic int64     `json:"totalDownloadTraffic,omitempty" column:"total_download" gorm:"column:total_download"`
	Status               int64     `json:"status,omitempty" column:"status" gorm:"column:status"`
}

type UserSpec struct {
	DetailInfo UserInfo `json:"detailInfo,omitempty"`
}

type User struct {
	unversioned.TypeMeta `json:",inline"`
	prototype.ObjectMeta `json:"metadata,omitempty"`

	Spec UserSpec `json:"spec,omitempty"`
}

type UserList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []User `json:"spec,omitempty"`
}
