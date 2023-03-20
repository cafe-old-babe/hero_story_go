package async_op

import (
	"hero_story/comm/log"
	"hero_story/comm/main_thread"
)

type worker struct {
	taskQ chan func()
}

func (w *worker) process(asyncOp *func(), continueWith *func()) {
	if *asyncOp == nil {
		return
	}
	if nil == w.taskQ {
		log.Error("worker.process: taskQ is nil")
		return
	}
	w.taskQ <- func() {
		(*asyncOp)()
		if continueWith != nil {
			main_thread.Process(*continueWith)
		}
	}
}

func (w *worker) loopExecTask() {
	if nil == w.taskQ {
		log.Error("worker.loopExecTask: taskQ is nil")
		return
	}
	for {
		task := <-w.taskQ
		if nil != task {
			task()
		}
	}
}
