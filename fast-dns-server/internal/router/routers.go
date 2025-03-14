package router

import (
	"fast-dns-server/internal/config"
	"log"
)

func (this *ApiInstance) RunApiGateway() {
	this.Instance.GET("/api/config", this.HandleFetchConfig)

	if err := this.Instance.Run(config.RootCfg.Details.Management.BackendListenAddr); err != nil {
		log.Print("run management backend server err, please check your config: ")
	}
	log.Println("management backend server has started successfully.")
}
