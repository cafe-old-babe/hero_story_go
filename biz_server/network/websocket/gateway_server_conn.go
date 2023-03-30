package websocket

import (
	"encoding/binary"
	"github.com/gorilla/websocket"
	"hero_story/biz_server/handler"
	"hero_story/biz_server/msg"
	"hero_story/comm/log"
	"hero_story/comm/main_thread"
)

type GatewayServerConn struct {
	GatewayServerId int32
	WsConn          *websocket.Conn
	sendMsgQ        chan *msg.InternalServerMsg
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
			ctx := &innerCmdContextImpl{
				gatewayServerId:   innerMsg.GatewayServerId,
				remoteSessionId:   innerMsg.SessionId,
				userId:            innerMsg.UserId,
				GatewayServerConn: conn,
			}
			main_thread.Process(func() {
				cmdHandlerFunc(ctx, message)
			})

		}()
	}
	//handler.OnUserQuitHandler(ctx)
}
