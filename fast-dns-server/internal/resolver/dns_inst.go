package resolver

import (
	"fast-dns-server/internal/config"
	"github.com/bluele/gcache"
	"github.com/miekg/dns"
	"time"
)

type DnsServerInst struct {
	Id           int64
	resolverInst dns.Server
	cache        gcache.Cache // 使用 gcache 来缓存 DNS 查询结果

}

func NewDnsServerInst(id int64, addr, proc string) *DnsServerInst {
	return &DnsServerInst{
		Id: id,
		resolverInst: dns.Server{
			Addr: addr,
			Net:  proc,
		},
		cache: gcache.New(int(config.RootCfg.Details.Config.CacheSize)).
			Expiration(time.Second * time.Duration(1000)).
			ARC(). // 使用 ARC 算法 (可选，LRU 或 LFU 也可以)
			Build(),
	}
}
