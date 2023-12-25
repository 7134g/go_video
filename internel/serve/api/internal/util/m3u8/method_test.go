package m3u8

import (
	"path/filepath"
	"testing"
)

func TestMergeFiles(t *testing.T) {
	t.Log(filepath.Split("E:\\documents\\Go\\_programe\\go_video\\download\\test"))
	err := MergeFiles("E:\\documents\\Go\\_programe\\go_video\\download\\test")
	t.Log(err)
}
