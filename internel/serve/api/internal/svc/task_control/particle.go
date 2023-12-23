package task_control

import (
	"bytes"
	"dv/base"
	"dv/internel/serve/api/internal/svc/task_control/m3u8"
	"dv/internel/serve/api/internal/table"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type download struct {
	key string // 任务标识

	fileDir  string // 目录
	fileName string // 文件名
	fileSize int64  // 目前文件大小
}

func buildKey(id uint, name string) string {
	return fmt.Sprintf("%d_%s", id, name)
}

func (d *download) getM3u8File(client *http.Client, _url string, header http.Header) ([]*m3u8.Segment, error) {
	// 构建请求
	buf := bytes.NewBuffer(nil)
	if err := d.get(client, _url, header, buf); err != nil {
		return nil, err
	}

	m3u8Data, err := m3u8.ParseM3u8Data(buf)
	if err != nil {
		return nil, err
	}

	if len(m3u8Data.MasterPlaylist) != 0 {
		// 下载最高清的视频
		index := m3u8Data.GetMaxBandWidth()
		if index < 0 {
			return nil, errors.New("解析失败")
		}
		link, err := url.Parse(_url)
		if err != nil {
			return nil, err
		}
		link, err = link.Parse(m3u8Data.MasterPlaylist[index].URI)
		if err != nil {
			return nil, err
		}
		return d.getM3u8File(client, link.String(), header)
	}

	for _, key := range m3u8Data.Keys {
		if key.Method == m3u8.CryptMethodNONE {
			continue
		}
		// 获取加密密匙
		link, err := url.Parse(_url)
		if err != nil {
			return nil, err
		}
		aesUrl, err := link.Parse(key.URI)
		if err != nil {
			return nil, err
		}
		aesBuf := bytes.NewBuffer(nil)
		if err := d.get(client, aesUrl.String(), header, aesBuf); err != nil {
			return nil, err
		}

		table.CryptoVideoTable.Set(d.key, aesBuf.Bytes())
		break
	}

	return m3u8Data.Segments, nil
}

func (d *download) get(client *http.Client, _url string, header http.Header, write io.Writer) error {
	// 构建请求
	req, err := http.NewRequest(http.MethodGet, _url, nil)
	if err != nil {
		return err
	}
	req.Header = header
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("%d %s", resp.StatusCode, resp.Status))
	}

	if resp.ContentLength <= 0 {
		return errors.New(fmt.Sprintf("%s 跳过数据内容大于等于文件大小，因此不下载\n", d.fileName))
	}

	return d.rw(resp.Body, write)
}

func (d *download) rw(read io.Reader, write io.Writer) error {
	bs := make([]byte, 1048576) // 每次读取http内容的大小(1mb)

	for {
		rn, err := read.Read(bs)
		if err != nil {
			if err == io.EOF {
				// 完成
				d.fileSize += int64(rn)
				_, _ = write.Write(bs[:rn])
				return nil
			}
			return err
		}

		d.fileSize += int64(rn)
		_, err = write.Write(bs[:rn])
		if err != nil {
			return err
		}
	}
}

func (d *download) decode(data []byte) []byte {
	if key, ok := table.CryptoVideoTable.Get(d.key); ok {
		return base.AESDecrypt(data, key)
	}

	return data
}
