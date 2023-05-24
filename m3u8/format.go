package m3u8

import (
	"fmt"
	"os"
	"path/filepath"
)

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
