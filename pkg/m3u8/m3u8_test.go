package m3u8

import (
	"bytes"
	"testing"
)

func TestParseM3u8Data_Simple(t *testing.T) {
	data := `#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:10
#EXTINF:9.009,
segment0.ts
#EXTINF:9.009,
segment1.ts
#EXT-X-ENDLIST`

	m, err := ParseM3u8Data(bytes.NewReader([]byte(data)))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if len(m.Segments) != 2 {
		t.Errorf("expected 2 segments, got %d", len(m.Segments))
	}
	if m.Segments[0].URI != "segment0.ts" {
		t.Errorf("expected segment0.ts, got %s", m.Segments[0].URI)
	}
}

func TestSegmentExtension(t *testing.T) {
	cases := map[string]string{
		"https://example.com/a/segment0.ts":           ".ts",
		"https://example.com/a/segment0.ts?token=abc": ".ts",
		"https://example.com/a/segment0.M4A#frag":     ".m4a",
		"https://example.com/a/segment0.m4s":          ".m4s",
		"https://example.com/a/segment0.aac":          ".aac",
		"https://example.com/a/segment0.mp4":          ".mp4",
		"https://example.com/a/segment0":              ".ts",
		"https://example.com/a/segment0.php":          ".ts",
		"segment0.m4a":                                ".m4a",
	}
	for in, want := range cases {
		if got := SegmentExtension(in); got != want {
			t.Errorf("SegmentExtension(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestIsSegmentFile(t *testing.T) {
	cases := map[string]bool{
		"000000.ts":  true,
		"000123.m4a": true,
		"000001.m4s": true,
		"init.mp4":   false,
		"tomato.txt": false,
		"000000.txt": false,
		"00000.ts":   false,
		"000000.TS":  true,
		"local.m3u8": false,
	}
	for in, want := range cases {
		if got := IsSegmentFile(in); got != want {
			t.Errorf("IsSegmentFile(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestParseM3u8Data_ExtXMap(t *testing.T) {
	data := `#EXTM3U
#EXT-X-VERSION:7
#EXT-X-TARGETDURATION:10
#EXT-X-KEY:METHOD=AES-128,URI="key.key",IV=0x00000000000000000000000000000001
#EXT-X-MAP:URI="init.mp4"
#EXTINF:10,
seg0.m4s
#EXTINF:10,
seg1.m4s
#EXT-X-ENDLIST`

	m, err := ParseM3u8Data(bytes.NewReader([]byte(data)))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if m.InitSegment == nil {
		t.Fatal("expected init segment")
	}
	if m.InitSegment.URI != "init.mp4" {
		t.Errorf("expected init.mp4, got %s", m.InitSegment.URI)
	}
	if m.InitSegment.KeyIndex != 1 {
		t.Errorf("expected init key index 1, got %d", m.InitSegment.KeyIndex)
	}
	if len(m.Segments) != 2 {
		t.Errorf("expected 2 segments, got %d", len(m.Segments))
	}
}

func TestParseM3u8Data_ExtXMapMissingURI(t *testing.T) {
	data := `#EXTM3U
#EXT-X-MAP:BYTERANGE="100@0"
#EXTINF:10,
seg0.m4s`

	if _, err := ParseM3u8Data(bytes.NewReader([]byte(data))); err == nil {
		t.Error("expected error for EXT-X-MAP without URI")
	}
}

func TestParseM3u8Data_MasterPlaylist(t *testing.T) {
	data := `#EXTM3U
#EXT-X-STREAM-INF:BANDWIDTH=1280000
low.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=2560000
high.m3u8`

	m, err := ParseM3u8Data(bytes.NewReader([]byte(data)))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if !m.IsMasterPlaylist() {
		t.Error("expected master playlist")
	}
	if len(m.MasterPlaylist) != 2 {
		t.Errorf("expected 2 streams, got %d", len(m.MasterPlaylist))
	}

	idx := m.GetMaxBandWidth()
	if idx != 1 {
		t.Errorf("expected max bandwidth index 1, got %d", idx)
	}
}
