package api

import (
	"gofreezer/pkg/api"
	"gofreezer/pkg/api/unversioned"
	"time"
)

/*DROP TABLE IF EXISTS `usertest`;
CREATE TABLE `usertest` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,
  `status` TINYINT(1) NOT NULL DEFAULT '1',
	`login_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	`count` bigint(63) NOT NULL DEFAULT '0',
  UNIQUE KEY `server` (`name`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;*/

type LoginUser struct {
	ID        int64     `json:"id,omitempty" freezer:"column:id" gorm:"column:id"`
	UserName  string    `json:"userName,omitempty" freezer:"column:name" gorm:"column:name"`
	LoginTime time.Time `json:"loginTime,omitempty" freezer:"column:login_time" gorm:"column:login_time"`
	Count     int32     `json:"count,omitempty" freezer:"column:count" gorm:"column:count"`
	Status    bool      `json:"status,omitempty" freezer:"column:status" gorm:"column:status"`
}

type LoginSpec struct {
	AuthName string `json:"authname,omitempty"`
	Auth     string `json:"auth,ommitempty"`
	Token    string `json:"token,omitempty"`
}

type Login struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata,omitempty"`

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
	api.ObjectMeta       `json:"metadata,omitempty"`

	Spec UserTokenSpec `json:"spec,omitempty" freezer:"table:user_token"`
}

type UserTokenList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []UserToken `json:"items"`
}

//User is a mysql users map
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

type UserSpec struct {
	DetailInfo UserInfo `json:"detailInfo,omitempty" freezer:"table:user"`
}

type User struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata,omitempty"`

	Spec UserSpec `json:"spec,omitempty"`
}

type UserList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []User `json:"spec,omitempty"`
}
