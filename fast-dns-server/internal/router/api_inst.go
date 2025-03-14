package router

import "github.com/gin-gonic/gin"

type ApiInstance struct {
	Id       int64
	RunMode  string
	Instance *gin.Engine
}

func NewApiInstance(id int64, runMode string) *ApiInstance {
	gin.SetMode(runMode)
	return &ApiInstance{
		Id:       id,
		RunMode:  runMode,
		Instance: gin.Default(),
	}
}
