package apiext

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
	AuthName string    `json:"authname,omitempty"`
	Auth     string    `json:"auth,ommitempty"`
	Token    string    `json:"token,omitempty"`
	User     LoginUser `json:"user,omitempty" freezer:"table:usertest"`
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
