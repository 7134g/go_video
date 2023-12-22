package m3u8

import (
	"dv/config"
	"dv/table"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func ExtractContain(name, u string) (*M3u8, error) {
	// 构建请求
	res, err := http.NewRequest(http.MethodGet, u, nil)
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

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
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
		link, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		link, err = link.Parse(decode.MasterPlaylist[index].URI)
		if err != nil {
			return nil, err
		}
		return ExtractContain(name, link.String())
	}

	for _, key := range decode.Keys {
		if key.Method == CryptMethodNONE {
			continue
		}
		// 获取加密密匙
		link, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		aesLink, err := link.Parse(key.URI)
		if err != nil {
			return nil, err
		}
		resp, err := config.Client.Get(aesLink.String())
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			return nil, err
		}
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		table.CryptoVideoTable.Set(name, b)
		break
	}

	return decode, nil
}

const (
	hour   = 3600
	minute = 60
)

// CalculationTime 计算播放总时长
func CalculationTime(d float32) string {
	t := int(d)

	h := t / hour              // 计算小时数
	m := (t - h*hour) / minute // 计算分钟数
	s := t - h*hour - m*minute // 计算剩余的秒数

	return fmt.Sprintf("%d h %d m %d s", h, m, s)
}

// 获取目录下的文件列表
func getFilesInDir(dirname string) ([]string, error) {
	var files []string

	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, info.Name())
		}
		return nil
	})

	return files, err
}

func MergeFiles(saveDir string) error {
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
