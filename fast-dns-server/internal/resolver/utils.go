package resolver

import "fast-dns-server/internal/utils"

// GetCacheInfo 获取缓存信息
func (this *DnsServerInst) GetCacheInfo() map[string]any {
	return map[string]any{
		"cache_hits":   utils.GlobeUtils.CacheHits,          // 缓存命中
		"cache_misses": utils.GlobeUtils.CacheMisses,        // 缓存未命中,
		"cache_rate":   utils.GlobeUtils.GetCacheHitRatio(), // 缓存命中率
	}
}

// 权重计算函数
func calculateWeight(server *DnsServer) float64 {
	// 权重 = (成功次数 + 1) / (失败次数 + 1) * (1 / 平均响应时间)
	successFactor := float64(server.Stats.SuccessCount + 1)
	failureFactor := float64(server.Stats.FailureCount + 1)
	timeFactor := 1 / server.Stats.AvgResponseTime.Seconds()

	// 权重
	weight := successFactor / failureFactor * timeFactor
	return weight
}
