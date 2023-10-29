package websocket

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"hero_story/biz_server/msg"
	"hero_story/comm/log"
)

type innerCmdContextImpl struct {
	gatewayServerId int32
	remoteSessionId int32
	userId          int64
	clientIpAddr    string
	lastActiveTime  int64
	*GatewayServerConn
}

// BindUserId 绑定用户
func (ctx *innerCmdContextImpl) BindUserId(val int64) {
	ctx.userId = val
}

// GetUserId 获取用户
func (ctx *innerCmdContextImpl) GetUserId() int64 {
	return ctx.userId
}

// GetClientIpAddr 获取客户端ip
func (ctx *innerCmdContextImpl) GetClientIpAddr() string {
	return ctx.clientIpAddr
}

// Write 写入消息
func (ctx *innerCmdContextImpl) Write(msgObj protoreflect.ProtoMessage) {
	if msgObj == nil {
		return
	}
	byteArray, err := msg.Encode(&msgObj)
	if nil != err {
		log.Error("encode msg fail, err:%v", err)
		return
	}
	innerMsg := &msg.InternalServerMsg{
		GatewayServerId: ctx.gatewayServerId,
		SessionId:       ctx.remoteSessionId,
		UserId:          ctx.userId,
		MsgData:         byteArray,
	}
	ctx.GatewayServerConn.sendMsgQ <- innerMsg
}

// SendError 错误消息
func (ctx *innerCmdContextImpl) SendError(errorCode int, errorInfo string) {

}

// Disconnect 断开连接
func (ctx *innerCmdContextImpl) Disconnect() {

}
