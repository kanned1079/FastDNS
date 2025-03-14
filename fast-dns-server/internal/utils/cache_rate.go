package utils

func (this *Utils) GetCacheHitRatio() float64 {
	total := this.CacheHits + this.CacheMisses
	if total == 0 {
		return 0
	}
	return float64(this.CacheHits) / float64(total)
}
