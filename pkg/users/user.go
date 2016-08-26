package users

import "shadowsocks-go/pkg/storage"

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
}

func get(handle storage.Interface, startUserID, endUserID int64) ([]User, error) {
	var users []User
	ctx := createContextWithValue(userTableName)

	fileds := []string{"id", "passwd", "port", "method", "enable", "transfer_enable", "u", "d"}
	query := string("id BETWEEN ? AND ?")
	queryArgs := []interface{}{startUserID, endUserID}
	selection := NewSelection(fileds, query, queryArgs)

	err := handle.GetToList(ctx, selection, &users)
	return users, err
}

func updateTraffic(handle storage.Interface, userID int64, upload, download int64) error {

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
