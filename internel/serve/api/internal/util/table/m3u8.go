package table

import "sync"

var M3u8DownloadSpeed = base{
	lock: sync.RWMutex{},
	body: map[string]any{},
}
