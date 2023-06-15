package table

import (
	"dv/base"
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

func AddErrorTask(key string, value interface{}) {
	errorTask.Store(key, value)
}

func GetListErrorTask(key string) []interface{} {
	ts := make([]interface{}, 0)
	errorTask.Range(func(k, v any) bool {
		if base.ReplaceName(k.(string)) == key {
			ts = append(ts, v)
			errorTask.Delete(k)
		}

		return true
	})

	return ts
}

func RangeErrorTask() []interface{} {
	ts := make([]interface{}, 0)
	errorTask.Range(func(key, value interface{}) bool {
		ts = append(ts, key)
		return true
	})
	return ts
}
