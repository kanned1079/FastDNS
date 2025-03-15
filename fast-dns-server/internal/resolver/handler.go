package resolver

import (
	"fast-dns-server/internal/config"
	"fast-dns-server/internal/utils"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var msgPool = sync.Pool{
	New: func() any {
		return new(dns.Msg)
	},
}

// GetDNSInfo 解析 DNS 请求
func (this *DnsServerInst) GetDNSInfo(domain string) ([]dns.RR, bool) {
	// 先检查缓存
	log.Println("domain: ", domain)
	if data, err := this.cache.Get(domain); err == nil {
		if rrResult, ok := data.([]dns.RR); ok {
			return rrResult, true
		}
	}

	// 动态加载配置信息
	dnsMode := config.RootCfg.Details.Dns.Mode
	log.Println("请求方式 ", dnsMode)
	// 如果缓存中没有数据，根据模式进行 DNS 查询
	var result []dns.RR
	var err error
	switch dnsMode {
	case "balance":
		result, err = this.queryDnsWithBalance(domain)
	case "parallel":
		result, err = this.queryDnsParallel(domain)
	default:
		log.Printf("Unsupported DNS query mode: %s", dnsMode)
	}

	if err != nil {
		log.Printf("DNS query failed for %s: %v", domain, err)
		return nil, false
	}

	// 将查询结果存入缓存
	if len(result) > 0 {
		this.cache.Set(domain, result)
	}

	return result, false
}

// 加权随机选择 DNS 服务器
func (this *DnsServerInst) queryDnsBalance(domain string) ([]dns.RR, error) {
	// 动态加载配置信息
	dnsList := config.RootCfg.Details.Dns.List
	backupSecondaryList := config.RootCfg.Details.Dns.BackupSecondaryList

	var dnsServer string
	if len(dnsList) > 0 {
		dnsServer = dnsList[0] // 简化逻辑，选择第一个 DNS
	} else {
		dnsServer = backupSecondaryList[0] // 使用备用 DNS
	}

	// 根据协议选择查询方式
	if strings.HasPrefix(dnsServer, "tls://") {
		// 使用 DNS over TLS
		return this.queryDnsOverTls(domain, dnsServer[5:])
	} else {
		// 普通 DNS 查询
		return this.queryDns(domain, dnsServer)
	}
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

// 通过加权随机算法查询 DNS
func (this *DnsServerInst) queryDnsWithBalance(domain string) ([]dns.RR, error) {
	// 动态加载配置信息
	dnsList := config.RootCfg.Details.Dns.List
	backupSecondaryList := config.RootCfg.Details.Dns.BackupSecondaryList

	// 确保 dnsList 或 backupSecondaryList 至少有一个有有效的 DNS 服务器
	var dnsServer string
	if len(dnsList) > 0 {
		dnsServer = dnsList[0] // 从 dnsList 选择第一个 DNS
	} else if len(backupSecondaryList) > 0 {
		dnsServer = backupSecondaryList[0] // 从 backupSecondaryList 选择第一个备用 DNS
	} else {
		return nil, fmt.Errorf("no DNS servers available") // 如果都为空，返回错误
	}

	// 根据协议选择查询方式
	if strings.HasPrefix(dnsServer, "tls://") {
		// 使用 DNS over TLS
		return this.queryDnsOverTls(domain, dnsServer[5:])
	} else {
		// 普通 DNS 查询
		return this.queryDns(domain, dnsServer)
	}
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

// 通过并行模式查询所有 DNS 服务器
func (this *DnsServerInst) queryDnsParallel(domain string) ([]dns.RR, error) {
	// 动态加载配置信息
	dnsList := config.RootCfg.Details.Dns.List
	if len(config.RootCfg.Details.Dns.List) <= 0 {
		dnsList = append(dnsList, config.RootCfg.Details.Dns.BackupSecondaryList...)
	}
	var wg sync.WaitGroup
	var result []dns.RR
	var err error

	// 并行查询所有 DNS 服务器
	var mu sync.Mutex
	var once sync.Once

	for _, dnsServer := range dnsList {
		wg.Add(1)
		go func(server string) {
			//log.Println("当前请求的上游服务器 ", server)
			defer wg.Done()

			var res []dns.RR
			var queryErr error
			if strings.HasPrefix(server, "tls://") {
				// 使用 DNS over TLS
				res, queryErr = this.queryDnsOverTls(domain, server[5:])
			} else {
				// 普通 DNS 查询
				res, queryErr = this.queryDns(domain, server)
			}

			// 获取到结果就锁定并返回
			if queryErr == nil && len(res) > 0 {
				once.Do(func() {
					mu.Lock()
					if result == nil { // 第一个成功的结果
						result = res
					}
					mu.Unlock()
				})
			} else if queryErr != nil && err == nil {
				err = queryErr
			}
		}(dnsServer)
	}

	wg.Wait()
	if result != nil {
		return result, nil
	}
	return nil, fmt.Errorf("failed to query DNS for %s", domain)
}

// 普通 DNS 查询
//func (this *DnsServerInst) queryDns(domain, dnsServer string) ([]dns.RR, error) {
//	// 从池中获取 dns.Msg 对象
//	message := msgPool.Get().(*dns.Msg)
//
//	// 手动重置 message 字段
//	message.Answer = nil            // 清空已有的答案
//	message.Ns = nil                // 清空权威域名服务器列表
//	message.Extra = nil             // 清空额外信息
//	message.Question = nil          // 清空查询问题
//	message.Id = 0                  // 重置查询 ID
//	message.RecursionDesired = true // 保持递归查询标志
//
//	// 设置查询问题
//	message.SetQuestion(dns.Fqdn(domain), dns.TypeA)
//
//	// 查询 DNS
//	response, _, err := this.client.Exchange(message, dnsServer+":53")
//	if err != nil {
//		// 使用完后将 message 归还池中
//		msgPool.Put(message)
//		return nil, err
//	}
//
//	// 使用完后将 message 归还池中
//	msgPool.Put(message)
//
//	// 返回查询结果
//	return response.Answer, nil
//
//}

func (this *DnsServerInst) queryDns(domain, dnsServer string) ([]dns.RR, error) {
	// 从池中获取 dns.Msg 对象
	message := msgPool.Get().(*dns.Msg)

	// 手动重置 message 字段
	message.Answer = nil            // 清空已有的答案
	message.Ns = nil                // 清空权威域名服务器列表
	message.Extra = nil             // 清空额外信息
	message.Question = nil          // 清空查询问题
	message.Id = 0                  // 重置查询 ID
	message.RecursionDesired = true // 保持递归查询标志

	// 设置查询问题
	message.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	// 根据配置启用 EDNS
	if config.RootCfg.Details.Config.Ends {
		// 启用 EDNS0，设置最大接收缓冲区为 4096 字节
		message.SetEdns0(4096, true)
	}

	// 根据配置启用 DNSSEC
	if config.RootCfg.Details.Config.Dnssec {
		// 启用 DNSSEC，设置 AD (Authenticated Data) 标志
		message.SetEdns0(4096, true)
		message.AuthenticatedData = true
	}

	// 查询 DNS
	response, _, err := this.client.Exchange(message, dnsServer+":53")
	if err != nil {
		// 使用完后将 message 归还池中
		msgPool.Put(message)
		return nil, err
	}

	// 使用完后将 message 归还池中
	msgPool.Put(message)

	// 返回查询结果
	return response.Answer, nil
}

// DNS over TLS 查询
//func (this *DnsServerInst) queryDnsOverTls(domain, dnsServer string) ([]dns.RR, error) {
//	// 从池中获取 dns.Msg 对象
//	message := msgPool.Get().(*dns.Msg)
//
//	// 手动重置 message 字段
//	message.Answer = nil            // 清空已有的答案
//	message.Ns = nil                // 清空权威域名服务器列表
//	message.Extra = nil             // 清空额外信息
//	message.Question = nil          // 清空查询问题
//	message.Id = 0                  // 重置查询 ID
//	message.RecursionDesired = true // 保持递归查询标志
//
//	message.SetQuestion(dns.Fqdn(domain), dns.TypeA)
//
//	// 使用 DNS over TLS 进行查询
//	tlsAddr := dnsServer + ":853"
//	response, _, err := this.client.Exchange(message, tlsAddr)
//	if err != nil {
//		// 使用完后将 message 归还池中
//		msgPool.Put(message)
//		return nil, err
//	}
//
//	// 使用完后将 message 归还池中
//	msgPool.Put(message)
//
//	return response.Answer, nil
//}

func (this *DnsServerInst) queryDnsOverTls(domain, dnsServer string) ([]dns.RR, error) {
	// 从池中获取 dns.Msg 对象
	message := msgPool.Get().(*dns.Msg)

	// 手动重置 message 字段
	message.Answer = nil            // 清空已有的答案
	message.Ns = nil                // 清空权威域名服务器列表
	message.Extra = nil             // 清空额外信息
	message.Question = nil          // 清空查询问题
	message.Id = 0                  // 重置查询 ID
	message.RecursionDesired = true // 保持递归查询标志

	message.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	// 根据配置启用 EDNS
	if config.RootCfg.Details.Config.Ends {
		// 启用 EDNS0，设置最大接收缓冲区为 4096 字节
		message.SetEdns0(4096, true)
	}

	// 根据配置启用 DNSSEC
	if config.RootCfg.Details.Config.Dnssec {
		// 启用 DNSSEC，设置 AD (Authenticated Data) 标志
		message.SetEdns0(4096, true)
		message.AuthenticatedData = true
	}

	// 使用 DNS over TLS 进行查询
	tlsAddr := dnsServer + ":853"
	response, _, err := this.client.Exchange(message, tlsAddr)
	if err != nil {
		// 使用完后将 message 归还池中
		msgPool.Put(message)
		return nil, err
	}

	// 使用完后将 message 归还池中
	msgPool.Put(message)

	return response.Answer, nil
}

// resolver方法
// 解析 DNS 请求
var startTime time.Time

func (this *DnsServerInst) resolver(domain string, qtype uint16) ([]dns.RR, error) {
	// 记录开始时间
	startTime = time.Now()

	// 调用 GetDNSInfo 来处理缓存和 DNS 查询
	result, found := this.GetDNSInfo(domain)
	if found {
		// 如果缓存中找到数据，返回缓存结果
		utils.GlobeUtils.CacheHits++
		log.Printf("从缓存中获取: %s, 耗时: %v", domain, time.Since(startTime)) // 打印从缓存获取的时间
		return result, nil
	}

	// 如果缓存中没有该数据，增加未命中次数
	utils.GlobeUtils.CacheMisses++

	// 打印请求到返回的总耗时
	log.Printf("从请求到返回的总耗时: %v", time.Since(startTime))

	// 返回查询结果
	return result, nil
}

func (this *DnsServerInst) handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	// 创建一个响应消息
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	// 创建一个 goroutine 来处理 DNS 请求
	var wg sync.WaitGroup
	mu := sync.Mutex{}
	for _, question := range r.Question {
		wg.Add(1)
		go func(q dns.Question) {
			defer wg.Done()
			// 如果是 A 记录的查询请求
			if q.Qtype == dns.TypeA {
				// 使用 GetDNSInfo 来解析 DNS 查询
				answers, err := this.resolver(q.Name, q.Qtype)
				if err != nil {
					// 如果解析失败，可以返回一个自定义的错误信息或继续返回空值
					log.Printf("DNS 解析失败: %v\n", err)
					return
				}

				// 锁住并安全地修改响应
				mu.Lock()
				m.Answer = append(m.Answer, answers...)
				mu.Unlock()
			}
		}(question)
	}

	// 等待所有的 goroutine 完成
	wg.Wait()

	// 将响应发送回客户端
	w.WriteMsg(m)
}
