package users

import (
	"shadowsocks-go/pkg/storage"
	"shadowsocks-go/pkg/storage/storagebackend"
	"shadowsocks-go/pkg/storage/storagebackend/factory"

	"golang.org/x/net/context"
)

const (
	userTableName = "user"
	nodeTableName = "ss_node"
	ctxKey        = "table"
)

func createContextWithValue(val string) context.Context {
	ctx := context.WithValue(context.TODO(), ctxKey, val)

	return ctx
}

func newStorage(c storagebackend.Config) (storage.Interface, error) {
	return factory.Create(c)
}
