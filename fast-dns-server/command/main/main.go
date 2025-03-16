package main

import (
	"fast-dns-server/internal/config"
	"fast-dns-server/internal/logger"
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

	//log.Println(config.RootCfg)

	logger.MyLogger = logger.NewLogger("./app.log")

}

func main() {
	//exec.Run()
	App := model.App{
		Id:          1,
		ApiGateway:  router.NewApiInstance(1, gin.DebugMode),
		DnsResolver: resolver.NewDnsServerInst(1, config.RootCfg.Details.Management.DnsServerListenAddr, "udp"),
		//Logger:      logger.NewLogger("app.log", 300),
	}
	go App.ApiGateway.RunApiGateway()
	//go utils.GlobeUtils.ShowLogEveryInterval()
	go logger.MyLogger.StartCheckSizeInterval()
	go App.DnsResolver.RunDnsResolver()
	select {}
}
