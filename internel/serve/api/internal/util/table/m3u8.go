package table

import "sync"

// M3u8DownloadSpeed 存放下载进度
var M3u8DownloadSpeed = cmpMap[uint]{
	lock: sync.RWMutex{},
	body: map[string]uint{},
}

// CryptoVideoTable 存放视频加密的密钥
var CryptoVideoTable = sliceMap[[]byte]{
	lock: sync.RWMutex{},
	body: make(map[string][]byte),
}
