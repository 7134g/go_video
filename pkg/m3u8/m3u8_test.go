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
