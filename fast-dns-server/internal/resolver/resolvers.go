package resolver

import (
	"github.com/miekg/dns"
	"log"
)

func (this *DnsServerInst) RunDnsResolver() {
	dns.HandleFunc(".", this.handleDNSRequest) // 监听所有域名的请求
	log.Println("Starting dns server")
	go func() {
		log.Println("start dns server tcp")
		if err := this.tpcResolverInst.ListenAndServe(); err != nil {
			log.Print(err)
		}
	}()
	go func() {
		log.Println("start dns server udp")
		if err := this.udpResolverInst.ListenAndServe(); err != nil {
			log.Print(err)
		}
	}()
	select {}
}
