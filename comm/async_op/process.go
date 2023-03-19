package async_op

import "sync"

var workerArray = [2048]*worker{}

var workerLocker = &sync.Mutex{}

func Process(bindId int, asyncOp *func(), continueWith *func()) {
	if asyncOp == nil {
		return
	}
	w := getCurrentWorker(bindId)
	if w == nil {
		return
	}

	w.process(asyncOp, continueWith)
}

// 根据id获取worker
func getCurrentWorker(bindId int) (w *worker) {

	if bindId < 0 {
		bindId = -bindId
	}
	i := bindId & len(workerArray)
	w = workerArray[i]
	if w != nil {
		return
	}
	workerLocker.Lock()
	defer workerLocker.Unlock()
	if w = workerArray[i]; w != nil {
		return
	}
	w = &worker{taskQ: make(chan func(), 2048)}
	go w.loopExecTask()
	workerArray[i] = w
	return
}
