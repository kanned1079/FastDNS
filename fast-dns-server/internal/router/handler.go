package router

import (
	"fast-dns-server/internal/config"
	"fast-dns-server/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
)

func (this *ApiInstance) HandleFetchConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"config": config.RootCfg.Details,
		"server": map[string]any{
			"os":      runtime.GOOS,
			"arch":    runtime.GOARCH,
			"routine": runtime.NumGoroutine(),
		},
		"cache": map[string]any{
			"cache_rate": utils.GlobeUtils.GetCacheHitRatio(),
			"cache_miss": utils.GlobeUtils.CacheMisses,
			"cache_hint": utils.GlobeUtils.CacheHits,
		},
	})
}
