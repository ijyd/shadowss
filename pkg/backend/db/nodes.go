package db

import (
	"shadowsocks-go/pkg/storage"
	"shadowsocks-go/pkg/util/network"

	"github.com/golang/glog"
)

//User is a mysql users map
type Node struct {
	ID          int64  `column:"id"`
	Name        string `column:"name"`
	Type        int64  `column:"type"`
	Host        string `column:"server"`
	Method      string `column:"method"`
	Status      string `column:"status"`      //traffic for per user
	StartUserID int64  `column:"startUserID"` //upload traffic for per user
	EndUserID   int64  `column:"endUserID"`   //download traffic for per user
}

func GetNodesByUserID(handle storage.Interface, uid int64) ([]Node, error) {

	fileds := []string{"id", "name", "type", "server", "method", "status", "startUserID", "endUserID"}
	query := string(" endUserID >= ? AND startUserID <= ?")
	queryArgs := []interface{}{uid, uid}
	selection := NewSelection(fileds, query, queryArgs)

	ctx := createContextWithValue(nodeTableName)

	var nodes []Node
	err := handle.GetToList(ctx, selection, &nodes)
	if err != nil {
		return nil, err
	}

	return nodes, err
}

func Getnodes(handle storage.Interface) (*Node, error) {
	host, err := network.ExternalIP()
	glog.Infof("Got host ip is %v\r\n", host)

	fileds := []string{"id", "name", "type", "server", "method", "status", "startUserID", "endUserID"}
	query := string("server = ?")
	queryArgs := []interface{}{host}
	selection := NewSelection(fileds, query, queryArgs)

	ctx := createContextWithValue(nodeTableName)

	var nodes []Node
	err = handle.GetToList(ctx, selection, &nodes)
	if nodes == nil {
		return nil, err
	}

	return &nodes[0], err
}
