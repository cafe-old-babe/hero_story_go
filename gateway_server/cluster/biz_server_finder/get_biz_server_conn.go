package biz_server_finder

import (
	"github.com/gorilla/websocket"
	"hero_story/biz_server/msg"
	"hero_story/comm/log"
	"sync"
)

var (
	bizServerConn *websocket.Conn
	locker        = &sync.Mutex{}
)

// GetBizServerConn 获取游戏服连接
func GetBizServerConn() (*websocket.Conn, error) {
	if bizServerConn != nil {
		return bizServerConn, nil
	}
	locker.Lock()
	defer locker.Unlock()
	if bizServerConn != nil {
		return bizServerConn, nil
	}
	//创建游戏服务器连接
	newConn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:12345/websocket", nil)
	if err != nil {
		return nil, err
	}
	bizServerConn = newConn
	// region 循环读取游戏服发来的消息,转发给客户端
	go func() {
		for {
			//读取游戏服务器返回的数据
			msgType, msgData, err := bizServerConn.ReadMessage()
			if err != nil {
				log.Error("从服务器读取消息失败: %+v", err)
			}

			innerMsg := &msg.InternalServerMsg{}
			innerMsg.FromByteArray(msgData)
			log.Info("从游戏服务器读取消息: sessionId: %d, userId: %d, type: %v, data: %v",
				innerMsg.SessionId, innerMsg.UserId, msgType, msgData)

			//ctx.Conn 网关服务器到游戏客户端连接 todo
			if err = bizServerConn.WriteMessage(msgType, innerMsg.MsgData); nil != err {
				log.Error("网关服务器到发送消息失败: %+v", err)
			}
		}
	}()
	//endregion
	return bizServerConn, nil
}
