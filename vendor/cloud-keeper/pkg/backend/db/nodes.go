package db

import (
	"fmt"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/pagination"

	"golib/pkg/storage"
	"golib/pkg/util/network"

	"github.com/golang/glog"
)

var nodeFileds = []string{"id", "name", "enableota",
	"server", "method", "status", "traffic_rate", "description",
	"traffic_limit", "upload", "download", "location", "vps_server_id", "vps_server_name", "total_upload", "total_download"}

func GetNodesByUserID(handle storage.Interface, uid int64) ([]api.NodeServer, error) {

	query := string(" endUserID >= ? AND startUserID <= ?")
	queryArgs := []interface{}{uid, uid}
	selection := NewSelection(nodeFileds, query, queryArgs)

	ctx := createContextWithValue(nodeTableName)

	var nodes []api.NodeServer
	err := handle.GetToList(ctx, selection, &nodes)
	if err != nil {
		return nil, err
	}

	if len(nodes) > 0 {
		return nodes, nil
	} else {
		return nil, fmt.Errorf("not found")
	}

}

func GetNodeByName(handle storage.Interface, name string) (*api.NodeServer, error) {
	query := string("name = ?")
	queryArgs := []interface{}{name}
	selection := NewSelection(nodeFileds, query, queryArgs)

	ctx := createContextWithValue(nodeTableName)

	var nodes []api.NodeServer
	err := handle.GetToList(ctx, selection, &nodes)
	if err != nil {
		return nil, err
	}

	if len(nodes) > 0 {
		return &nodes[0], nil
	} else {
		return nil, fmt.Errorf("not found")
	}
}

func Getnodes(handle storage.Interface) (*api.NodeServer, error) {
	host, err := network.ExternalIP()
	glog.Infof("Got host ip is %v\r\n", host)

	query := string("server = ?")
	queryArgs := []interface{}{host}
	selection := NewSelection(nodeFileds, query, queryArgs)

	ctx := createContextWithValue(nodeTableName)

	var nodes []api.NodeServer
	err = handle.GetToList(ctx, selection, &nodes)
	if nodes == nil {
		return nil, err
	}

	return &nodes[0], err
}

func GetNodes(handle storage.Interface, page pagination.Pager) ([]api.NodeServer, error) {

	// query := string("server = ?")
	// queryArgs := []interface{}{host}
	ctx := createContextWithValue(nodeTableName)

	selection, err := buildListSelecttion(ctx, handle, page, nodeFileds)
	if err != nil {
		return nil, err
	}

	var nodes []api.NodeServer
	err = handle.GetToList(ctx, selection, &nodes)
	if nodes == nil {
		return nil, err
	}

	return nodes, err
}

func CreateNode(handle storage.Interface, detail api.NodeServer) error {

	ctx := createContextWithValue(nodeTableName)

	err := handle.Create(ctx, detail.Name, &detail, &detail)
	if err != nil {
		glog.Errorf("create a node record failure %v\r\n", err)
	}
	return err
}

func DeleteNode(handle storage.Interface, name string) error {

	ctx := createContextWithValue(nodeTableName)

	var node api.NodeServer
	query := string("name = ?")
	queryArgs := []interface{}{name}
	selection := NewSelection(nil, query, queryArgs)

	err := handle.Delete(ctx, selection, &node)
	if err != nil {
		glog.Errorf("delete a apikey record failure %v\r\n", err)
	}
	return err
}

func UpdateNodeTraffic(handle storage.Interface, userID int64, totalUpload, totalDownload, upload, download int64) error {

	node := &api.NodeServer{
		ID:                   userID,
		Upload:               upload,
		Download:             download,
		TotalUploadTraffic:   totalUpload,
		TotalDownloadTraffic: totalDownload,
	}

	conditionFields := string("id")
	updateFields := []string{"upload", "download", "total_upload", "total_download"}

	ctx := createContextWithValue(nodeTableName)
	err := handle.GuaranteedUpdate(ctx, conditionFields, updateFields, node)
	return err
}

func UpdateNodeStatus(handle storage.Interface, nodeID int64, status bool) error {

	var statusInt int64
	if status {
		statusInt = 1
	} else {
		statusInt = 0
	}

	node := &api.NodeServer{
		ID:     nodeID,
		Status: statusInt,
	}

	conditionFields := string("id")
	updateFields := []string{"status"}

	glog.V(5).Infof("update node status %v\r\n", node)

	ctx := createContextWithValue(nodeTableName)
	err := handle.GuaranteedUpdate(ctx, conditionFields, updateFields, node)
	return err
}

func UpdateNode(handle storage.Interface, detail api.NodeServer) error {

	conditionFields := string("id")
	updateFields := nodeFileds

	ctx := createContextWithValue(nodeTableName)
	err := handle.GuaranteedUpdate(ctx, conditionFields, updateFields, detail)
	return err
}
