package m3u8

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMergeFilesNative_NoSegments(t *testing.T) {
	dir := t.TempDir()
	if err := MergeFilesNative(dir); err == nil {
		t.Error("expected error for directory with no .ts segments")
	}
}

func TestMergeFilesNative_IgnoresNonTs(t *testing.T) {
	dir := t.TempDir()
	// 仅有非 .ts 文件（例如 ffmpeg 版遗留的 list.txt），应视为无分片。
	if err := os.WriteFile(filepath.Join(dir, "list.txt"), []byte("file '000000.ts'\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := MergeFilesNative(dir); err == nil {
		t.Error("expected error when only non-ts files present")
	}
}

func TestMergeFilesNative_InvalidTs(t *testing.T) {
	dir := t.TempDir()
	// 非法 TS 内容：解封装不出任何音视频流，应报错且不残留 mp4。
	if err := os.WriteFile(filepath.Join(dir, "000000.ts"), []byte("not a real mpeg-ts stream"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := MergeFilesNative(dir); err == nil {
		t.Error("expected error for invalid ts content")
	}
	if _, err := os.Stat(filepath.Join(dir, "..", filepath.Base(dir)+".mp4")); err == nil {
		t.Error("expected no leftover mp4 on failure")
	}
}
