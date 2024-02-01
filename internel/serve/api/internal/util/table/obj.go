package table

import "sync"

// DownloadDataLen 存放下载进度
var DownloadDataLen = cmpMap[uint, uint]{
	lock: sync.RWMutex{},
	body: map[uint]uint{},
}

// CryptoVideoTable 存放视频加密的密钥
var CryptoVideoTable = sliceMap[uint, []byte]{
	lock: sync.RWMutex{},
	body: make(map[uint][]byte),
}

var ProxyCatchUrl = cmpMap[string, uint]{
	lock: sync.RWMutex{},
	body: make(map[string]uint),
}

var ProxyCatchHtml = cmpMap[string, string]{
	lock: sync.RWMutex{},
	body: make(map[string]string),
}

var DownloadTaskScore = cmpMap[uint, uint]{
	lock: sync.RWMutex{},
	body: make(map[uint]uint),
}
