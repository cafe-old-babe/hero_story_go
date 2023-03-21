package broadcaster

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"hero_story/biz_server/base"
)

var innerMap = make(map[int32]base.MyCmdContext)

func Broadcast(msgObj protoreflect.ProtoMessage) {
	if nil == msgObj {
		return
	}
	for _, ctx := range innerMap {
		if ctx != nil {
			ctx.Write(msgObj)
		}
	}
}

func AddCmdCtx(sessionId int32, ctx base.MyCmdContext) {
	if ctx == nil {
		return
	}
	innerMap[sessionId] = ctx
}

func RemoveCmdCtxBySessionId(sessionId int32) {
	if sessionId <= 0 {
		return
	}
	delete(innerMap, sessionId)
}
