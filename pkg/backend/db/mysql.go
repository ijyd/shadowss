package db

import (
	"shadowsocks-go/pkg/storage"
	"shadowsocks-go/pkg/storage/storagebackend"
	"shadowsocks-go/pkg/storage/storagebackend/factory"

	"golang.org/x/net/context"
)

const (
	userTableName      = "user"
	nodeTableName      = "ss_node"
	userTokeTableName  = "user_token"
	apiServerTableName = "api_server"
	ctxKey             = "table"
)

func createContextWithValue(val string) context.Context {
	ctx := context.WithValue(context.TODO(), ctxKey, val)

	return ctx
}

func NewStorage(c storagebackend.Config) (storage.Interface, error) {
	return factory.Create(c)
}
