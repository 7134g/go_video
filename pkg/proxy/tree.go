package proxy

import (
	"net/http"
)

type webTree struct {
	urlAddress string
	content    []byte
	List       []*webTree
}

var webMap = map[string]*webTree{}

func addNewRequest(req *http.Request) {
	// 将每一个http内容和请求内容存放在 webTree

}
