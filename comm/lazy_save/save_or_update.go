package lazy_save

import (
	"hero_story/comm/log"
	"sync"
	"time"
)

var lsoMap = &sync.Map{}

func init() {
	startSave()
}

func SaveOrUpdate(lso lazySaveObj) {
	if lso == nil {
		return
	}
	log.Info("记录延时保存对象, lsoId = %s", lso.GetLsoId())
	nowTime := time.Now().UnixMilli()
	existRecord, _ := lsoMap.Load(lso.GetLsoId())
	if existRecord != nil {
		existRecord.(*lazySaveRecord).setLastUpdateTime(nowTime)
		return
	}

	lsoMap.Store(lso.GetLsoId(), &lazySaveRecord{
		lsoRef:         &lso,
		lastUpdateTime: nowTime,
	})
}

func startSave() {
	/*
		此处有问题
		1. userdao.SaveOrUpdate(lso.(*userdata.User)) 这步是同步操作,有缺陷
		3. 只需要关心什么时候存就行,还需要知道怎么存,不易扩展,如果后续添加新的需要延时保存的对象,还需要改造,comm是框架代码,不应该掺加业务代码
	*/
	go func() {
		for {
			time.Sleep(time.Second)
			nowTime := time.Now().UnixMilli()
			deleteArray := make([]string, 1024)
			lsoMap.Range(func(key, value interface{}) bool {
				if value == nil {
					return true
				}
				if record, ok := value.(*lazySaveRecord); ok {
					if (nowTime - record.getLastUpdateTime()) < (20 * 1000) {
						return true
					}
					log.Info("执行延时保存对象, lsoId = %s", (*record.lsoRef).GetLsoId())
					//userdao.SaveOrUpdate(lso.(*userdata.User))
					(*record.lsoRef).SaveOrUpdate()
					deleteArray = append(deleteArray, key.(string))
				}
				return true
			})
			for _, key := range deleteArray {
				lsoMap.Delete(key)
			}
		}
	}()
}
