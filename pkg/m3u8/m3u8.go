// Partial reference https://github.com/grafov/m3u8/blob/master/reader.go
package m3u8

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type (
	PlaylistType string
	CryptMethod  string
)

const (
	PlaylistTypeVOD   PlaylistType = "VOD"
	PlaylistTypeEvent PlaylistType = "EVENT"

	CryptMethodAES  CryptMethod = "AES-128"
	CryptMethodNONE CryptMethod = "NONE"
)

// regex pattern for extracting `key=value` parameters from a line
var linePattern = regexp.MustCompile(`([a-zA-Z-]+)=("[^"]+"|[^",]+)`)

type M3u8 struct {
	Version        int8   // EXT-X-VERSION:version
	MediaSequence  uint64 // Default 0, #EXT-X-MEDIA-SEQUENCE:sequence
	Segments       []*Segment
	MasterPlaylist []*MasterPlaylist
	Keys           map[int]*Key
	EndList        bool         // #EXT-X-ENDLIST
	PlaylistType   PlaylistType // VOD or EVENT
	TargetDuration float64      // #EXT-X-TARGETDURATION:duration
}

type Segment struct {
	URI      string
	KeyIndex int
	Title    string  // #EXTINF: duration,<title>
	Duration float32 // #EXTINF: duration,<title>
	Length   uint64  // #EXT-X-BYTERANGE: length[@offset]
	Offset   uint64  // #EXT-X-BYTERANGE: length[@offset]
}

// #EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=240000,RESOLUTION=416x234,CODECS="avc1.42e00a,mp4a.40.2"
type MasterPlaylist struct {
	URI        string
	BandWidth  uint32
	Resolution string
	Codecs     string
	ProgramID  uint32
}

// #EXT-X-KEY:METHOD=AES-128,URI="key.key"
type Key struct {
	// 'AES-128' or 'NONE'
	// If the encryption method is NONE, the URI and the IV attributes MUST NOT be present
	Method CryptMethod
	URI    string
	IV     string
	// KeyData 存储下载后的密钥数据（16字节）
	KeyData []byte
}

// IsMasterPlaylist 判断是否为 Master Playlist
func (m *M3u8) IsMasterPlaylist() bool {
	return len(m.MasterPlaylist) > 0
}

// HasEncryption 判断是否包含加密分片
func (m *M3u8) HasEncryption() bool {
	for _, key := range m.Keys {
		if key.Method == CryptMethodAES {
			return true
		}
	}
	return false
}

// GetSegmentKey 获取分片对应的加密密钥
func (m *M3u8) GetSegmentKey(seg *Segment) *Key {
	if seg.KeyIndex == 0 {
		return nil
	}
	return m.Keys[seg.KeyIndex]
}

// TotalDuration 计算总时长
func (m *M3u8) TotalDuration() float32 {
	var total float32
	for _, seg := range m.Segments {
		total += seg.Duration
	}
	return total
}

func parseMasterPlaylist(line string) (*MasterPlaylist, error) {
	params := parseLineParameters(line)
	if len(params) == 0 {
		return nil, errors.New("empty parameter")
	}
	mp := new(MasterPlaylist)
	for k, v := range params {
		switch k {
		case "BANDWIDTH":
			v, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				return nil, err
			}
			mp.BandWidth = uint32(v)
		case "RESOLUTION":
			mp.Resolution = v
		case "PROGRAM-ID":
			v, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				return nil, err
			}
			mp.ProgramID = uint32(v)
		case "CODECS":
			mp.Codecs = v
		}
	}
	return mp, nil
}

// parseLineParameters extra parameters in string `line`
func parseLineParameters(line string) map[string]string {
	r := linePattern.FindAllStringSubmatch(line, -1)
	params := make(map[string]string)
	for _, arr := range r {
		params[arr[1]] = strings.Trim(arr[2], "\"")
	}
	return params
}

func (m *M3u8) GetMaxBandWidth() int {
	if len(m.MasterPlaylist) == 0 {
		return -1
	}

	var maxIndex int
	for i := 0; i < len(m.MasterPlaylist); i++ {
		if m.MasterPlaylist[i].BandWidth > m.MasterPlaylist[maxIndex].BandWidth {
			maxIndex = i
		}
	}

	return maxIndex
}
