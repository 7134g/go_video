package m3u8

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func MergeFilesFfmpeg(savePath, taskName, ffmpeg string) error {
	info, err := os.Stat(ffmpeg)
	if err != nil || info == nil {
		return errors.New("ffmpeg error")
	}

	listName := fmt.Sprintf("%s.txt", taskName)
	listFilePath, err := os.Create(listName)
	if err != nil {
		return err
	}
	_ = filepath.Walk(savePath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		line := fmt.Sprintf("file '%s/%s'\n", savePath, info.Name())
		_, _ = listFilePath.Write([]byte(line))
		return nil
	})
	_ = listFilePath.Close()
	defer func() {
		_ = os.Remove(listName)
	}()

	// ffmpeg -f concat -safe 0 -i list.txt  -acodec copy -vcodec copy -absf aac_adtstoasc -y output.mp4
	outPath := filepath.Join(savePath, "../", fmt.Sprintf("%s.mp4", taskName))
	cmd := exec.Command(
		ffmpeg,
		"-f",
		"concat",
		"-safe",
		"0",
		"-i",
		listName,
		"-acodec",
		"copy",
		"-vcodec",
		"copy",
		"-absf",
		"aac_adtstoasc",
		"-y",
		outPath,
	)

	return cmd.Run()
}
