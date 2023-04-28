package base

import "net/http"

type Downloader struct {
	Logger

	TaskName string // 文件名
	SaveDir  string // 安装目录
	Link     string // http地址
	script   string // 类型

	header http.Header
	client *http.Client
}

func (d *Downloader) SetHeader(header http.Header) {
	d.header = header
}

func (d *Downloader) GetHeader() http.Header {
	return d.header
}

func (d *Downloader) SetClient(client *http.Client) {
	d.client = client
}

func (d *Downloader) GetClient() *http.Client {
	return d.client
}

func (d *Downloader) SetScript(script string) {
	d.script = script
}

func (d *Downloader) GetScript() string {
	return d.script
}
