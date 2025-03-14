package main

import (
	"fast-dns-server/internal/config"
	"fast-dns-server/internal/exec"
	"fast-dns-server/internal/model"
	"fast-dns-server/internal/resolver"
	"fast-dns-server/internal/router"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	log.Print("Init")

	if err := config.RootCfg.ReadConfigFile("./config/config.yaml"); err != nil {
		log.Println("read config file err: ", err)
		return
	}

	log.Println(config.RootCfg)

}

func main() {
	exec.Run()
	App := model.App{
		Id:          1,
		ApiGateway:  router.NewApiInstance(1, gin.DebugMode),
		DnsResolver: resolver.NewDnsServerInst(1, config.RootCfg.Details.Management.DnsServerListenAddr, "udp"),
	}
	go App.ApiGateway.RunApiGateway()
	App.DnsResolver.RunDnsResolver()
}
