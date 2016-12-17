package batchshadowss

import (
	"cloud-keeper/pkg/ansible"
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/account"
	"cloud-keeper/pkg/registry/core/account/accserver"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/runtime"

	"github.com/golang/glog"
)

type BatchShadowssREST struct {
	acc    account.Registry
	accsrv accserver.Registry
}

func NewREST(acc account.Registry, accsrv accserver.Registry) *BatchShadowssREST {
	return &BatchShadowssREST{
		acc:    acc,
		accsrv: accsrv,
	}
}

func (*BatchShadowssREST) New() runtime.Object {
	return &api.BatchShadowss{}
}

// func (r *BatchShadowssREST) upgradeShadowss(ctx freezerapi.Context, accName string) error {
// 	accsrvListOptions := &freezerapi.ListOptions{}
// 	accsrvListOptions.PageSelector = pages.SelectorFromSet(1, 10)
// 	accsrvListOptions.FieldSelector = fields.SelectorFromSet(map[string]string{
// 		"metadata.name": accName,
// 	})
//
// 	pagecnt := 1
// 	var hosts []string
// 	for {
// 		accsrvlist, err := r.accsrv.ListAccServers(ctx, accsrvListOptions)
// 		if err != nil {
// 			glog.Warningf("get account(%v) servers list error %v\r\n", accName, err)
// 			return err
// 		}
//
// 		for _, v := range accsrvlist.Items {
// 			if len(v.Spec.Vultr.IPV4Addr) > 0 {
// 				hosts = append(hosts, v.Spec.Vultr.IPV4Addr)
// 			} else if len(v.Spec.DigitalOcean.IPV4Addr) > 0 {
// 				hosts = append(hosts, v.Spec.DigitalOcean.IPV4Addr)
// 			} else {
// 				glog.Warningf("not found any valid server ip(%v) \r\n", v.Name)
// 				continue
// 			}
// 		}
//
// 		has, page, _ := accsrvListOptions.PageSelector.LastPage()
// 		if has {
// 			if page > uint64(pagecnt) {
// 				accsrvListOptions.PageSelector = pages.SelectorFromSet(uint64(pagecnt)+1, 10)
// 			} else {
// 				break
// 			}
// 		} else {
// 			break
// 		}
// 		pagecnt++
//
// 	}
//
// 	if len(hosts) > 0 {
// 		glog.Infof("upgrade shadowss host list(%v)", hosts)
// 		go ansible.UpgradeShadowss(hosts)
// 	}
//
// 	return nil
// }

func (r *BatchShadowssREST) upgradeShadowss(ctx freezerapi.Context, target []api.TargetAccServer) error {

	var hosts []string
	for _, v := range target {
		hosts = append(hosts, v.Host)
	}

	if len(hosts) > 0 {
		glog.Infof("upgrade shadowss host list(%v)", hosts)
		go ansible.UpgradeShadowss(hosts)
	}

	return nil
}

func (r *BatchShadowssREST) Create(ctx freezerapi.Context, obj runtime.Object) (runtime.Object, error) {

	shadowss := obj.(*api.BatchShadowss)

	r.upgradeShadowss(ctx, shadowss.Spec.Target)

	return shadowss, nil
}
