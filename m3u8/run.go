package m3u8

import (
	"bufio"
	"dv/base"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Execute m3u8 下载
func Execute(filename string, params []string) error {
	// m3u8.exe %url% --workDir E:\recreation\github --enableDelAfterDone --saveName %filename%
	cmd := exec.Command(params[0], params[1:]...)
	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	reader := bufio.NewReader(out)
	for true {
		line, _, readErr := reader.ReadLine()
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			log.Fatalln(err)
		}
		gbkData, err := simplifiedchinese.GBK.NewDecoder().Bytes(line)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("doing......", filename) // 进度信息
		fmt.Println(string(gbkData))
	}
	fmt.Println("done ========》", filename)

	return nil
}

type Dm3u8 struct {
	base.Downloader

	M3u8BaseLink    string // 分片所使用的基础链接地址
	fileCount       int    // 已下载的分片数量
	fileFutureCount int    // 需要下载的所有分片数

	Crypto       []byte      // 加密密匙
	CryptoMethod CryptMethod // 加密方式
}

func NewDownloader(taskName, saveDir, httpUrl string) Dm3u8 {
	m := Dm3u8{}
	m.TaskName = taskName
	m.SaveDir = filepath.Join(saveDir, taskName)
	m.Link = httpUrl
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
	res.Header = d.GetHeader()
	resp, err := d.GetClient().Do(res)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	// 处理请求
	info, err := os.Stat(d.SaveDir)
	if err != nil || !info.IsDir() {
		_ = os.MkdirAll(d.SaveDir, os.ModeDir)
	}

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
		d.CryptoMethod = key.Method
		resp, err := d.GetClient().Get(d.M3u8BaseLink + key.URI)
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
		d.Crypto = b
		break
	}

	return decode, nil
}

func (d Dm3u8) MergeFiles() error {
	outputFilepath := filepath.Join(d.SaveDir, "../", d.TaskName+".mp4")
	outputFile, err := os.Create(outputFilepath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	files, err := getFilesInDir(d.SaveDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		inputFilepath := filepath.Join(d.SaveDir, file)
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
