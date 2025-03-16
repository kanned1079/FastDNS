package resolver

import (
	"fast-dns-server/internal/config"
	"github.com/miekg/dns"
	"sync"
)

var msgPool = sync.Pool{
	New: func() any {
		return new(dns.Msg)
	},
}

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
	if config.RootCfg.Details.Config.Edns {
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
	if config.RootCfg.Details.Config.Edns {
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
