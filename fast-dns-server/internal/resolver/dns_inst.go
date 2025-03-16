package resolver

import (
	"fast-dns-server/internal/config"
	"github.com/bluele/gcache"
	"github.com/miekg/dns"
	"sync"
	"time"
)

type DnsServerInst struct {
	Id              int64
	tpcResolverInst dns.Server
	udpResolverInst dns.Server
	client          *dns.Client
	cache           gcache.Cache // 使用 gcache 来缓存 DNS 查询结果
	cacheMutex      sync.Mutex
}

// DnsServerStats DNS server performance stats
type DnsServerStats struct {
	SuccessCount    int
	FailureCount    int
	AvgResponseTime time.Duration
}

// DnsServer DNS server selection (weighted random selection)
type DnsServer struct {
	Address string
	Stats   DnsServerStats
}

var dnsServersStats = make(map[string]*DnsServer)

func NewDnsServerInst(id int64, addr string) *DnsServerInst {
	return &DnsServerInst{
		Id: id,
		tpcResolverInst: dns.Server{
			Addr: addr,
			Net:  "tcp",
		},
		udpResolverInst: dns.Server{
			Addr: addr,
			Net:  "udp",
		},
		cache: gcache.New(int(config.RootCfg.Details.Config.CacheSize)).
			Expiration(time.Second * time.Duration(1000)).
			ARC(). // 使用 ARC 算法
			Build(),
		client: new(dns.Client),
	}
}
