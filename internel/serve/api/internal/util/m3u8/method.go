package m3u8

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ParseM3u8Data(reader io.Reader) (*M3u8, error) {
	s := bufio.NewScanner(reader)
	var lines []string
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	var (
		i     = 0
		count = len(lines)
		m3u8  = &M3u8{
			Keys: make(map[int]*Key),
		}
		keyIndex = 0

		key     *Key
		seg     *Segment
		extInf  bool
		extByte bool
	)

	for ; i < count; i++ {
		line := strings.TrimSpace(lines[i])
		if i == 0 {
			if "#EXTM3U" != line {
				return nil, fmt.Errorf("invalid m3u8, missing #EXTM3U in line 1")
			}
			continue
		}
		switch {
		case line == "":
			continue
		case strings.HasPrefix(line, "#EXT-X-PLAYLIST-TYPE:"):
			if _, err := fmt.Sscanf(line, "#EXT-X-PLAYLIST-TYPE:%s", &m3u8.PlaylistType); err != nil {
				return nil, err
			}
			isValid := m3u8.PlaylistType == "" || m3u8.PlaylistType == PlaylistTypeVOD || m3u8.PlaylistType == PlaylistTypeEvent
			if !isValid {
				return nil, fmt.Errorf("invalid playlist type: %s, line: %d", m3u8.PlaylistType, i+1)
			}
		case strings.HasPrefix(line, "#EXT-X-TARGETDURATION:"):
			if _, err := fmt.Sscanf(line, "#EXT-X-TARGETDURATION:%f", &m3u8.TargetDuration); err != nil {
				return nil, err
			}
		case strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:"):
			if _, err := fmt.Sscanf(line, "#EXT-X-MEDIA-SEQUENCE:%d", &m3u8.MediaSequence); err != nil {
				return nil, err
			}
		case strings.HasPrefix(line, "#EXT-X-VERSION:"):
			if _, err := fmt.Sscanf(line, "#EXT-X-VERSION:%d", &m3u8.Version); err != nil {
				return nil, err
			}
		// Parse master playlist
		case strings.HasPrefix(line, "#EXT-X-STREAM-INF:"):
			mp, err := parseMasterPlaylist(line)
			if err != nil {
				return nil, err
			}
			i++
			mp.URI = lines[i]
			if mp.URI == "" || strings.HasPrefix(mp.URI, "#") {
				return nil, fmt.Errorf("invalid EXT-X-STREAM-INF URI, line: %d", i+1)
			}
			m3u8.MasterPlaylist = append(m3u8.MasterPlaylist, mp)
			continue
		case strings.HasPrefix(line, "#EXTINF:"):
			if extInf {
				return nil, fmt.Errorf("duplicate EXTINF: %s, line: %d", line, i+1)
			}
			if seg == nil {
				seg = new(Segment)
			}
			var s string
			if _, err := fmt.Sscanf(line, "#EXTINF:%s", &s); err != nil {
				return nil, err
			}
			if strings.Contains(s, ",") {
				split := strings.Split(s, ",")
				seg.Title = split[1]
				s = split[0]
			}
			df, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return nil, err
			}
			seg.Duration = float32(df)
			seg.KeyIndex = keyIndex
			extInf = true
		case strings.HasPrefix(line, "#EXT-X-BYTERANGE:"):
			if extByte {
				return nil, fmt.Errorf("duplicate EXT-X-BYTERANGE: %s, line: %d", line, i+1)
			}
			if seg == nil {
				seg = new(Segment)
			}
			var b string
			if _, err := fmt.Sscanf(line, "#EXT-X-BYTERANGE:%s", &b); err != nil {
				return nil, err
			}
			if b == "" {
				return nil, fmt.Errorf("invalid EXT-X-BYTERANGE, line: %d", i+1)
			}
			if strings.Contains(b, "@") {
				split := strings.Split(b, "@")
				offset, err := strconv.ParseUint(split[1], 10, 64)
				if err != nil {
					return nil, err
				}
				seg.Offset = uint64(offset)
				b = split[0]
			}
			length, err := strconv.ParseUint(b, 10, 64)
			if err != nil {
				return nil, err
			}
			seg.Length = uint64(length)
			extByte = true
		// Parse segments URI
		case !strings.HasPrefix(line, "#"):
			if extInf {
				if seg == nil {
					return nil, fmt.Errorf("invalid line: %s", line)
				}
				seg.URI = line
				extByte = false
				extInf = false
				m3u8.Segments = append(m3u8.Segments, seg)
				seg = nil
				continue
			}
		// Parse key
		case strings.HasPrefix(line, "#EXT-X-KEY"):
			params := parseLineParameters(line)
			if len(params) == 0 {
				return nil, fmt.Errorf("invalid EXT-X-KEY: %s, line: %d", line, i+1)
			}
			method := CryptMethod(params["METHOD"])
			if method != "" && method != CryptMethodAES && method != CryptMethodNONE {
				return nil, fmt.Errorf("invalid EXT-X-KEY method: %s, line: %d", method, i+1)
			}
			keyIndex++
			key = new(Key)
			key.Method = method
			key.URI = params["URI"]
			key.IV = params["IV"]
			m3u8.Keys[keyIndex] = key
		case line == "#EndList":
			m3u8.EndList = true
		default:
			continue
		}
	}

	return m3u8, nil
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

func MergeFiles(saveDir string) error {
	path, dirName := filepath.Split(saveDir)
	outputFilepath := filepath.Join(path, dirName+".mp4")
	outputFile, err := os.Create(outputFilepath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	var size int64 = 0
	err = filepath.Walk(saveDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		size += info.Size()

		inputFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer inputFile.Close()

		if _, err := io.Copy(outputFile, inputFile); err != nil {
			return err
		}
		return nil
	})

	//fileSize, _ := table.M3u8DownloadDataLen.Get(dirName)
	//fmt.Println(size, fileSize)

	return err
}
