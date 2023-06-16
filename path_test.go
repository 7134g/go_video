package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestPath(t *testing.T) {
	path := `E:\documents\Go\_programe\go_video\download\最差劲`
	list, _ := os.Create("list.txt")
	defer list.Close()

	//files, err := ioutil.ReadDir(path)
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//for _, f := range files {
	//	_, _ = list.Write([]byte(fmt.Sprintf("file '%s'\n", f.Name())))
	//}

	_ = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		//fmt.Println(fmt.Sprintf("file '%s'\n", info.Name()))
		_, _ = list.Write([]byte(fmt.Sprintf("file './最差劲/%s'\n", info.Name())))
		return nil
	})
}
