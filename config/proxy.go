package config

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
)

var Client = http.DefaultClient

func httpProxy(proxy string) func(*http.Request) (*url.URL, error) {
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		log.Fatalln(err)
	}
	return http.ProxyURL(proxyUrl)
}

func GetHttpProxyClient(proxy string) *http.Client {
	return &http.Client{
		//Timeout: time.Second * 5,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			}, // 使用环境变量的代理
			Proxy: httpProxy(proxy),
		},
	}
}
