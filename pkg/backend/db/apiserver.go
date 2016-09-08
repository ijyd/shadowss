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
	Host       string    `column:"host"`
	Port       int64     `column:"port"`
	Status     string    `colume:"status"`
	CreateTime time.Time `column:"create_time" gorm:"column:created_time"`
}

func GetApiServers(handle storage.Interface) ([]APIServers, error) {

	fileds := []string{"id", "host", "port", "status", "create_time"}
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

func CreateAPIServer(handle storage.Interface, host string, port int64, isEnable bool) error {

	ctx := createContextWithValue(apiServerTableName)

	var status string
	if isEnable {
		status = string("Enable")
	} else {
		status = string("Disable")
	}

	server := &APIServers{
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
