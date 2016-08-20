package storage

//GetResult get storage  result by this handle
type GetResult func(x, y int) error

//Interface implement a storeage backend
type Interface interface {
	GetToList(table string, fields []string, result interface{}) error
}
