package base

type AsyncBizResult struct {
	returnObj interface{}
	// 1. 如果不在continueWith上调用的话,会有多线程问题读写的问题,完成函数跨线程(携程)了,2.会有不被调用的风险,赋值与调用错位,3.多次调用
	completeFunc func()
}

func (r *AsyncBizResult) GetReturnObj() interface{} {
	return r.returnObj
}

func (r *AsyncBizResult) SetReturnObj(val interface{}) {
	r.returnObj = val
}
func (r AsyncBizResult) OnComplete(val func()) {
	r.completeFunc = val
}
func (r *AsyncBizResult) DoComplete() {
	if nil != r.completeFunc {
		r.completeFunc()
	}
}
