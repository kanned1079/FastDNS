package resolver

import (
	"fast-dns-server/internal/utils"
	"math/rand"
	"time"
)

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

// 获取加权随机选择的 DNS 服务器
func (this *DnsServerInst) getWeightedDnsServer() string {
	var totalWeight float64
	var weightedServers []DnsServer

	// 计算所有服务器的总权重
	for _, server := range dnsServersStats {
		weight := calculateWeight(server)
		totalWeight += weight
		weightedServers = append(weightedServers, *server)
	}

	// 随机选择一个服务器
	rand.Seed(time.Now().UnixNano())
	randomValue := rand.Float64() * totalWeight

	// 根据权重随机选择 DNS 服务器
	for _, server := range weightedServers {
		randomValue -= calculateWeight(&server)
		if randomValue <= 0 {
			// 返回选中的服务器地址
			return server.Address
		}
	}

	// 如果没有选择到，则返回第一个 DNS 服务器
	return weightedServers[0].Address
}

// 更新 DNS 服务器的查询统计信息
func (this *DnsServerInst) updateDnsStats(server string, success bool, responseTime time.Duration) {
	stats, exists := dnsServersStats[server]
	if !exists {
		stats = &DnsServer{Address: server}
		dnsServersStats[server] = stats
	}

	// 更新统计数据
	if success {
		stats.Stats.SuccessCount++
		stats.Stats.AvgResponseTime = (stats.Stats.AvgResponseTime*time.Duration(stats.Stats.SuccessCount-1) + responseTime) / time.Duration(stats.Stats.SuccessCount)
	} else {
		stats.Stats.FailureCount++
	}
}
