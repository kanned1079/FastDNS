package utils

var GlobeUtils Utils

type Utils struct {
	CacheHits      int     // 命中次数
	CacheMisses    int     // 未命中次数
	CacheHitRatio  float64 // 缓存命中率
	colorCode      string  // 日志标题颜色
	currentTimeStr string  // 当前时间
}
