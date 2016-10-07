package common

import (
	"cloud-keeper/pkg/backend"
	"cloud-keeper/pkg/etcdhelper"
)

var Storage *backend.Backend
var EtcdStorage *etcdhelper.EtcdHelper
