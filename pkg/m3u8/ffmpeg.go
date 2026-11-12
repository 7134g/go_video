package m3u8

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

// MergeFilesFfmpeg 把 savePath 下的分片合并为上级目录的 <dirName>.mp4（纯音频为 .m4a）。
//
// 分片按 "%06d.<ext>" 命名，补零后的文件名字典序即播放顺序。
// 存在 init.mp4（fMP4）时生成本地 m3u8 交给 ffmpeg 的 HLS demuxer 合并，
// 否则使用 concat demuxer。
func MergeFilesFfmpeg(savePath, ffmpeg string) error {
	info, err := os.Stat(ffmpeg)
	if err != nil || info == nil {
		return errors.New("ffmpeg error")
	}

	segments, hasInit, audioOnly, err := collectSegments(savePath)
	if err != nil {
		return err
	}

	dirName := filepath.Base(savePath)
	ext := ".mp4"
	if audioOnly {
		ext = ".m4a"
	}
	outPath := filepath.Join(filepath.Dir(savePath), dirName+ext)

	var cmdList []string
	var tempFile string
	if hasInit {
		tempFile = filepath.Join(savePath, "local.m3u8")
		if err := writeLocalPlaylist(tempFile, segments); err != nil {
			return err
		}
		cmdList = []string{
			"-allowed_extensions", "ALL",
			"-i", tempFile,
			"-c", "copy",
			"-y", outPath,
		}
	} else {
		tempFile = filepath.Join(savePath, dirName+".txt")
		if err := writeConcatList(tempFile, segments); err != nil {
			return err
		}
		cmdList = []string{
			"-f", "concat",
			"-safe", "0",
			"-i", tempFile,
			"-acodec", "copy",
			"-vcodec", "copy",
		}
		if needAdtsToAsc(segments) {
			cmdList = append(cmdList, "-bsf:a", "aac_adtstoasc")
		}
		cmdList = append(cmdList, "-y", outPath)
	}
	defer os.Remove(tempFile)

	cmd := exec.Command(ffmpeg, cmdList...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg 合并失败: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// collectSegments 返回按播放顺序排列的分片文件名，以及是否存在 init 分片、是否纯音频。
func collectSegments(savePath string) (segments []string, hasInit bool, audioOnly bool, err error) {
	entries, err := os.ReadDir(savePath)
	if err != nil {
		return nil, false, false, err
	}

	audioOnly = true
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if name == InitSegmentName {
			hasInit = true
			continue
		}
		if !IsSegmentFile(name) {
			continue
		}
		segments = append(segments, name)
		ext := strings.ToLower(path.Ext(name))
		if ext != ".m4a" && ext != ".aac" {
			audioOnly = false
		}
	}

	if len(segments) == 0 {
		return nil, hasInit, false, errors.New("未找到分片文件，无法合并")
	}
	if !hasInit {
		for _, name := range segments {
			if strings.EqualFold(path.Ext(name), ".m4s") {
				return nil, hasInit, false, errors.New("分片为 fMP4 碎片但缺少 EXT-X-MAP 初始化分片，无法合并")
			}
		}
	}
	return segments, hasInit, audioOnly, nil
}

// writeLocalPlaylist 生成供 ffmpeg HLS demuxer 读取的本地播放列表。
// 时长为固定值即可（VOD + ENDLIST 不依赖它等待分片）。
func writeLocalPlaylist(listPath string, segments []string) error {
	var b strings.Builder
	b.WriteString("#EXTM3U\n")
	b.WriteString("#EXT-X-VERSION:7\n")
	b.WriteString("#EXT-X-TARGETDURATION:10\n")
	fmt.Fprintf(&b, "#EXT-X-MAP:URI=\"%s\"\n", InitSegmentName)
	for _, name := range segments {
		b.WriteString("#EXTINF:10,\n")
		b.WriteString(name + "\n")
	}
	b.WriteString("#EXT-X-ENDLIST\n")
	return os.WriteFile(listPath, []byte(b.String()), 0644)
}

// writeConcatList 生成 ffmpeg concat demuxer 使用的文件列表。
func writeConcatList(listPath string, segments []string) error {
	var b strings.Builder
	for _, name := range segments {
		fmt.Fprintf(&b, "file '%s'\n", name)
	}
	return os.WriteFile(listPath, []byte(b.String()), 0644)
}

// needAdtsToAsc 判断是否需要把 ADTS AAC（TS / .aac）转换成 MP4 的 ASC 格式。
func needAdtsToAsc(segments []string) bool {
	for _, name := range segments {
		ext := strings.ToLower(path.Ext(name))
		if ext == ".ts" || ext == ".aac" {
			return true
		}
	}
	return false
}
