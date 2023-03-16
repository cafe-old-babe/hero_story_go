package handler

import (
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/types/dynamicpb"
)

// CmdHandlerFunc 消息处理函数
//  websocket.conn 不足:1.分布式架构不适用;2.不支持限流,消息包很小,可能被攻击者攻击;3.channel无法存储业务数据
type CmdHandlerFunc func(conn *websocket.Conn, pbMsgObj *dynamicpb.Message) //java==> Function f; f.apply()
// key->cmdCode
var cmdHandlerMap = make(map[uint16]CmdHandlerFunc) //注册函数

// CreateCmdHandler 根据消息代号创建指令处理器
func CreateCmdHandler(msgCode uint16) CmdHandlerFunc {
	return cmdHandlerMap[msgCode]
}
