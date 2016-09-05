package shadowssapi

import (
	"crypto/rand"
	"fmt"
	"shadowsocks-go/pkg/backend/db"

	"github.com/golang/glog"
)

// var cache = cacheutil.NewCache(maxLoginCacheSize)
// var cacheIndex uint64

func randBearerToken() (string, error) {
	token := make([]byte, 16)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", token), err
}

func addToken(user *db.User) (string, error) {

	token, err := randBearerToken()

	err = Storage.AddUserToken(token, user.ID, user.MacAddr)
	if err != nil {
		glog.Errorf("store token failure %v", err)
	}

	return token, err
}

func CheckToken(token string) (*db.User, error) {
	uid, err := Storage.GetUserIDByToken(token)

	user, err := Storage.GetUserByID(uid)
	if err != nil {
		return nil, err
	}
	glog.V(5).Infof("Got Users %+v", user)

	return user, err
}
