package main_thread

import "hero_story/comm/log"

const maxQueueSize = 2048

var (
	mainQueue = make(chan func(), maxQueueSize)
)

func init() {
	//消费mainQueue队列并且执行
	go func() {
		for {
			task := <-mainQueue
			if nil == task {
				continue
			}
			func() {
				defer func() {
					if err := recover(); err != nil {
						log.Error("发生异常, %+v", err)
					}
				}()
				task()
			}()

		}
	}()
}

func Process(task func()) { //task --> interface Runnable {}
	if task == nil {
		return
	}
	mainQueue <- task
	/*if !started {
		startLocker.Lock()
		defer startLocker.Unlock()
		if !started {
			started = true
			go execute()
		}
	}*/
}

//按顺序执行
/*func execute() {
	for {
		task := <-mainQueue
		if nil != task {
			task()
		}
	}
}*/
