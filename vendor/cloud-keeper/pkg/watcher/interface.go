package watcher

import "gofreezer/pkg/runtime"

//Interface implements your watch event
type Interface interface {
	AddObj(runtime.Object)
	ModifyObj(runtime.Object)
	DelObj(runtime.Object)
	Error(runtime.Object)
}
