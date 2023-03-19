package async_op

var workerArray = [2048]*worker{}

func Process(bindId int, asyncOp *func(), continueWith *func()) {
	if asyncOp == nil {
		return
	}
	go workerArray[bindId].loopExecTask()
	workerArray[bindId].process(asyncOp, continueWith)
}
