package db

import (
	"fmt"
	"shadowsocks-go/pkg/storage"
	"time"

	"github.com/golang/glog"
)

//User is a mysql users map
type UserToken struct {
	ID         int64  `column:"id"`
	Token      string `column:"token"`
	UserID     int64  `column:"user_id"`
	CreateTime int    `column:"create_time"`
	ExpireTime int    `column:"expire_time"`
	MacAddr    string `column:"macAddr" gorm:"column:macAddr"`
}

func CheckToken(handle storage.Interface, token string) (int64, error) {

	fileds := []string{"id", "token", "user_id", "create_time", "expire_time", "macAddr"}
	query := string("token = ?")
	queryArgs := []interface{}{token}
	selection := NewSelection(fileds, query, queryArgs)

	ctx := createContextWithValue(userTokeTableName)

	var usertoken []UserToken
	err := handle.GetToList(ctx, selection, &usertoken)
	if err != nil {
		return 0, err
	}

	if len(usertoken) > 0 {
		return usertoken[0].UserID, nil
	} else {
		return 0, fmt.Errorf("not found")
	}

}

func CreateToken(handle storage.Interface, token string, uid int64, macAddr string) error {

	ctx := createContextWithValue(userTokeTableName)
	userToken := &UserToken{
		Token:      token,
		UserID:     uid,
		MacAddr:    macAddr,
		CreateTime: int(time.Now().Unix()),
		ExpireTime: int(time.Now().Add(time.Duration(1) * time.Hour).Unix()),
	}

	err := handle.Create(ctx, token, userToken, userToken)
	if err != nil {
		glog.Errorf("create a token record failure %v\r\n", err)
	}
	return err
}
