package web

type User struct {
	id             int
	port           int
	passwd         string
	method         string
	enable         int
	transferEnable int `json:"transfer_enable"`
	u              int
	d              int
}
