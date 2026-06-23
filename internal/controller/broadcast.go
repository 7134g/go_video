package controller

import (
	"sync"
)

type Message struct {
	TaskID  uint   `json:"task_id"`
	Message string `json:"message"`
}

var (
	msgListeners      []chan Message
	msgMu             sync.RWMutex
	progressListeners []chan ProgressInfo
	progressMu        sync.RWMutex
)

// BroadcastMessage 向所有 WebSocket 监听者非阻塞广播一条消息。
// 监听者 channel 已满时会丢弃该消息——优先保护 controller 不被慢消费者拖死。
func BroadcastMessage(taskID uint, msg string) {
	msgMu.RLock()
	defer msgMu.RUnlock()

	m := Message{TaskID: taskID, Message: msg}
	for _, ch := range msgListeners {
		select {
		case ch <- m:
		default:
		}
	}
}

func AddMessageListener(ch chan Message) {
	msgMu.Lock()
	defer msgMu.Unlock()
	msgListeners = append(msgListeners, ch)
}

func RemoveMessageListener(ch chan Message) {
	msgMu.Lock()
	defer msgMu.Unlock()
	for i, c := range msgListeners {
		if c == ch {
			msgListeners = append(msgListeners[:i], msgListeners[i+1:]...)
			break
		}
	}
}

// BroadcastProgress 向所有 WebSocket 监听者非阻塞广播一条进度更新。
func BroadcastProgress(info ProgressInfo) {
	progressMu.RLock()
	defer progressMu.RUnlock()
	for _, ch := range progressListeners {
		select {
		case ch <- info:
		default:
		}
	}
}

func AddProgressListener(ch chan ProgressInfo) {
	progressMu.Lock()
	defer progressMu.Unlock()
	progressListeners = append(progressListeners, ch)
}

func RemoveProgressListener(ch chan ProgressInfo) {
	progressMu.Lock()
	defer progressMu.Unlock()
	for i, c := range progressListeners {
		if c == ch {
			progressListeners = append(progressListeners[:i], progressListeners[i+1:]...)
			break
		}
	}
}
