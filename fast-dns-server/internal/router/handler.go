package router

import (
	"fast-dns-server/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (this *ApiInstance) HandleFetchConfig(ctx *gin.Context) {
	// ...
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"config": map[string]any{
			"dns1": "223.5.5.5",
			"dns2": "223.5.5.5",
		},
		"cache_rate": utils.GlobeUtils.GetCacheHitRatio(),
	})
}
