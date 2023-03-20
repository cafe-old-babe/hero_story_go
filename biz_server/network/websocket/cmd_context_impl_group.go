package websocket

import "google.golang.org/protobuf/reflect/protoreflect"

type cmdContextImplGroup struct {
	innerMap map[int32]*CmdContextImpl
}

var cmdContextImplGroupInstance = &cmdContextImplGroup{
	innerMap: make(map[int32]*CmdContextImpl),
}

func GetCmdContextImplGroup() *cmdContextImplGroup {
	return cmdContextImplGroupInstance
}

func (group *cmdContextImplGroup) Add(ctx *CmdContextImpl) {
	if ctx == nil {
		return
	}
	group.innerMap[ctx.SessionId] = ctx
}

func (group *cmdContextImplGroup) RemoveBySessionId(sessionId int32) {
	if sessionId <= 0 {
		return
	}
	delete(group.innerMap, sessionId)
}

func (group *cmdContextImplGroup) Broadcast(msgObj protoreflect.ProtoMessage) {
	if nil == msgObj {
		return
	}
	for _, ctx := range group.innerMap {
		if ctx != nil {
			ctx.Write(msgObj)
		}
	}
}
