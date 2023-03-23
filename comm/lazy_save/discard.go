package lazy_save

import "hero_story/comm/log"

func Discard(lso lazySaveObj) {
	if lso == nil {
		return
	}
	log.Info("放弃延时保存, lsoId: %d", lso.GetLsoId())

	lsoMap.Delete(lso.GetLsoId())
}
