package common

import "net/http"

const (
	Token = "Bearer 49acafe7e63682e1e6b6983580c4ee56"
)

func AddAuthHttpHeader(req *http.Request) {
	req.Header.Add("Authorization", Token)
}
