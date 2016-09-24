package etcdhelper

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/api/v1"

	"github.com/golang/glog"

	"gofreezer/pkg/genericstoragecodec"
	storageoptions "gofreezer/pkg/genericstoragecodec/options"
)

//EtcdHelper wrap etcd storage and codec
type EtcdHelper struct {
	StorageCodec *genericstoragecodec.GenericStorageCodec
}

//NewEtcdHelper create a EtcdHelper
func NewEtcdHelper(options *storageoptions.StorageOptions) *EtcdHelper {

	storageVersion := v1.SchemeGroupVersion

	storeCodec, err := genericstoragecodec.NewGenericStorageCodec(options, api.Codecs, storageVersion)
	if err != nil {
		glog.Errorf("new genericstoragecodec error %v\r\n", err)
		return nil
	}

	return &EtcdHelper{
		StorageCodec: storeCodec,
	}

}
