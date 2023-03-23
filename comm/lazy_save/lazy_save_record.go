package lazy_save

type lazySaveRecord struct {
	lsoRef         *lazySaveObj
	lastUpdateTime int64
}

func (r *lazySaveRecord) getLastUpdateTime() int64 {
	return r.lastUpdateTime
}
func (r *lazySaveRecord) setLastUpdateTime(val int64) {
	r.lastUpdateTime = val
}
