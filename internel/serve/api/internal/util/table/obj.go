package table

import "sync"

// CryptoVideoTable 存放视频加密的密钥
var CryptoVideoTable = sliceMap[uint, []byte]{
	lock: sync.RWMutex{},
	body: make(map[uint][]byte),
}

var ProxyCatchUrl = cmpMap[string, uint]{
	lock: sync.RWMutex{},
	body: make(map[string]uint),
}

// ProxyCatchHtmlTitle 用于获取title
var ProxyCatchHtmlTitle = cmpMap[string, string]{
	lock: sync.RWMutex{},
	body: make(map[string]string),
}

var DownloadTimeSince = cmpMap[uint, uint]{
	lock: sync.RWMutex{},
	body: make(map[uint]uint),
}

// DownloadTaskByteLength 当前已经下载长度
var DownloadTaskByteLength = cmpMap[uint, uint]{
	lock: sync.RWMutex{},
	body: make(map[uint]uint),
}

// DownloadTaskMaxLength 每个任务的文件大小
var DownloadTaskMaxLength = cmpMap[uint, uint]{
	lock: sync.RWMutex{},
	body: make(map[uint]uint),
}
