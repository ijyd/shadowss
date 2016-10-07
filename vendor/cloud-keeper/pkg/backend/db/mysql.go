package db

import (
	"golib/pkg/storage"
	"golib/pkg/storage/storagebackend"
	"golib/pkg/storage/storagebackend/factory"

	"golang.org/x/net/context"
)

const (
	userTableName      = "user"
	apiKeyTableName    = "vps_server_account"
	nodeTableName      = "ss_node"
	userTokeTableName  = "user_token"
	apiServerTableName = "api_server"

	ctxKey = "table"
)

func createContextWithValue(val string) context.Context {
	ctx := context.WithValue(context.TODO(), ctxKey, val)

	return ctx
}

func NewStorage(c storagebackend.Config) (storage.Interface, error) {
	return factory.Create(c)
}
