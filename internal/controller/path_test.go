package controller

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestSanitizeName(t *testing.T) {
	cases := map[string]string{
		"normal title":           "normal title",
		"../etc/passwd":          "_etc_passwd", // ../ 头被 Trim(".") 去掉，/→_
		"a/b\\c":                 "a_b_c",
		`"<>:?*|`:                "_______",
		"":                       "untitled",
		"   ":                    "untitled",
		"...":                    "untitled",
		strings.Repeat("x", 200): strings.Repeat("x", 120),
	}
	for in, want := range cases {
		got := sanitizeName(in)
		if got != want {
			t.Errorf("sanitizeName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSafeJoin_NoEscape(t *testing.T) {
	dir := t.TempDir()
	// 即使 task.Name 完全是 ../，结果也不能逃出 dir。
	full := safeJoin(dir, "../../../etc/passwd")
	absDir, _ := filepath.Abs(dir)
	absFull, _ := filepath.Abs(full)
	if !strings.HasPrefix(absFull, absDir) {
		t.Errorf("safeJoin escaped: %s not under %s", absFull, absDir)
	}
}
