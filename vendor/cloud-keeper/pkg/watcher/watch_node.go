package watcher

// import (
// 	"cloud-keeper/pkg/etcdhelper"
// 	prototype "gofreezer/pkg/api"
// 	"gofreezer/pkg/storage"
// 	"gofreezer/pkg/watch"
//
// 	"github.com/golang/glog"
// )
//
// //WatchNodeUsersLoop create a routine for sync users from backend storage
// func WatchNodeUsersLoop(key string, helper *etcdhelper.EtcdHelper, callback Interface) {
// 	ctx := prototype.NewContext()
// 	resourceVer := string("0")
//
// 	watcher, err := helper.StorageCodec.Storage.WatchList(ctx, key, resourceVer, storage.Everything)
//
// 	if err != nil {
// 		glog.Fatalf("Unexpected error: %v", err)
// 	}
// 	defer watcher.Stop()
//
// 	for {
// 		select {
// 		case event, ok := <-watcher.ResultChan():
// 			if !ok {
// 				glog.Errorf("Unexpected channel close")
// 				return
// 			}
//
// 			glog.V(5).Infof("Got event  %#v", event.Type)
// 			switch event.Type {
// 			case watch.Added:
// 				gotObject := event.Object.(*api.NodeUser)
// 				callback.AddObj(gotObject)
// 			case watch.Modified:
// 				gotObject := event.Object.(*api.NodeUser)
// 				callback.ModifyObj(gotObject)
// 			case watch.Deleted:
// 				gotObject := event.Object.(*api.NodeUser)
// 				callback.DelObj(gotObject)
// 			case watch.Error:
// 				callback.Error(event.Object)
// 				return
// 			default:
// 				glog.Errorf("UnExpected: %#v, got: %#v", event.Type, event.Object)
// 			}
//
// 		}
// 	}
//
// }
