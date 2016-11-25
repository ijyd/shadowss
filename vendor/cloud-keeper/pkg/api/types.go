package api

import (
	prototype "gofreezer/pkg/api"
	"gofreezer/pkg/api/unversioned"
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

type LoginList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []Login `json:"items"`
}

//User is a mysql users map
type UserTokenSpec struct {
	ID         int64            `json:"id,omitempty" freezer:"column:id" gorm:"column:id"`
	Token      string           `json:"token,omitempty" freezer:"column:token" gorm:"column:token"`
	UserID     int64            `json:"userID,omitempty" freezer:"column:user_id" gorm:"column:user_id"`
	CreateTime unversioned.Time `json:"createTime,omitempty" freezer:"column:create_time" gorm:"column:create_time"`
	ExpireTime unversioned.Time `json:"expireTime,omitempty" freezer:"column:expire_time" gorm:"column:expire_time"`
	Name       string           `json:"name,omitempty" freezer:"column:name;resoucekey" gorm:"column:name"`
}

type UserToken struct {
	unversioned.TypeMeta `json:",inline"`
	prototype.ObjectMeta `json:"metadata,omitempty"`

	Spec UserTokenSpec `json:"spec,omitempty" freezer:"table:user_token"`
}

type UserTokenList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []UserToken `json:"items"`
}

type VultrServerInfo struct {
	CreatedTime string `json:"createdTime,omitempty"`
	Location    string `json:"location,omitempty"`
	Name        string `json:"name,omitempty"`
	Status      string `json:"status,omitempty"`

	IPV4Addr    string `json:"ipv4Addr,omitempty"`
	IPV4NetMask string `json:"ipv4NetMask,omitempty"`
	IPV4Gateway string `json:"ipv4Gateway,omitempty"`

	PendingCharges float64 `json:"pendingCharges,omitempty"`

	CostPerMonth     string  `json:"costPerMonth,omitempty"`
	AllowedBandWidth float64 `json:"allowedBandwidth,omitempty"`
	CurrentBandwidth float64 `json:"currentBandwidth,omitempty"`
}

type DGServerInfo struct {
	Location    string `json:"location,omitempty"`
	Name        string `json:"name,omitempty"`
	Status      string `json:"status,omitempty"`
	CreatedTime string `json:"createdTime,omitempty"`

	IPV4Addr    string `json:"ipv4Addr,omitempty"`
	IPV4NetMask string `json:"ipv4NetMask,omitempty"`
	IPV4Gateway string `json:"ipv4Gateway,omitempty"`

	PriceMonthly float64 `json:"priceMonthly,omitempty"`
	PriceHourly  float64 `json:"priceHourly,omitempty"`
}

type AccServerSpec struct {
	ID       string `json:"id,omitempty"`
	Size     string `json:"size,omitempty"`
	Region   string `json:"region,omitempty"`
	Image    string `json:"image,omitempty"`
	SSHKeyID string `json:"sshKeyID,omitempty"`
	Name     string `json:"name,omitempty"`

	AccName string `json:"accName,omitempty"`

	DigitalOcean DGServerInfo    `json:"digitalocean,omitempty"`
	Vultr        VultrServerInfo `json:"vultr,omitempty"`
}

type AccServer struct {
	unversioned.TypeMeta `json:",inline"`
	prototype.ObjectMeta `json:"metadata,omitempty"`

	Spec AccServerSpec `json:"spec,omitempty"`
}

type AccServerList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []AccServer `json:"items"`
}

type AccServerDeploySS struct {
	HostList  []string          `json:"hostList,omitempty"`
	Attribute map[string]string `json:"attr,omitempty"`
}

type AccExecSpec struct {
	SSHKey  string            `json:"sshKey,omitempty"`
	Command string            `json:"command,omitempty"`
	ID      int64             `json:"id,omitempty"`
	Deploy  AccServerDeploySS `json:"deploySS,omitempty"`
	AccName string            `json:"accName,omitempty"`
}

type AccExec struct {
	unversioned.TypeMeta `json:",inline"`
	prototype.ObjectMeta `json:"metadata,omitempty"`

	Spec AccExecSpec `json:"spec,omitempty"`
}

type SSHKey struct {
	KeyID       string `json:"keyID,omitempty"`
	Name        string `json:"name,omitempty"`
	Key         string `json:"key,omitempty"`
	FingerPrint string `json:"fingerprint,omitempty"`
}

type AccSSHKeySpec struct {
	Keys    []SSHKey `json:"keys,omitempty"`
	AccName string   `json:"accName,omitempty"`
}

type AccSSHKey struct {
	unversioned.TypeMeta `json:",inline"`
	prototype.ObjectMeta `json:"metadata,omitempty"`

	Spec AccSSHKeySpec `json:"spec,omitempty"`
}

//
// type AccSSHKeyList struct {
// 	unversioned.TypeMeta `json:",inline"`
// 	unversioned.ListMeta `json:"metadata,omitempty"`
//
// 	Items []AccSSHKey `json:"items"`
// }

type VultrAccountInfo struct {
	Balance           float64 `json:"balance,omitempty"`
	PendingCharges    float64 `json:"pendingCharges,omitempty"`
	LastPaymentDate   string  `json:"lastPaymentDate,omitempty"`
	LastPaymentAmount float64 `json:"lastPaymentAmount,omitempty"`
}

type DGAccountInfo struct {
	DropletLimit  int    `json:"dropletLimit,omitempty"`
	Email         string `json:"email,omitempty"`
	UUID          string `json:"uuid,omitempty"`
	EmailVerified bool   `json:"emailVerified,omitempty"`
}

type AccountInfoSpec struct {
	DigitalOcean DGAccountInfo    `json:"digitalocean,omitempty"`
	Vultr        VultrAccountInfo `json:"vultr,omitempty"`
	AccName      string           `json:"accName,omitempty"`
}

type AccountInfo struct {
	unversioned.TypeMeta `json:",inline"`
	prototype.ObjectMeta `json:"metadata,omitempty"`

	Spec AccountInfoSpec `json:"spec,omitempty"`
}

const (
	OperatorVultr        = "Vultr"
	OperatorDigitalOcean = "DigitalOcean"
)

type OperatorType string

type AccountDetail struct {
	ID             int64            `json:"-" freezer:"column:id"`
	Name           string           `json:"name,omitempty" freezer:"column:name;resoucekey"`
	Operators      string           `json:"operators,omitempty" freezer:"column:operators"`
	Key            string           `json:"key,omitempty" freezer:"column:api_key" gorm:"column:api_key"`
	Descryption    string           `json:"descryption,omitempty" freezer:"column:descryption"`
	CreditCeilings float64          `json:"creditCeilings,omitempty" freezer:"column:credit_ceilings"`
	Lables         string           `json:"lables,omitempty" freezer:"column:lables"`
	CreateTime     unversioned.Time `json:"creationTime,omitempty" freezer:"column:expire_time" gorm:"column:expire_time"`
	ExpireTime     unversioned.Time `json:"expire,omitempty" freezer:"column:created_time" gorm:"column:created_time"`
}

// AccountSpec of Vultr account
type AccountSpec struct {
	AccDetail AccountDetail `json:"account,omitempty" freezer:"table:vps_server_account"`
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
	NodeLablesVPSID       = "vpsID"
	NodeLablesVPSOP       = "vpsOperator"
	NodeLablesVPSName     = "vpsName"
	NodeLablesVPSIP       = "vpsIP"
)

const (
	NodeAnnotationUserCnt    = "userCount"
	NodeAnnotationRefreshCnt = "Refresh"
	NodeAnnotationVersion    = "version"
)

type NodeServer struct {
	ID                   int64  `json:"id" freezer:"column:id" gorm:"column:id"`
	Name                 string `json:"name,omitempty" freezer:"column:name;resoucekey" gorm:"column:name"`
	EnableOTA            bool   `json:"enableOTA" freezer:"column:enableota" gorm:"column:enableota"`
	Host                 string `json:"host,omitempty" freezer:"column:server" gorm:"column:server"`
	Method               string `json:"method" freezer:"column:method" gorm:"column:method"`
	Status               int64  `json:"status,omitempty" freezer:"column:status" gorm:"column:status"`
	Location             string `json:"location,omitempty" freezer:"column:location" gorm:"column:location"`
	AccServerID          int64  `json:"accServerID,omitempty" freezer:"column:vps_server_id" gorm:"column:vps_server_id"`
	AccServerName        string `json:"accServerName,omitempty" freezer:"column:vps_server_name" gorm:"column:vps_server_name"`
	Description          string `json:"description,omitempty" freezer:"column:description" gorm:"column:description"`
	TrafficLimit         int64  `json:"trafficLimit,omitempty" freezer:"column:traffic_limit" gorm:"column:traffic_limit"`
	Upload               int64  `json:"upload,omitempty" freezer:"column:upload" gorm:"column:upload"`
	Download             int64  `json:"download,omitempty" freezer:"column:download" gorm:"column:download"`
	TrafficRate          int64  `json:"trafficRate,omitempty" freezer:"column:traffic_rate" gorm:"column:traffic_rate"`
	TotalUploadTraffic   int64  `json:"totalUploadTraffic,omitempty" freezer:"column:total_upload" gorm:"column:total_upload"`
	TotalDownloadTraffic int64  `json:"totalDownloadTraffic,omitempty" freezer:"column:total_download" gorm:"column:total_download"`
	CustomMethod         int    `json:"customMethod,omitempty" freezer:"column:custom_method" gorm:"column:custom_method"`
}

type NodeSpec struct {
	Server NodeServer              `json:"server,omitempty" freezer:"table:vps_new_node"`
	Users  map[string]NodeUserSpec `json:"users,omitempty"`
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

// const (
// 	UserServicetDefaultNode = "default"
// )

type UserServiceSpec struct {
	NodeName  string         `json:"nodeName,omitempty"`
	Host      string         `json:"host,omitempty"`
	UserRefer UserReferences `json:"userRefer,omitempty"`
	Delete    bool           `json:"delete,omitempty"`
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

type UserServiceBindingNodesSpec struct {
	NodeUserReference map[string]NodeReferences `json:"nodeUserReference,omitempty"`
	NodeCnt           uint                      `json:"nodecnt,omitempty"`
	Status            bool                      `json:"status,omitempty"`
}

type UserServiceBindingNodes struct {
	unversioned.TypeMeta `json:",inline"`
	prototype.ObjectMeta `json:"metadata,omitempty"`

	Spec UserServiceBindingNodesSpec `json:"spec,omitempty"`
}

type UserInfo struct {
	ID                   int64            `json:"id,omitempty" freezer:"column:id"`
	Passwd               string           `json:"passwd,omitempty" freezer:"column:passwd"`
	Email                string           `json:"email,omitempty" freezer:"column:email"`
	EnableOTA            bool             `json:"enableOTA,omitempty" freezer:"column:enable_ota"`
	TrafficLimit         int64            `json:"trafficLimit,omitempty" freezer:"column:traffic_limit" gorm:"column:traffic_limit"` //traffic for per user
	UploadTraffic        int64            `json:"uploadTraffic,omitempty" freezer:"column:upload" gorm:"column:upload"`              //upload traffic for per user
	DownloadTraffic      int64            `json:"downloadTraffic,omitempty" freezer:"column:download" gorm:"column:download"`        //download traffic for per user
	Name                 string           `json:"name,omitempty" freezer:"column:user_name;resoucekey" gorm:"column:user_name"`
	ManagePasswd         string           `json:"managePasswd,omitempty" freezer:"column:manage_pass" gorm:"column:manage_pass"`
	ExpireTime           unversioned.Time `json:"expireTime,omitempty" freezer:"column:expire_time" gorm:"column:expire_time"`
	EmailVerify          bool             `json:"emailVerify,omitempty" freezer:"column:is_email_verify" gorm:"column:is_email_verify"`
	RegIPAddr            string           `json:"regIPAddr,omitempty" freezer:"column:reg_ip" gorm:"column:reg_ip"`
	RegDBTime            unversioned.Time `json:"regTime,omitempty" freezer:"column:reg_date" gorm:"column:reg_date"`
	Description          string           `json:"description,omitempty" freezer:"column:description" gorm:"column:description"`
	TrafficRate          float64          `json:"trafficRate,omitempty" freezer:"column:traffic_rate" gorm:"column:traffic_rate"`
	IsAdmin              bool             `json:"isAdmin,omitempty" freezer:"column:is_admin" gorm:"column:is_admin"`
	LastCheckInTime      unversioned.Time `json:"-" freezer:"column:last_check_in_time" gorm:"column:last_check_in_time"`
	LastResetPwdTime     unversioned.Time `json:"-" freezer:"column:last_reset_pass_time" gorm:"column:last_reset_pass_time"`
	TotalUploadTraffic   int64            `json:"totalUploadTraffic,omitempty" freezer:"column:total_upload" gorm:"column:total_upload"`
	TotalDownloadTraffic int64            `json:"totalDownloadTraffic,omitempty" freezer:"column:total_download" gorm:"column:total_download"`
	Status               bool             `json:"status,omitempty" freezer:"column:status" gorm:"column:status"`
}

type NodeReferences struct {
	Host string         `json:"host,omitempty"`
	User UserReferences `json:"user,omitempty"`
}

type BindingNodes struct {
	Nodes   map[string]NodeReferences `json:"nodes,omitempty"`
	NodeCnt uint                      `json:"nodecnt,omitempty"`
	Status  bool                      `json:"status,omitempty"`
}

type UserSpec struct {
	DetailInfo  UserInfo     `json:"detailInfo,omitempty" freezer:"table:vps_new_user"`
	UserService BindingNodes `json:"userService,omitempty"`
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

type UserPublicFileSpec struct {
	FileName    string `json:"file,omitempty"`
	Description string `json:"description,omitempty"`
}

type UserPublicFile struct {
	unversioned.TypeMeta `json:",inline"`
	prototype.ObjectMeta `json:"metadata,omitempty"`

	Spec UserPublicFileSpec `json:"spec,omitempty"`
}

type UserPublicFileList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []UserPublicFile `json:"items,omitempty"`
}

type ActiveAPINodeSpec struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Password string `jsong:"pwd,omitempty"`
	Method   string `json:"method,omitempty"`
}

type ActiveAPINode struct {
	unversioned.TypeMeta `json:",inline"`
	prototype.ObjectMeta `json:"metadata,omitempty"`

	Spec ActiveAPINodeSpec `json:"spec,omitempty"`
}

type ActiveAPINodeList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items       []ActiveAPINode `json:"items,omitempty"`
	EncryptData string          `json:"encData,omitempty"`
}
