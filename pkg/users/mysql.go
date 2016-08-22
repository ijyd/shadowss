package users

import (
	"shadowsocks-go/pkg/storage"
	"shadowsocks-go/pkg/storage/storagebackend"
	"shadowsocks-go/pkg/storage/storagebackend/factory"
)

//User is a mysql users map
type User struct {
	ID              int    `column:"id"`
	Port            int    `column:"port"`
	Passwd          string `column:"passwd"`
	Method          string `column:"method"`
	Enable          int    `column:"enable"`
	TrafficLimit    int    `column:"transfer_enable"` //traffic for per user
	UploadTraffic   int    `column:"u"`               //upload traffic for per user
	DownloadTraffic int    `column:"d"`               //download traffic for per user
}

const (
	userTableName = "user"
)

var userTableField = []string{"id", "passwd", "port", "method", "enable", "transfer_enable", "u", "d"}

func newStorage(c storagebackend.Config) (storage.Interface, error) {
	return factory.Create(c)
}

func get(handle storage.Interface) ([]User, error) {
	var users []User
	err := handle.GetToList(userTableName, userTableField, &users)
	return users, err
}

func updateTraffic(handle storage.Interface, userID int, upload, download int64) error {

	user := &User{
		ID:              userID,
		UploadTraffic:   int(upload),
		DownloadTraffic: int(download),
	}

	conditionFields := []string{"id"}
	updateFields := []string{"u", "d"}

	err := handle.GuaranteedUpdate(userTableName, conditionFields, updateFields, user)
	return err
}
