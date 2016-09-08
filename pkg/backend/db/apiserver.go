package db

import (
	"fmt"
	"shadowsocks-go/pkg/storage"
	"time"

	"github.com/golang/glog"
)

//User is a mysql users map
type APIServers struct {
	ID         int64     `column:"id"`
	Name       string    `column:"name"`
	Host       string    `column:"host"`
	Port       int64     `column:"port"`
	Status     string    `column:"status"`
	CreateTime time.Time `column:"created_time" gorm:"column:created_time"`
}

func GetApiServers(handle storage.Interface) ([]APIServers, error) {

	fileds := []string{"id", "name", "host", "port", "status", "created_time"}
	query := string("status = ?")
	queryArgs := []interface{}{string("Enable")}
	selection := NewSelection(fileds, query, queryArgs)

	ctx := createContextWithValue(apiServerTableName)

	var servers []APIServers
	err := handle.GetToList(ctx, selection, &servers)
	if err != nil {
		return nil, err
	}

	if len(servers) > 0 {
		return servers, nil
	} else {
		return nil, fmt.Errorf("not found")
	}

}

func CreateAPIServer(handle storage.Interface, name string, host string, port int64, isEnable bool) error {

	ctx := createContextWithValue(apiServerTableName)

	var status string
	if isEnable {
		status = string("Enable")
	} else {
		status = string("Disable")
	}

	server := &APIServers{
		Name:       name,
		Host:       host,
		Port:       port,
		Status:     status,
		CreateTime: time.Now(),
	}

	err := handle.Create(ctx, host, server, server)
	if err != nil {
		glog.Errorf("create a server record failure %v\r\n", err)
	}
	return err
}

func DeleteAPIServerByID(handle storage.Interface, id int64) error {

	ctx := createContextWithValue(apiServerTableName)

	var server APIServers
	query := string("id = ?")
	queryArgs := []interface{}{id}
	selection := NewSelection(nil, query, queryArgs)

	err := handle.Delete(ctx, selection, &server)
	if err != nil {
		glog.Errorf("delete a server record failure %v\r\n", err)
	}
	return err
}
