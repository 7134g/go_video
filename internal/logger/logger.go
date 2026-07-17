package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	LogDir        = "logs"
	logPattern    = "2006-01-02"
	maxLogAge     = 30 * 24 * time.Hour
	cleanInterval = time.Hour
)

type dateRotatingWriter struct {
	mu        sync.Mutex
	dir       string
	curDate   string
	file      *os.File
	lastClean time.Time
}

func NewDateRotatingWriter(dir string) io.WriteCloser {
	return &dateRotatingWriter{dir: dir}
}

func (w *dateRotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	today := now.Format(logPattern)

	if now.Sub(w.lastClean) > cleanInterval {
		cleanOldLogs(w.dir)
		w.lastClean = now
	}

	if today != w.curDate {
		if w.file != nil {
			w.file.Close()
		}
		if err := os.MkdirAll(w.dir, 0755); err != nil {
			return 0, err
		}
		f, err := os.OpenFile(filepath.Join(w.dir, today+".log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return 0, err
		}
		w.file = f
		w.curDate = today
	}

	return w.file.Write(p)
}

func (w *dateRotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

func cleanOldLogs(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	now := time.Now()
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".log") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if now.Sub(info.ModTime()) > maxLogAge {
			os.Remove(filepath.Join(dir, info.Name()))
		}
	}
}

type multiHandler struct {
	handlers []slog.Handler
}

func (h *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, r.Level) {
			if err := handler.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (h *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}

func Init(level slog.Level, logDir string) {
	opts := &slog.HandlerOptions{Level: level}

	fileWriter := NewDateRotatingWriter(logDir)
	fileHandler := slog.NewJSONHandler(fileWriter, opts)

	textHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:       level,
		AddSource:   false,
		ReplaceAttr: removeTime,
	})

	handler := &multiHandler{
		handlers: []slog.Handler{textHandler, fileHandler},
	}
	slog.SetDefault(slog.New(handler))
}

func InitWithLevel(levelStr string, logDir string) {
	level := parseLevel(levelStr)
	Init(level, logDir)
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func removeTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		return slog.Attr{}
	}
	return a
}
