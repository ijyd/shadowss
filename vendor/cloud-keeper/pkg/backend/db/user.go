package db

import (
	"fmt"
	"time"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/pagination"
	"golib/pkg/storage"

	"github.com/golang/glog"
)

var fileds = []string{"id", "passwd", "email", "enable_ota",
	"traffic_limit", "upload", "download", "user_name", "manage_pass",
	"expire_time", "reg_ip", "reg_date", "description", "is_admin", "total_upload", "total_download", "status"}

func GetUserByID(handle storage.Interface, id int64) (*api.UserInfo, error) {
	var users []api.UserInfo
	ctx := createContextWithValue(userTableName)

	query := string("id = ?")
	queryArgs := []interface{}{id}
	selection := NewSelection(fileds, query, queryArgs)

	err := handle.GetToList(ctx, selection, &users)
	if err != nil {
		return nil, err
	}

	if len(users) > 0 {
		return &users[0], err
	} else {
		return nil, fmt.Errorf("not found")
	}

}

func GetUserByName(handle storage.Interface, key string) (*api.UserInfo, error) {
	var users []api.UserInfo
	ctx := createContextWithValue(userTableName)

	query := string("user_name = ?")
	queryArgs := []interface{}{key}
	selection := NewSelection(fileds, query, queryArgs)

	err := handle.GetToList(ctx, selection, &users)
	if err != nil {
		glog.Errorf("handle get to list %v\r\n", err)
		return nil, err
	}

	if len(users) > 0 {
		return &users[0], err
	} else {
		return nil, fmt.Errorf("not found")
	}

}

func GetUsersByLimit(handle storage.Interface, limitRecord int64, skip int64) ([]api.UserInfo, error) {
	var users []api.UserInfo
	ctx := createContextWithValue(userTableName)

	skipVal := skip
	sortVal := string("id")
	limitVal := limitRecord
	selection := NewPageSelection(fileds, nil, nil, sortVal, limitVal, skipVal)

	err := handle.GetToList(ctx, selection, &users)
	return users, err
}

func GetUserList(handle storage.Interface, page pagination.Pager) ([]api.UserInfo, error) {
	ctx := createContextWithValue(userTableName)
	selection, err := buildListSelecttion(ctx, handle, page, fileds)
	if err != nil {
		return nil, err
	}

	var userInfo []api.UserInfo
	err = handle.GetToList(ctx, selection, &userInfo)
	if err != nil {
		return nil, err
	}

	if len(userInfo) > 0 {
		return userInfo, nil
	} else {
		return nil, fmt.Errorf("not found")
	}
}

//func CreateUser(handle storage.Interface, name string, host string, port int64, isEnable bool) error {
func CreateUser(handle storage.Interface, info api.UserInfo) error {

	ctx := createContextWithValue(userTableName)

	var Enable int64
	if info.EnableOTA == 1 {
		Enable = 1
	} else {
		Enable = 0
	}

	user := &api.UserInfo{
		ID:           info.ID,
		Name:         info.Name,
		EnableOTA:    Enable,
		Email:        info.Email,
		TrafficLimit: info.TrafficLimit,
		ManagePasswd: info.ManagePasswd,
		Passwd:       info.Passwd,
		ExpireTime:   info.ExpireTime,
		EmailVerify:  info.EmailVerify,
		RegIPAddr:    info.RegIPAddr,
		RegDBTime:    time.Now(),
		TrafficRate:  info.TrafficRate,
		Description:  info.Description,
	}

	err := handle.Create(ctx, info.Name, user, user)
	if err != nil {
		glog.Errorf("create a server record failure %v\r\n", err)
	}
	return err
}

func UpdateUserTraffic(handle storage.Interface, userID int64, totalUpload, totalDownload, upload, download int64) error {

	user := &api.UserInfo{
		ID:                   userID,
		UploadTraffic:        upload,
		DownloadTraffic:      download,
		TotalDownloadTraffic: totalDownload,
		TotalUploadTraffic:   totalUpload,
	}

	conditionFields := string("id")
	updateFields := []string{"upload", "download", "total_upload", "total_download"}

	ctx := createContextWithValue(userTableName)
	err := handle.GuaranteedUpdate(ctx, conditionFields, updateFields, user)
	return err
}

func UpdateUserNode(handle storage.Interface, userID int64, nodes string) error {

	// user := &api.UserInfo{
	// 	ID:    userID,
	// 	Nodes: nodes,
	// }
	//
	// conditionFields := string("id")
	// updateFields := []string{"nodes"}
	//
	// ctx := createContextWithValue(userTableName)
	// err := handle.GuaranteedUpdate(ctx, conditionFields, updateFields, user)
	return nil
}

func UpdateUserPort(handle storage.Interface, userID int64, port int64) error {

	// user := &api.UserInfo{
	// 	ID:   userID,
	// 	Port: port,
	// }
	//
	// conditionFields := string("id")
	// updateFields := []string{"port"}
	//
	// ctx := createContextWithValue(userTableName)
	// err := handle.GuaranteedUpdate(ctx, conditionFields, updateFields, user)
	return nil
}

func UpdateUserStatus(handle storage.Interface, userID int64, status bool) error {

	statusInt := 0
	if status {
		statusInt = 1
	}

	user := &api.UserInfo{
		ID:     userID,
		Status: int64(statusInt),
	}

	conditionFields := string("id")
	updateFields := []string{"status"}

	ctx := createContextWithValue(userTableName)
	err := handle.GuaranteedUpdate(ctx, conditionFields, updateFields, user)
	return err
}

func DeleteUserByName(handle storage.Interface, name string) error {

	ctx := createContextWithValue(userTableName)

	var user api.UserInfo
	query := string("user_name = ?")
	queryArgs := []interface{}{name}
	selection := NewSelection(nil, query, queryArgs)

	err := handle.Delete(ctx, selection, &user)
	if err != nil {
		glog.Errorf("delete a server record failure %v\r\n", err)
	}
	return err
}
