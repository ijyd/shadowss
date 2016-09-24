package controller

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/backend"
	"cloud-keeper/pkg/controller/apiserverctl"
	"cloud-keeper/pkg/controller/nodectl"
	"cloud-keeper/pkg/etcdhelper"
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/storage"
	"gofreezer/pkg/watch"

	"github.com/golang/glog"
)

const (
	perNodeUserLimit = 30
)

var AutoSchedule *NodeSchedule

func ControllerStart(helper *etcdhelper.EtcdHelper, be *backend.Backend, host string, port int) error {
	AutoSchedule = NewNodeSchedule(helper, be)

	//add apiserver node
	has := apiserverctl.CheckLocalAPIServer(helper)
	glog.V(5).Infof("has local server %v", has)
	if !has {
		_, err := apiserverctl.AddLocalAPIServer(be, helper, host, port, uint64(0), true, true)
		if err != nil {
			return err
		}
	}

	go manageNode(helper)

	return nil
}

func AllocNode(user *api.User) error {
	return AutoSchedule.AllocNode(user)
}

func manageNode(helper *etcdhelper.EtcdHelper) {
	watchKey := nodectl.PrefixNode
	ctx := prototype.NewContext()
	resourceVer := string("")

	glog.V(5).Infof("watch at %v with resource %v", watchKey, resourceVer)
	watcher, err := helper.StorageCodec.Storage.WatchList(ctx, watchKey, resourceVer, storage.Everything)

	if err != nil {
		glog.Fatalf("Unexpected error: %v", err)
	}
	defer watcher.Stop()

	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				glog.Errorf("Unexpected channel close")
				return
			}

			glog.V(5).Infof("Got event  %#v", event.Type)
			switch event.Type {
			case watch.Added:
				glog.V(5).Infof("Got Add  got: %#v", event.Object)
				gotObject, ok := event.Object.(*api.Node)
				if ok {
					AutoSchedule.NewNode(gotObject)
				} else {
					gotObject, ok := event.Object.(*api.NodeUser)
					if ok {
						AutoSchedule.NewNodeUser(gotObject)
					}
				}
			case watch.Modified:
				glog.V(5).Infof("Got modify  got: %#v", event.Object)
				gotObject, ok := event.Object.(*api.Node)
				if ok {
					AutoSchedule.UpdateNode(gotObject)
				} else {
					gotObject, ok := event.Object.(*api.NodeUser)
					if ok {
						AutoSchedule.UpdateNodeUser(gotObject)
					}
				}
			case watch.Deleted:
				glog.V(5).Infof("Got Deleted  got: %#v", event.Object)
				gotObject, ok := event.Object.(*api.Node)
				if ok {
					AutoSchedule.DelNode(gotObject)
				} else {
					gotObject, ok := event.Object.(*api.NodeUser)
					if ok {
						AutoSchedule.DelNodeUser(gotObject)
					}
				}
			case watch.Error:
				glog.V(5).Infof("Got Error  got: %#v", event.Object)
				return
			default:
				glog.Errorf("UnExpected: %#v, got: %#v", event.Type, event.Object)
			}

		}
	}
}
