package base

import (
	"hero_story/comm/main_thread"
	"sync/atomic"
)

type AsyncBizResult struct {
	returnObj interface{}
	// 1. 如果不在continueWith上调用的话,会有多线程问题读写的问题,完成函数跨线程(携程)了
	// 解决方式: 将该方法放到主线程执行.
	//,2.会有不被调用的风险,赋值与调用错位
	// OnComplete 设置后,检查是否设置了返回结果,如果设置了,直接执行
	//3.多次调用
	//  解决方式:记录是否调用过,标志位,多线程执行时会有跨线程的问题,completeFuncHasAlreadyBeenCalled,(hasReturnObj, hasCompleteFunc)
	completeFunc func()
	//默认值0,没有被调用
	completeFuncHasAlreadyBeenCalled int32
	//返回值与回调函数仅能设置一次
	hasReturnObj    int32
	hasCompleteFunc int32
}

func (r *AsyncBizResult) GetReturnObj() interface{} {
	return r.returnObj
}

func (r *AsyncBizResult) SetReturnObj(val interface{}) {
	if atomic.CompareAndSwapInt32(&r.hasReturnObj, 0, 1) {
		r.returnObj = val
		r.doComplete()
	}
}
func (r *AsyncBizResult) OnComplete(val func()) {
	if atomic.CompareAndSwapInt32(&r.hasCompleteFunc, 0, 1) {
		r.completeFunc = val
		if 1 == r.hasReturnObj {
			r.doComplete()
		}
	}
}

//doComplete 显示调用改为隐式调用
func (r *AsyncBizResult) doComplete() {
	if nil == r.completeFunc {
		return
	}

	//通过CAS原语比较标记值
	if atomic.CompareAndSwapInt32(&r.completeFuncHasAlreadyBeenCalled, 0, 1) {
		main_thread.Process(r.completeFunc)
	}
}
