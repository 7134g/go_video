package m3u8

import (
	"bufio"
	"dv/base"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
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

	fileCount       int // 已下载的分片数量
	fileFutureCount int // 需要下载的所有分片数
}

func NewDownloader(taskName, saveDir, httpUrl string) Dm3u8 {
	m := Dm3u8{}
	m.TaskName = taskName
	m.SaveDir = filepath.Join(saveDir, taskName)
	m.Link = httpUrl
	return m
}

// ExtractSegments 提取m3u8内所有内容
func (d Dm3u8) ExtractSegments() ([]string, error) {
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
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	body := string(bodyBytes)
	body = strings.TrimPrefix(body, "\n")
	body = strings.TrimSuffix(body, "\n")
	lines := strings.Split(body, "\n")
	baseurl := strings.TrimSuffix(d.Link, d.Link[strings.LastIndex(d.Link, "/")+1:])

	// #EXTINF  开始
	// #EXT-X-ENDLIST 结束
	if !strings.Contains(body, "#EXTINF:") {
		// 第一个m3u8拿到的是m3u8列表, 选择最后一个m3u8的链接地址
		m3u8Url := lines[len(lines)-1]
		d.Link = baseurl + m3u8Url
		segment, err := d.ExtractSegments()
		if err != nil {
			return nil, err
		}
		return segment, err
	}

	// 拼接所有分片下载地址
	segment := make([]string, 0)
	for _, line := range lines {
		if !strings.HasPrefix(line, "#") && strings.TrimSpace(line) != "" {
			segment = append(segment, baseurl+line)
		}
	}

	return segment, err
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

//func (d Dm3u8) Execute() error {
//	// 构建请求
//	res, err := http.NewRequest(http.MethodGet, d.Link, nil)
//	if err != nil {
//		return err
//	}
//	res.Header = d.GetHeader()
//	resp, err := d.GetClient().Do(res)
//	if resp != nil {
//		defer resp.Body.Close()
//	}
//	if err != nil {
//		return err
//	}
//
//	// 处理请求
//	dirname := filepath.Join(d.SaveDir, d.TaskName)
//	info, err := os.Stat(dirname)
//	if err != nil || !info.IsDir() {
//		_ = os.MkdirAll(dirname, os.ModeDir)
//	}
//	bodyBytes, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return err
//	}
//	lines := strings.Split(string(bodyBytes), "\n")
//	for index, line := range lines {
//		if !strings.HasPrefix(line, "#") && strings.TrimSpace(line) != "" {
//			if err := d.downloadFile(dirname, line); err != nil {
//				return err
//			}
//			d.Doing(d.TaskName, fmt.Sprintf("分片 %d 下载成功", index))
//			d.Doing(d.TaskName, fmt.Sprintf("已下载的分片数(%d/%d)", d.fileCount, d.fileFutureCount))
//		}
//	}
//
//	return nil
//}
//
//func (d *Dm3u8) downloadFile(dirname, link string) error {
//	// 构建请求
//	res, err := http.NewRequest(http.MethodGet, link, nil)
//	if err != nil {
//		return err
//	}
//	res.Header = d.GetHeader()
//	resp, err := d.GetClient().Do(res)
//	if resp != nil {
//		defer resp.Body.Close()
//	}
//	if err != nil {
//		return err
//	}
//
//	filename := filepath.Base(link)
//	videoPath := filepath.Join(dirname, filename)
//	f, err := os.Create(videoPath)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//
//	if _, err := io.Copy(f, resp.Body); err != nil {
//		return err
//	}
//	return nil
//}
