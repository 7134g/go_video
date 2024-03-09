package h5

import (
	_ "embed"
	"net/http"
)

var (
	//go:embed dist/index.html
	html []byte
	//go:embed dist/assets/index-03d09010.css
	css []byte
	//go:embed dist/assets/index-137bce34.js
	js []byte
)

func Css(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	_, _ = w.Write(css)
}

func Js(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	_, _ = w.Write(js)
}

func Html(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(html)
}
