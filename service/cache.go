package service

import "github.com/silentred/template/util"

var CacheInMem *util.MemCache

func GetMemCache() util.Cache {
	if CacheInMem == nil {
		CacheInMem = util.NewMemCache(3600)
	}
	return CacheInMem
}
