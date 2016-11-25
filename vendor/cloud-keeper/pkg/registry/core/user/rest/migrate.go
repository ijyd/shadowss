package rest

import (
	"cloud-keeper/pkg/api"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gofreezer/pkg/api/unversioned"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	freezerapi "gofreezer/pkg/api"

	"github.com/golang/glog"
)

//User is a mysql users map
type UserInfo1 struct {
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

type UserSpec1 struct {
	DetailInfo UserInfo1 `json:"detailInfo,omitempty"`
}

type User1 struct {
	unversioned.TypeMeta  `json:",inline"`
	freezerapi.ObjectMeta `json:"metadata,omitempty"`

	Spec UserSpec1 `json:"spec,omitempty"`
}

type User1List struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []User1 `json:"spec,omitempty"`
}

func RequestUserListFromEtcd() (*User1List, error) {
	secure := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
		DisableKeepAlives:  true,
	}
	timeout := 5 * time.Second
	client := &http.Client{
		Transport: secure,
		Timeout:   timeout,
	}

	url := fmt.Sprintf("https://54.191.184.140:18088/api/v1/users")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer d64eed709f69bef52b9c828e73198dc6")
	//req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("get usrs  error:%v", string(body))
	}
	userlist := &User1List{}

	err = json.Unmarshal(body, userlist)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error:%v", err)
	}

	return userlist, nil
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
	unversioned.TypeMeta  `json:",inline"`
	freezerapi.ObjectMeta `json:"metadata,omitempty"`

	Spec UserServiceSpec `json:"spec,omitempty"`
}

func RequestUserServiceFromEtcd(name string) (*UserService, error) {
	secure := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
		DisableKeepAlives:  true,
	}
	timeout := 5 * time.Second
	client := &http.Client{
		Transport: secure,
		Timeout:   timeout,
	}

	url := fmt.Sprintf("https://54.191.184.140:18088/api/v1/users/%s/bindingnodes", name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer d64eed709f69bef52b9c828e73198dc6")
	//req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("get usrs  error:%v", string(body))
	}
	usersrv := &UserService{}

	err = json.Unmarshal(body, usersrv)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error:%v", err)
	}

	return usersrv, nil
}

func (r *UserREST) RequestUserFromEtcd() error {
	userlist, err := RequestUserListFromEtcd()
	if err != nil {
		return err
	}

	for k, v := range userlist.Items {

		usersrv, err := RequestUserServiceFromEtcd(v.Spec.DetailInfo.Name)
		if err != nil {
			return err
		}
		glog.V(5).Infof("got %v:%+v\r\n", k, *usersrv)
		newUserName := strings.Replace(v.Spec.DetailInfo.Name, ":", "", -1)
		newUserName = strings.ToLower(newUserName)

		ctx := freezerapi.NewContext()
		user := &api.User{}
		user.Kind = "User"
		user.APIVersion = "v1"
		user.Name = newUserName

		user.Annotations = make(map[string]string)
		for ak, av := range v.Annotations {
			user.Annotations[ak] = av
		}

		user.Spec.DetailInfo.Description = v.Spec.DetailInfo.Description
		user.Spec.DetailInfo.DownloadTraffic = v.Spec.DetailInfo.DownloadTraffic
		user.Spec.DetailInfo.Email = v.Spec.DetailInfo.Email
		if v.Spec.DetailInfo.EmailVerify == 0 {
			user.Spec.DetailInfo.EmailVerify = false
		} else {
			user.Spec.DetailInfo.EmailVerify = true
		}

		if v.Spec.DetailInfo.EnableOTA == 0 {
			user.Spec.DetailInfo.EnableOTA = false
		} else {
			user.Spec.DetailInfo.EnableOTA = true
		}

		user.Spec.DetailInfo.ExpireTime = unversioned.NewTime(v.Spec.DetailInfo.ExpireTime)
		user.Spec.DetailInfo.ID = v.Spec.DetailInfo.ID
		if v.Spec.DetailInfo.IsAdmin == 0 {
			user.Spec.DetailInfo.IsAdmin = false
		} else {
			user.Spec.DetailInfo.IsAdmin = true
		}

		user.Spec.DetailInfo.LastCheckInTime = unversioned.NewTime(v.Spec.DetailInfo.LastCheckInTime)
		user.Spec.DetailInfo.RegDBTime = unversioned.NewTime(v.Spec.DetailInfo.RegDBTime)
		user.Spec.DetailInfo.LastResetPwdTime = unversioned.NewTime(v.Spec.DetailInfo.LastResetPwdTime)
		user.Spec.DetailInfo.ManagePasswd = v.Spec.DetailInfo.ManagePasswd
		user.Spec.DetailInfo.Name = newUserName
		user.Spec.DetailInfo.Passwd = v.Spec.DetailInfo.Passwd
		user.Spec.DetailInfo.RegIPAddr = v.Spec.DetailInfo.RegIPAddr
		if v.Spec.DetailInfo.Status == 0 {
			user.Spec.DetailInfo.Status = false
		} else {
			user.Spec.DetailInfo.Status = true
		}

		user.Spec.DetailInfo.TotalDownloadTraffic = v.Spec.DetailInfo.TotalDownloadTraffic
		user.Spec.DetailInfo.TotalUploadTraffic = v.Spec.DetailInfo.TotalUploadTraffic
		user.Spec.DetailInfo.UploadTraffic = v.Spec.DetailInfo.UploadTraffic
		user.Spec.DetailInfo.TrafficLimit = v.Spec.DetailInfo.TrafficLimit
		user.Spec.DetailInfo.TrafficRate = v.Spec.DetailInfo.TrafficRate

		user.Spec.UserService.NodeCnt = usersrv.Spec.NodeCnt
		user.Spec.UserService.Status = usersrv.Spec.Status

		user.Spec.UserService.Nodes = make(map[string]api.NodeReferences)
		for k, v := range usersrv.Spec.NodeUserReference {
			nodeName := strings.Replace(k, ":", "", -1)
			nodeRefer := api.NodeReferences{}
			nodeRefer.User.DownloadTraffic = v.User.DownloadTraffic
			nodeRefer.User.UploadTraffic = v.User.UploadTraffic
			nodeRefer.User.EnableOTA = v.User.EnableOTA
			nodeRefer.User.ID = v.User.ID
			nodeRefer.User.Method = v.User.Method
			nodeRefer.User.Name = newUserName
			nodeRefer.User.Password = v.User.Password
			nodeRefer.User.Port = v.User.Port
			nodeRefer.Host = v.Host
			user.Spec.UserService.Nodes[nodeName] = nodeRefer
		}

		_, err = r.dynamo.Create(ctx, user)
		if err != nil {
			return err
		}
	}

	return nil

}

func (r *UserREST) MigrateUserToDynamodb() error {

	err := r.RequestUserFromEtcd()
	glog.Infof("request error %v\r\n", err)

	return nil
}
