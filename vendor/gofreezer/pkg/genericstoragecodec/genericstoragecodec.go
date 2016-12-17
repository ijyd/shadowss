package genericstoragecodec

import (
	"fmt"
	"mime"

	"gofreezer/pkg/api"
	"gofreezer/pkg/genericstoragecodec/options"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/runtime/schema"
	"gofreezer/pkg/runtime/serializer/recognizer"
	"gofreezer/pkg/storage/etcds"
	"gofreezer/pkg/storage/storagebackend"
	"gofreezer/pkg/storage/storagebackend/factory"

	"github.com/golang/glog"
)

var memoryVersion = api.SchemeGroupVersion

type GenericStorageCodec struct {
	Storage        etcds.Interface
	Codecs         runtime.Codec
	DestroyStorage factory.DestroyFunc
}

func NewGenericStorageCodec(options *options.StorageOptions, ns runtime.StorageSerializer, storageVersion schema.GroupVersion) (*GenericStorageCodec, error) {
	storageConfig := options.StorageConfig
	codec, err := newStorageCodec(options.DefaultStorageMediaType, ns, storageVersion, memoryVersion, storageConfig)
	if err != nil {
		return nil, err
	}

	storageConfig.Codec = codec
	storageConfig.StorageVersion = storageVersion.Version

	storageHandle, destroy, err := factory.Create(storageConfig)
	if err != nil {
		return nil, err
	}

	storageEtcd := storageHandle.(etcds.Interface)

	return &GenericStorageCodec{
		Storage:        storageEtcd,
		Codecs:         codec,
		DestroyStorage: destroy,
	}, nil

}

// newStorageCodec assembles a storage codec for the provided storage media type, the provided serializer, and the requested
// storage and memory versions.
func newStorageCodec(storageMediaType string, ns runtime.StorageSerializer, storageVersion, memoryVersion schema.GroupVersion, config storagebackend.Config) (runtime.Codec, error) {
	mediaType, _, err := mime.ParseMediaType(storageMediaType)
	if err != nil {
		return nil, fmt.Errorf("%q is not a valid mime-type", storageMediaType)
	}
	serializer, ok := runtime.SerializerInfoForMediaType(ns.SupportedMediaTypes(), mediaType)
	if !ok {
		return nil, fmt.Errorf("unable to find serializer for %q", storageMediaType)
	}

	s := serializer.Serializer

	// etcd2 only supports string data - we must wrap any result before returning
	// TODO: storagebackend should return a boolean indicating whether it supports binary data
	if !serializer.EncodesAsText && (config.Type == storagebackend.StorageTypeUnset || config.Type == storagebackend.StorageTypeETCD2) {
		glog.V(4).Infof("Wrapping the underlying binary storage serializer with a base64 encoding for etcd2")
		s = runtime.NewBase64Serializer(s)
	}

	encoder := ns.EncoderForVersion(
		s,
		runtime.NewMultiGroupVersioner(
			storageVersion,
			schema.GroupKind{Group: storageVersion.Group},
			schema.GroupKind{Group: memoryVersion.Group},
		),
	)

	ds := recognizer.NewDecoder(s, ns.UniversalDeserializer())
	decoder := ns.DecoderToVersion(
		ds,
		runtime.NewMultiGroupVersioner(
			memoryVersion,
			schema.GroupKind{Group: memoryVersion.Group},
			schema.GroupKind{Group: storageVersion.Group},
		),
	)

	return runtime.NewCodec(encoder, decoder), nil
}
