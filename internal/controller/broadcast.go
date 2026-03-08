package controller

import (
	"sync"
)

type Message struct {
	TaskID  uint   `json:"task_id"`
	Message string `json:"message"`
}

var (
	msgListeners []chan Message
	msgMu        sync.RWMutex
)

func BroadcastMessage(taskID uint, msg string) {
	msgMu.RLock()
	defer msgMu.RUnlock()

	//fmt.Println("====================>", msg)

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
