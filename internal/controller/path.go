package controller

import (
	"path/filepath"
	"regexp"
	"strings"
)

// 不允许出现在文件/目录名里的字符（Windows 保留 + 路径分隔符 + 控制字符）。
var unsafeFileChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)

// sanitizeName 清洗任务名，使其可安全作为目录或文件名的一部分使用。
// - 移除路径分隔符与 Windows 保留字符
// - 去掉前后空白、点
// - 截断到 120 字节
// 不为空白结果保留 fallback。
func sanitizeName(raw string) string {
	s := unsafeFileChars.ReplaceAllString(raw, "_")
	s = strings.TrimSpace(s)
	s = strings.Trim(s, ".")
	if len(s) > 120 {
		s = s[:120]
	}
	if s == "" {
		s = "untitled"
	}
	return s
}

// safeJoin 把清洗后的 name 拼接到 dir 之下，并断言结果仍在 dir 内。
// 若越界（极端情况），返回 dir/untitled。
func safeJoin(dir, name string) string {
	clean := sanitizeName(name)
	full := filepath.Join(dir, clean)
	absDir, _ := filepath.Abs(dir)
	absFull, _ := filepath.Abs(full)
	if !strings.HasPrefix(absFull, absDir) {
		return filepath.Join(dir, "untitled")
	}
	return full
}
