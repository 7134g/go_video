package m3u8

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/yapingcat/gomedia/go-mp4"
	"github.com/yapingcat/gomedia/go-mpeg2"
)

// lazyFileReader 按顺序把多个文件当成一个连续流读取，任意时刻只持有一个打开的 fd，
// 避免一次性 open 上千个分片导致文件描述符耗尽。
type lazyFileReader struct {
	paths []string
	idx   int
	cur   *os.File
}

func (r *lazyFileReader) Read(p []byte) (int, error) {
	for {
		if r.cur == nil {
			if r.idx >= len(r.paths) {
				return 0, io.EOF
			}
			f, err := os.Open(r.paths[r.idx])
			if err != nil {
				return 0, err
			}
			r.cur = f
			r.idx++
		}
		n, err := r.cur.Read(p)
		if errors.Is(err, io.EOF) {
			_ = r.cur.Close()
			r.cur = nil
			if n > 0 {
				return n, nil
			}
			continue // 当前文件读完，接着读下一个
		}
		return n, err
	}
}

func (r *lazyFileReader) Close() error {
	if r.cur != nil {
		err := r.cur.Close()
		r.cur = nil
		return err
	}
	return nil
}

// MergeFilesNative 纯 Go 实现：把 savePath 下所有 .ts 分片 remux 成 <dirName>.mp4，
// 不依赖 ffmpeg，不转码，仅做 MPEG-TS → MP4 的容器转换。
//
// 与 MergeFilesFfmpeg 一致地假设分段文件名按字典序即为播放顺序（controller 用 "%06d.ts" 命名）。
// 支持 H264/H265 视频与 AAC/MP3 音频；遇到库不支持的编码时返回错误，由调用方回退到 ffmpeg。
func MergeFilesNative(savePath string) (err error) {
	segments, err := filepath.Glob(filepath.Join(savePath, "*.ts"))
	if err != nil {
		return err
	}
	if len(segments) == 0 {
		return errors.New("no ts segments found")
	}
	sort.Strings(segments)

	dirName := filepath.Base(filepath.Clean(savePath))
	outPath := filepath.Join(savePath, "../", fmt.Sprintf("%s.mp4", dirName))
	out, err := os.OpenFile(outPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	// 统一收尾：关闭文件；任一步出错则删除半成品 mp4，避免回退 ffmpeg 时残留。
	defer func() {
		if cerr := out.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			_ = os.Remove(outPath)
		}
	}()

	muxer, err := mp4.CreateMp4Muxer(out)
	if err != nil {
		return err
	}

	var (
		vtid, atid         uint32
		hasVideo, hasAudio bool
		writeErr           error
	)
	demuxer := mpeg2.NewTSDemuxer()
	demuxer.OnFrame = func(cid mpeg2.TS_STREAM_TYPE, frame []byte, pts uint64, dts uint64) {
		if writeErr != nil {
			return
		}
		switch cid {
		case mpeg2.TS_STREAM_H264:
			if !hasVideo {
				vtid = muxer.AddVideoTrack(mp4.MP4_CODEC_H264)
				hasVideo = true
			}
			writeErr = muxer.Write(vtid, frame, pts, dts)
		case mpeg2.TS_STREAM_H265:
			if !hasVideo {
				vtid = muxer.AddVideoTrack(mp4.MP4_CODEC_H265)
				hasVideo = true
			}
			writeErr = muxer.Write(vtid, frame, pts, dts)
		case mpeg2.TS_STREAM_AAC:
			if !hasAudio {
				atid = muxer.AddAudioTrack(mp4.MP4_CODEC_AAC)
				hasAudio = true
			}
			writeErr = muxer.Write(atid, frame, pts, dts)
		case mpeg2.TS_STREAM_AUDIO_MPEG1, mpeg2.TS_STREAM_AUDIO_MPEG2:
			if !hasAudio {
				atid = muxer.AddAudioTrack(mp4.MP4_CODEC_MP3)
				hasAudio = true
			}
			writeErr = muxer.Write(atid, frame, pts, dts)
		}
	}

	src := &lazyFileReader{paths: segments}
	defer func() { _ = src.Close() }()

	if err = demuxer.Input(src); err != nil {
		return err
	}
	if writeErr != nil {
		return writeErr
	}
	if !hasVideo && !hasAudio {
		return errors.New("no supported audio/video stream in ts segments")
	}
	return muxer.WriteTrailer()
}
