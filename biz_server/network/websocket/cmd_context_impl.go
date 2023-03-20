package websocket

import (
	"encoding/binary"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/reflect/protoreflect"
	"hero_story/biz_server/handler"
	"hero_story/biz_server/msg"
	"hero_story/comm/log"
	"hero_story/comm/main_thread"
	"time"
)

const oneSecondMilli = 1000
const limitOnceByByteCount = 64 * 1024
const limitByPacketCountPerSecond = 16

// CmdContextImpl Implemented on MyCmdContext
type CmdContextImpl struct {
	userId       int64
	clientIpAddr string
	Conn         *websocket.Conn
	sendMsgQueue chan *protoreflect.ProtoMessage
	SessionId    int32
}

// BindUserId 绑定用户
func (ctx *CmdContextImpl) BindUserId(val int64) {
	ctx.userId = val
}

// GetUserId 获取用户
func (ctx *CmdContextImpl) GetUserId() int64 {
	return ctx.userId
}

// GetClientIpAddr 获取客户端ip
func (ctx *CmdContextImpl) GetClientIpAddr() string {
	return ctx.clientIpAddr
}

// Write 写入消息
func (ctx *CmdContextImpl) Write(msgObj protoreflect.ProtoMessage) {
	if nil == msgObj || ctx.Conn == nil || ctx.sendMsgQueue == nil {
		return
	}
	ctx.sendMsgQueue <- &msgObj

}

// SendError 错误消息
func (ctx *CmdContextImpl) SendError(errorCode int, errorInfo string) {

}

// Disconnect 断开连接
func (ctx *CmdContextImpl) Disconnect() {
	if ctx == nil || ctx.Conn == nil {

		return
	}
	_ = ctx.Conn.Close()
}

func (ctx *CmdContextImpl) LoopSendMsg() {
	if ctx.sendMsgQueue != nil {
		return
	}
	ctx.sendMsgQueue = make(chan *protoreflect.ProtoMessage, 1024)
	for {
		msgObj := <-ctx.sendMsgQueue
		if msgObj == nil {
			continue
		}
		byteArray, err := msg.Encode(msgObj)
		if err != nil {
			log.Error("[websocket] Encode msg error: %v", err)
			return
		}
		if err := ctx.Conn.WriteMessage(websocket.BinaryMessage, byteArray); err != nil {
			log.Error("[websocket] WriteMessage error: %v", err)
		}
	}
}

func (ctx *CmdContextImpl) LoopReadMsg() {
	if nil == ctx.Conn {
		return
	}
	ctx.Conn.SetReadLimit(limitOnceByByteCount)
	//limit the count of packet
	t0 := int64(0)
	counter := 0
	for {

		_, msgData, err := ctx.Conn.ReadMessage()
		if err != nil {
			log.Error("websocket readMessage error: %v ", err)
			break
		}
		log.Info("msgData: %v", msgData)

		t1 := time.Now().UnixMilli()
		if (t1 - t0) > oneSecondMilli {
			t0 = t1
			counter = 0
		}

		if counter >= limitByPacketCountPerSecond {
			log.Error("消息过于频繁, 拒绝处理, userId: %v, clientIp: %s", ctx.GetUserId(), ctx.GetClientIpAddr())
			continue
		}
		counter++

		msgCode := binary.BigEndian.Uint16(msgData[2:4])
		message, err := msg.Decode(msgData[4:], int16(msgCode))
		if err != nil {
			log.Error("message message msgCode: %d, err: %+v", msgCode, err)
			continue
		}
		log.Info("收到客户端消息,,msgCode: %d, message Name: %v", msgCode, message.Descriptor().Name())

		cmdHandlerFunc := handler.CreateCmdHandler(msgCode)
		if cmdHandlerFunc == nil {
			log.Error("没有查询到指令处理函数,msgCode: %d", msgCode)
			continue
		}
		main_thread.Process(func() {
			cmdHandlerFunc(ctx, message)
		})

	}
}
