package table

import "sync"

// DownloadDataLen 存放下载进度
var DownloadDataLen = cmpMap[uint]{
	lock: sync.RWMutex{},
	body: map[string]uint{},
}

// CryptoVideoTable 存放视频加密的密钥
var CryptoVideoTable = sliceMap[[]byte]{
	lock: sync.RWMutex{},
	body: make(map[string][]byte),
}

// TitleData 标题
var TitleData = cmpMap[string]{
	lock: sync.RWMutex{},
	body: make(map[string]string),
}
