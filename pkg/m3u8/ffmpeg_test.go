package m3u8

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func writeTestFiles(t *testing.T, dir string, names ...string) {
	t.Helper()
	for _, name := range names {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0644); err != nil {
			t.Fatalf("write %s failed: %v", name, err)
		}
	}
}

func TestCollectSegments(t *testing.T) {
	dir := t.TempDir()
	writeTestFiles(t, dir, "000001.ts", "000000.m4a", "init.mp4", "list.txt", ".DS_Store")

	segments, hasInit, audioOnly, err := collectSegments(dir)
	if err != nil {
		t.Fatalf("collect failed: %v", err)
	}
	want := []string{"000000.m4a", "000001.ts"}
	if !reflect.DeepEqual(segments, want) {
		t.Errorf("segments = %v, want %v", segments, want)
	}
	if !hasInit {
		t.Error("expected hasInit")
	}
	if audioOnly {
		t.Error("expected mixed audio/video, audioOnly should be false")
	}
}

func TestCollectSegments_AudioOnly(t *testing.T) {
	dir := t.TempDir()
	writeTestFiles(t, dir, "000000.m4a", "000001.m4a", "000002.aac")

	_, hasInit, audioOnly, err := collectSegments(dir)
	if err != nil {
		t.Fatalf("collect failed: %v", err)
	}
	if hasInit {
		t.Error("expected no init segment")
	}
	if !audioOnly {
		t.Error("expected audioOnly")
	}
}

func TestCollectSegments_Fmp4WithoutInit(t *testing.T) {
	dir := t.TempDir()
	writeTestFiles(t, dir, "000000.m4s")

	if _, _, _, err := collectSegments(dir); err == nil {
		t.Error("expected error for fMP4 fragments without init segment")
	}
}

func TestCollectSegments_Empty(t *testing.T) {
	dir := t.TempDir()

	if _, _, _, err := collectSegments(dir); err == nil {
		t.Error("expected error for empty directory")
	}
}
