package router

import (
	"fast-dns-server/internal/config"
	"log"
)

func (this *ApiInstance) RunApiGateway() {
	// 提供静态文件服务，确保能够访问到 dist 下的文件
	//this.Instance.StaticFS("/static/dist", http.FS(assets.StaticContext))
	this.Instance.Use(this.ProtocolAllowance())

	//// 提供前端页面，假设你的 index.html 在 dist 目录下
	//this.Instance.GET("/", func(c *gin.Context) {
	//	data, err := assets.StaticContext.ReadFile("dist/index.html")
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{
	//			"error": "Unable to read the index.html file",
	//		})
	//		return
	//	}
	//	c.Data(http.StatusOK, "text/html", data)
	//})

	// 捕获所有前端路由（用于 SPA）
	//this.Instance.NoRoute(func(c *gin.Context) {
	//	data, err := assets.StaticContext.ReadFile("dist/index.html")
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{
	//			"error": "Unable to read the index.html file",
	//		})
	//		return
	//	}
	//	c.Data(http.StatusOK, "text/html", data)
	//})

	// API 端点
	this.Instance.GET("/api/statistics", this.HandleFetchConfig)

	// 启动服务
	if err := this.Instance.Run(config.RootCfg.Details.Management.BackendListenAddr); err != nil {
		log.Print("run management backend server err, please check your config: ")
	}
	log.Println("management backend server has started successfully.")
}
