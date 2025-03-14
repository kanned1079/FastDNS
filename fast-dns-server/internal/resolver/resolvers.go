package resolver

import (
	"github.com/miekg/dns"
	"log"
)

func (this *DnsServerInst) RunDnsResolver() {
	dns.HandleFunc(".", this.handleDNSRequest) // 监听所有域名的请求
	log.Println("Starting dns server")
	if err := this.resolverInst.ListenAndServe(); err != nil {
		log.Print(err)
	}
}
