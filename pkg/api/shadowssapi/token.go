package shadowssapi

import (
	"crypto/rand"
	"fmt"
	"shadowsocks-go/pkg/backend/db"
	cacheutil "shadowsocks-go/pkg/util/cache"
	"strings"

	"github.com/golang/glog"
)

var cache = cacheutil.NewCache(maxLoginCacheSize)
var cacheIndex uint64

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
	} else {
		cache.Add(uint64(user.ID), user)
	}

	return token, err
}

func CheckToken(input string) (*db.User, error) {

	parts := strings.Split(input, " ")
	if len(parts) < 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, fmt.Errorf("invalid token")
	}

	token := parts[1]
	glog.V(5).Infof("check token %v\r\n", token)
	uid, err := Storage.GetUserIDByToken(token)
	if err != nil {
		return nil, err
	}

	var user *db.User
	obj, found := cache.Get(uint64(uid))
	if !found {
		user, err = Storage.GetUserByID(uid)
		if err != nil {
			return nil, err
		}
		glog.V(5).Infof("Got Users %+v by db, missing cache\r\n", user)
	} else {
		user, _ = obj.(*db.User)
	}

	return user, err
}
