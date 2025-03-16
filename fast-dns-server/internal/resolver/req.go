package resolver

import (
	"fast-dns-server/internal/config"
	"fast-dns-server/internal/utils"
	"fmt"
	"github.com/miekg/dns"
)

// GetDNSInfo 解析 DNS 请求
func (this *DnsServerInst) GetDNSInfo(domain string) ([]dns.RR, bool) {
	// 先检查缓存
	utils.GlobeUtils.ShowStatueLog("cyan", "SEARCH", "domain: "+domain)
	if data, err := this.cache.Get(domain); err == nil {
		if rrResult, ok := data.([]dns.RR); ok {
			return rrResult, true
		}
	}

	// 动态加载配置信息
	dnsMode := config.RootCfg.Details.Dns.Mode
	// 如果缓存中没有数据，根据模式进行 DNS 查询
	var result []dns.RR
	var err error
	switch dnsMode {
	case "balance":
		result, err = this.queryDnsWithBalance(domain)
	case "parallel":
		result, err = this.queryDnsParallel(domain)
	default:
		utils.GlobeUtils.ShowStatueLog("warning", "FAILURE", fmt.Sprintf("Unsupported DNS query mode: %s", dnsMode))
	}

	if err != nil {
		utils.GlobeUtils.ShowStatueLog("error", "FAILURE", fmt.Sprintf("DNS query failed for %s: %v", domain, err))
		return nil, false
	}

	// 将查询结果存入缓存
	if len(result) > 0 {
		this.cache.Set(domain, result)
	}

	return result, false
}
