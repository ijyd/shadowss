package account

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/collector"
	"cloud-keeper/pkg/collector/collectorbackend"
	"cloud-keeper/pkg/collector/collectorbackend/factory"

	"github.com/golang/glog"

	"fmt"
	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/errors"
	"gofreezer/pkg/api/rest"
)

// Registry is an interface for things that know how to store node.
type Registry interface {
	GetAccount(ctx freezerapi.Context, name string) (*api.Account, error)
	ListAccounts(ctx freezerapi.Context, options *freezerapi.ListOptions) (*api.AccountList, error)
	CreateAccount(ctx freezerapi.Context, acc *api.Account) (*api.Account, error)
	DeleteAccount(ctx freezerapi.Context, name string) error
	//CreateCloudProvider(ctx freezerapi.Context) error
	GetCloudProvider(ctx freezerapi.Context, name string) (collector.Collector, error)
}

// storage puts strong typing around storage calls
type storage struct {
	rest.Creater
	rest.Updater
	rest.Deleter
	rest.Getter
	rest.Lister
	collectors map[string]collector.Collector
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(c rest.Creater, u rest.Updater, d rest.Deleter, g rest.Getter, l rest.Lister) Registry {
	return &storage{
		Creater:    c,
		Updater:    u,
		Deleter:    d,
		Getter:     g,
		Lister:     l,
		collectors: make(map[string]collector.Collector),
	}
}

func (s *storage) GetAccount(ctx freezerapi.Context, name string) (*api.Account, error) {
	obj, err := s.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return obj.(*api.Account), nil
}

func (s *storage) ListAccounts(ctx freezerapi.Context, options *freezerapi.ListOptions) (*api.AccountList, error) {
	obj, err := s.List(ctx, options)
	if err != nil {
		return nil, err
	}
	return obj.(*api.AccountList), nil
}

func (s *storage) CreateAccount(ctx freezerapi.Context, acc *api.Account) (*api.Account, error) {
	obj, err := s.Create(ctx, acc)
	if err != nil {
		return nil, err
	}
	return obj.(*api.Account), nil
}

func (s *storage) DeleteAccount(ctx freezerapi.Context, name string) error {
	_, err := s.Delete(ctx, name)
	return err
}

func (s *storage) GetCloudProvider(ctx freezerapi.Context, name string) (collector.Collector, error) {

	glog.Infof("Get cloud provider %s\r\n", name)
	collectorHandler, ok := s.collectors[name]
	if ok {
		return collectorHandler, nil
	} else {
		acc, err := s.GetAccount(ctx, name)
		if err != nil {
			return nil, err
		}

		cfg := collectorbackend.Config{
			Type:   string(acc.Spec.AccDetail.Operators),
			APIKey: acc.Spec.AccDetail.Key,
		}
		collectorHandle, err := factory.Create(cfg)
		if err == nil {
			s.collectors[name] = collectorHandle
			return collectorHandle, nil
		}
	}
	return nil, errors.NewBadRequest(fmt.Sprintf("not found cloud provider by account(%v)", name))
}
