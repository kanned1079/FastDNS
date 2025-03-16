package router

import (
	"fast-dns-server/internal/config"
	"fast-dns-server/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (this *ApiInstance) HandleFetchConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"config": config.RootCfg.Details,
		"cache": map[string]any{
			"cache_rate": utils.GlobeUtils.GetCacheHitRatio(),
			"cache_miss": utils.GlobeUtils.CacheMisses,
			"cache_hint": utils.GlobeUtils.CacheHits,
		},
	})
}
