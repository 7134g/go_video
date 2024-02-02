package files

import (
	"os"
	"path/filepath"
)

// MakeDir 创建目录
func MakeDir(_dir string) bool {
	if IsExist(_dir) {
		return true
	}
	err := os.MkdirAll(_dir, os.ModePerm)
	if err != nil {
		return false
	} else {
		return true
	}
}

// IsExist 判断文件或目录是否存在
func IsExist(f string) bool {
	info, err := os.Stat(f)
	return err == nil || info != nil
}

func GetFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	MakeDir(dir)

	fl, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return fl, nil
}
