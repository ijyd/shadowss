package db

import (
	"fmt"
	"shadowsocks-go/pkg/storage"
)

//User is a mysql users map
type User struct {
	ID              int64  `column:"id"`
	Port            int64  `column:"port"`
	Passwd          string `column:"passwd"`
	Method          string `column:"method"`
	Enable          int64  `column:"enable"`
	TrafficLimit    int64  `column:"transfer_enable"` //traffic for per user
	UploadTraffic   int64  `column:"u"`               //upload traffic for per user
	DownloadTraffic int64  `column:"d"`               //download traffic for per user
	Name            string `column:"user_name"`
	MacAddr         string `column:"macAddr"`
}

func GetServUsers(handle storage.Interface, startUserID, endUserID int64) ([]User, error) {
	var users []User
	ctx := createContextWithValue(userTableName)

	fileds := []string{"id", "passwd", "port", "method", "enable", "transfer_enable", "u", "d"}
	query := string("id BETWEEN ? AND ?")
	queryArgs := []interface{}{startUserID, endUserID}
	selection := NewSelection(fileds, query, queryArgs)

	err := handle.GetToList(ctx, selection, &users)
	return users, err
}

func GetUser(handle storage.Interface, key string) (*User, error) {
	var users []User
	ctx := createContextWithValue(userTableName)

	fileds := []string{"id", "passwd", "port", "method", "enable", "user_name", "macAddr"}
	query := string("user_name = ?")
	queryArgs := []interface{}{key}
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

func GetUserByID(handle storage.Interface, id int64) (*User, error) {
	var users []User
	ctx := createContextWithValue(userTableName)

	fileds := []string{"id", "passwd", "port", "method", "enable", "email"}
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

func UpdateUserTraffic(handle storage.Interface, userID int64, upload, download int64) error {

	user := &User{
		ID:              userID,
		UploadTraffic:   upload,
		DownloadTraffic: download,
	}

	conditionFields := string("id")
	updateFields := []string{"u", "d"}

	ctx := createContextWithValue(userTableName)
	err := handle.GuaranteedUpdate(ctx, conditionFields, updateFields, user)
	return err
}
