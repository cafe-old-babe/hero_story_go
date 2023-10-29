package websocket

import (
	"encoding/binary"
	"fmt"
	"github.com/gorilla/websocket"
	"hero_story/biz_server/handler"
	"hero_story/biz_server/msg"
	"hero_story/comm/log"
	"hero_story/comm/main_thread"
	"sync"
	"time"
)

type GatewayServerConn struct {
	GatewayServerId     int32
	WsConn              *websocket.Conn
	sendMsgQ            chan *msg.InternalServerMsg
	ctxMap              *sync.Map
	ctxMapLastClearTime int64
}

func (conn *GatewayServerConn) LoopSendMsg() {
	conn.sendMsgQ = make(chan *msg.InternalServerMsg, 64)
	go func() {
		for {

			msgObj := <-conn.sendMsgQ
			if nil == msgObj {
				continue
			}
			array := msgObj.ToByteArray()
			if err := conn.WsConn.WriteMessage(websocket.BinaryMessage, array); nil != err {
				log.Error("gateway_server send err: %+v", err)
			}
		}

	}()
}

func (conn *GatewayServerConn) LoopReadMsg() {
	if nil == conn.WsConn {
		return
	}

	conn.ctxMap = &sync.Map{}

	for {
		_, msgData, err := conn.WsConn.ReadMessage()
		if err != nil {
			log.Error("websocket readMessage error: %v ", err)
			break
		}
		log.Info("msgData: %v", msgData)

		func() {
			defer func() {
				if e := recover(); e != nil {
					log.Error("解析消息出错：%+v", e)
				}
			}()
			innerMsg := &msg.InternalServerMsg{}
			innerMsg.FromByteArray(msgData)
			realMsgData := innerMsg.MsgData

			msgCode := binary.BigEndian.Uint16(realMsgData[2:4])
			message, err := msg.Decode(realMsgData[4:], int16(msgCode))
			if err != nil {
				log.Error("message message msgCode: %d, err: %+v", msgCode, err)
				return
			}
			log.Info("收到客户端消息,remoteSessionId: %d, userId: %d, msgCode: %d, message Name: %v",
				innerMsg.SessionId, innerMsg.UserId, msgCode, message.Descriptor().Name())

			cmdHandlerFunc := handler.CreateCmdHandler(msgCode)
			if cmdHandlerFunc == nil {
				log.Error("没有查询到指令处理函数,msgCode: %d", msgCode)
				return
			}
			// 获取唯一的会话ID
			sessionUid := fmt.Sprintf("%d_%d", innerMsg.GatewayServerId, innerMsg.SessionId)
			ctx, _ := conn.ctxMap.LoadOrStore(sessionUid, &innerCmdContextImpl{
				gatewayServerId:   innerMsg.GatewayServerId,
				remoteSessionId:   innerMsg.SessionId,
				userId:            innerMsg.UserId,
				GatewayServerConn: conn,
			})
			if nil == ctx {
				log.Error("loadOrStore after ctx still nil")
				return
			}

			impl := ctx.(*innerCmdContextImpl)
			impl.lastActiveTime = time.Now().UnixMilli()
			main_thread.Process(func() {
				cmdHandlerFunc(impl, message)
			})
			// 判断ctxMap里有没有长时间没有发送消息的用户
			// 如果有,就删除掉
			conn.clearCtxMap()

		}()
	}
	//handler.OnUserQuitHandler(ctx)
}

// clearCtxMap
func (conn *GatewayServerConn) clearCtxMap() {
	nowTime := time.Now().UnixMilli()
	timeout := int64(2 * time.Second)
	if nowTime-conn.ctxMapLastClearTime < timeout {
		// 如果上次清除时间未超过2分钟,跳过
		return
	}
	conn.ctxMapLastClearTime = nowTime

	deleteUidSlice := make([]interface{}, 64)
	conn.ctxMap.Range(func(key, val any) bool {
		if nil == key || nil == val {
			return true
		}
		curCtx := val.(*innerCmdContextImpl)
		if nowTime-curCtx.lastActiveTime < timeout {
			return true
		}
		deleteUidSlice = append(deleteUidSlice, key)
		return true
	})
	for _, sessionId := range deleteUidSlice {
		if nil == sessionId {
			continue
		}
		// 删除
		conn.ctxMap.Delete(sessionId)
	}

}
