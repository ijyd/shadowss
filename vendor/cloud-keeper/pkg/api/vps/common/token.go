package common

import (
	"crypto/rand"
	"fmt"
	"strings"

	"cloud-keeper/pkg/api"
	cacheutil "golib/pkg/util/cache"

	"github.com/golang/glog"
)

const (
	maxLoginCacheSize = 32 * 4
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

func AddToken(user *api.UserInfo) (string, error) {

	token, err := randBearerToken()
	if err != nil {
		glog.Errorf("generate token failure %v\r\n", err)
		return string(""), nil
	}

	id, err := Storage.GetTokenIDByUID(user.ID)
	if err != nil {

		err = Storage.AddUserToken(token, user.ID, user.Name)
		if err != nil {
			glog.Errorf("store token failure %v", err)
		} else {
			cache.Add(uint64(user.ID), user)
		}
	} else {
		Storage.UpdateToken(token, id)
	}

	return token, err
}

func CheckToken(input string) (*api.UserInfo, error) {

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

	glog.V(5).Infof("get user id %v\r\n", uid)
	var user *api.UserInfo
	obj, found := cache.Get(uint64(uid))
	if !found {
		user, err = Storage.GetUserByID(uid)
		if err != nil {
			return nil, err
		}
		glog.V(5).Infof("Got Users %+v by db, missing cache\r\n", user)
	} else {
		user, _ = obj.(*api.UserInfo)
	}

	return user, err
}
