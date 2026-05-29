package m3u8

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

// MergeFilesFfmpeg 使用 ffmpeg concat demuxer 把 savePath 下所有分段合并为 <dirName>.mp4。
//
// 隐式假设：分段文件名按字典序即为播放顺序。controller 写入时用 "%06d.ts" 命名，
// 与 filepath.Walk 的字典序遍历刚好一致；若调用方改了命名规则需重新审视。
func MergeFilesFfmpeg(savePath, ffmpeg string) error {
	info, err := os.Stat(ffmpeg)
	if err != nil || info == nil {
		return errors.New("ffmpeg error")
	}

	_, dirName := filepath.Split(savePath)
	//h := md5.New()
	//h.Write([]byte(dirName))
	//name := hex.EncodeToString(h.Sum(nil))
	listNamePath := filepath.Join(savePath, fmt.Sprintf("%s.txt", dirName))
	listFilePath, err := os.Create(listNamePath)
	if err != nil {
		return err
	}
	_ = filepath.Walk(savePath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		line := fmt.Sprintf("file '%s'\n", info.Name())
		_, _ = listFilePath.Write([]byte(line))
		return nil
	})
	_ = listFilePath.Close()
	//defer os.Remove(listNamePath)

	// ffmpeg -f concat -safe 0 -i list.txt  -acodec copy -vcodec copy -absf aac_adtstoasc -y output.mp4
	outPath := filepath.Join(savePath, "../", fmt.Sprintf("%s.mp4", dirName))
	cmdList := []string{
		"-f",
		"concat",
		"-safe",
		"0",
		"-i",
		listNamePath,
		"-acodec",
		"copy",
		"-vcodec",
		"copy",
		"-absf",
		"aac_adtstoasc",
		"-y",
		outPath,
	}
	cmd := exec.Command(
		ffmpeg,
		cmdList...,
	)

	cmdString := ffmpeg
	for _, s := range cmdList {
		cmdString = fmt.Sprintf("%s %s", cmdString, s)
	}
	fmt.Println(cmdString)

	return cmd.Run()
}
