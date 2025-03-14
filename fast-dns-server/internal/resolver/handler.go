package resolver

import (
	"fast-dns-server/internal/utils"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"time"
)

// 获取缓存命中率
//func (this *DnsServerInst) GetCacheHitRatio() float64 {
//	total := this.CacheHits + this.CacheMisses
//	if total == 0 {
//		return 0
//	}
//	return float64(this.CacheHits) / float64(total)
//}

// resolver 将请求转发到上游 DNS 服务器并缓存结果
func (this *DnsServerInst) resolver(domain string, qtype uint16) ([]dns.RR, error) {
	// 检查缓存
	if cachedResult, err := this.cache.Get(domain); err == nil {
		// 如果缓存中存在查询结果
		utils.GlobeUtils.CacheHits++ // 增加缓存命中次数
		log.Println("从缓存中获取:", domain)
		return cachedResult.([]dns.RR), nil
	}

	// 如果缓存中没有该数据，增加未命中次数
	utils.GlobeUtils.CacheMisses++

	// 创建 DNS 请求
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), qtype)
	m.RecursionDesired = true

	// 创建一个 DNS 客户端，并设置超时
	c := &dns.Client{Timeout: 5 * time.Second}

	response, _, err := c.Exchange(m, "223.5.5.5:53")
	if err != nil {
		log.Printf("[ERROR] DNS 请求失败: %v\n", err)
		return nil, err
	}

	if response == nil {
		log.Printf("[ERROR] 没有收到响应\n")
		return nil, fmt.Errorf("no response")
	}
	log.Println("上游 DNS 查询结果:", response.Answer)

	// 将结果存入缓存
	this.cache.Set(domain, response.Answer)

	return response.Answer, nil
}

// 处理 DNS 请求

// 处理 DNS 请求
func (this *DnsServerInst) handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	// 创建一个响应消息
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	// 处理请求的每一个问题
	for _, question := range r.Question {
		// 如果是 A 记录的查询请求
		if question.Qtype == dns.TypeA {
			// 使用上游 DNS 进行解析
			answers, err := this.resolver(question.Name, question.Qtype)
			if err != nil {
				// 如果解析失败，可以返回一个自定义的错误信息或继续返回空值
				log.Printf("DNS 解析失败: %v\n", err)
				continue
			}

			// 将上游 DNS 服务器返回的结果添加到响应中
			m.Answer = append(m.Answer, answers...)
		}
	}

	// 将响应发送回客户端
	w.WriteMsg(m)
}
