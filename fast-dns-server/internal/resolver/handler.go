package resolver

import (
	"errors"
	"fast-dns-server/internal/config"
	"fast-dns-server/internal/utils"
	"fmt"
	"github.com/miekg/dns"
	"strings"
	"sync"
	"time"
)

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
		utils.GlobeUtils.ShowStatueLog("warning", "ATTENTION", "no DNS servers available")
		return nil, errors.New("no DNS servers available") // 如果都为空，返回错误
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
	utils.GlobeUtils.ShowStatueLog("warning", "ATTENTION", "failed to query DNS for "+domain)
	return nil, errors.New("failed to query DNS for " + domain)
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
		utils.GlobeUtils.ShowStatueLog("success", "FOUND", fmt.Sprintf("fetch from cache: %s latency: %v", domain, time.Since(startTime)))

		return result, nil
	}

	// 如果缓存中没有该数据，增加未命中次数
	utils.GlobeUtils.CacheMisses++

	// 打印请求到返回的总耗时
	//log.Printf("[CACHE NOT EXIST] no cache found, all latency: %v", time.Since(startTime))
	utils.GlobeUtils.ShowStatueLog("primary", "RESOLVED", fmt.Sprintf("no cache found, query upstream server success, all latency: %v", time.Since(startTime)))

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
					//log.Printf("[ERROR] cannot resolve %v\n", err)
					utils.GlobeUtils.ShowStatueLog("error", "FAILURE", "cannot resolve: "+err.Error())
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
