package db

import (
	"fmt"
	"time"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/pagination"

	"golib/pkg/storage"

	"github.com/golang/glog"
)

func GetApiServers(handle storage.Interface, page pagination.Pager) ([]api.APIServerInfor, error) {

	ctx := createContextWithValue(apiServerTableName)
	fileds := []string{"id", "name", "host", "port", "status", "created_time"}
	selection, err := buildListSelecttion(ctx, handle, page, fileds)
	if err != nil {
		return nil, err
	}

	var servers []api.APIServerInfor
	err = handle.GetToList(ctx, selection, &servers)
	if err != nil {
		return nil, err
	}

	if len(servers) > 0 {
		return servers, nil
	} else {
		return nil, fmt.Errorf("not found")
	}

}

//func CreateAPIServer(handle storage.Interface, name string, host string, port int64, isEnable bool) error {
func CreateAPIServer(handle storage.Interface, info api.APIServerInfor) error {

	ctx := createContextWithValue(apiServerTableName)

	server := &api.APIServerInfor{
		Name:       info.Name,
		Host:       info.Host,
		Port:       info.Port,
		Status:     info.Status,
		CreateTime: time.Now(),
	}

	err := handle.Create(ctx, info.Name, server, server)
	if err != nil {
		glog.Errorf("create a server record failure %v\r\n", err)
	}
	return err
}

func DeleteAPIServerByName(handle storage.Interface, name string) error {

	ctx := createContextWithValue(apiServerTableName)

	var server api.APIServerInfor
	query := string("name = ?")
	queryArgs := []interface{}{name}
	selection := NewSelection(nil, query, queryArgs)

	err := handle.Delete(ctx, selection, &server)
	if err != nil {
		glog.Errorf("delete a server record failure %v\r\n", err)
	}
	return err
}
