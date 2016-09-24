package watcher

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/etcdhelper"
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/storage"
	"gofreezer/pkg/watch"

	"github.com/golang/glog"
)

//WatchNodeUsersLoop create a routine for sync users from backend storage
func WatchNodeUsersLoop(key string, helper *etcdhelper.EtcdHelper, callback Interface) {
	//watchKey := "/" + "NodeConfig" + "/test"
	ctx := prototype.NewContext()
	resourceVer := string("")

	glog.V(5).Infof("watch at %v with resource %v", key, resourceVer)
	//watcher, err := helper.StorageCodec.Storage.Watch(ctx, key, resourceVer, storage.Everything)
	watcher, err := helper.StorageCodec.Storage.WatchList(ctx, key, resourceVer, storage.Everything)

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
				gotObject := event.Object.(*api.NodeUser)
				callback.AddObj(gotObject)
			case watch.Modified:
				glog.V(5).Infof("Got modify  got: %#v", event.Object)
				gotObject := event.Object.(*api.NodeUser)
				glog.Errorf("*******Got modify event not accept here %+v******", gotObject)
				//callback.ModifyObj(gotObject)
			case watch.Deleted:
				glog.V(5).Infof("Got Deleted  got: %#v", event.Object)
				gotObject := event.Object.(*api.NodeUser)
				callback.DelObj(gotObject)
			case watch.Error:
				glog.Errorf("Got Error  got: %#v", event.Object)
				return
			default:
				//gotObject := event.Object.(*api.NodeConfig)
				glog.Errorf("UnExpected: %#v, got: %#v", event.Type, event.Object)
			}

		}
	}

}
