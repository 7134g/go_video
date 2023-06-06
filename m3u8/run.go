package m3u8

import (
	"dv/config"
	"dv/table"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Dm3u8 struct {
	Name string
	Link string

	M3u8BaseLink string // 分片所使用的基础链接地址
}

func NewDownloader(name, link string) Dm3u8 {
	m := Dm3u8{}
	m.Name = name
	m.Link = link
	return m
}

// ExtractContain 提取m3u8内所有内容
func (d *Dm3u8) ExtractContain() (*M3u8, error) {
	// 设置当前基础链接
	d.M3u8BaseLink = strings.TrimSuffix(d.Link, d.Link[strings.LastIndex(d.Link, "/")+1:])

	// 构建请求
	res, err := http.NewRequest(http.MethodGet, d.Link, nil)
	if err != nil {
		return nil, err
	}
	res.Header = config.Header
	resp, err := config.Client.Do(res)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	// 处理请求
	decode, err := parse(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(decode.MasterPlaylist) != 0 {
		// 下载最高清的视频
		index := decode.GetMaxBandWidth()
		if index < 0 {
			return nil, errors.New("解析失败")
		}
		d.Link = d.M3u8BaseLink + decode.MasterPlaylist[index].URI
		return d.ExtractContain()
	}

	for _, key := range decode.Keys {
		if key.Method == CryptMethodNONE {
			continue
		}
		// 获取加密密匙
		resp, err := config.Client.Get(d.M3u8BaseLink + key.URI)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		table.CryptoVideoTable.Set(d.Name, b)
		break
	}

	return decode, nil
}

func (d Dm3u8) MergeFiles(saveDir string) error {
	outputFilepath := filepath.Join(saveDir, "../", d.Name+".mp4")
	outputFile, err := os.Create(outputFilepath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	files, err := getFilesInDir(saveDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		inputFilepath := filepath.Join(saveDir, file)
		inputFile, err := os.Open(inputFilepath)
		if err != nil {
			return err
		}
		defer inputFile.Close()

		if _, err := io.Copy(outputFile, inputFile); err != nil {
			return err
		}
	}

	return nil
}
