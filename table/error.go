package table

import (
	"dv/config"
	"sync"
)

var errorCountTable sync.Map

func IncreaseErrorCount(key string) {
	count, exist := errorCountTable.Load(key)
	if exist {
		errorCountTable.Store(key, count.(uint)+1)
	} else {
		errorCountTable.Store(key, uint(1))
	}
}

func GetErrorCount(key string) uint {
	count, exist := errorCountTable.Load(key)
	if !exist {
		return 0
	}

	return count.(uint)
}

func IsMaxErrorCount(key string) bool {
	count, exist := errorCountTable.Load(key)
	if !exist {
		return false
	}

	return count.(uint) > config.GetConfig().TaskErrorMaxCount
}

var errorTask sync.Map

func AddErrorTask(key string) {
	errorTask.Store(key, struct{}{})
}

func RangeErrorTask() []string {
	ts := make([]string, 0)
	errorTask.Range(func(key, value interface{}) bool {
		ts = append(ts, key.(string))
		return true
	})
	return ts
}
