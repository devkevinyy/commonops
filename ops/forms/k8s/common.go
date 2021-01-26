package k8s_structs

import (
	"encoding/json"
	"fmt"
	"github.com/chujieyang/commonops/ops/untils/ws"
	"github.com/gorilla/websocket"
	"k8s.io/client-go/tools/remotecommand"
)

type Resp struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

func RespSuccess(data interface{}) (resp Resp) {
	return Resp{
		Code: 0,
		Msg: "success",
		Data: data,
	}
}

func RespError(code int, msg string) (resp Resp) {
	return Resp{
		Code: code,
		Msg: msg,
		Data: nil,
	}
}

type StreamHandler struct {
	WsConn *ws.WsConnection
	ResizeEvent chan remotecommand.TerminalSize
}

type xtermMessage struct {
	MsgType string `json:"type"` // 类型:resize客户端调整终端, input客户端输入
	Input string `json:"input"` // msgtype=input情况下使用
	Rows uint16 `json:"rows"` // msgtype=resize情况下使用
	Cols uint16 `json:"cols"` // msgtype=resize情况下使用
}

func (handler *StreamHandler) Next() (size *remotecommand.TerminalSize) {
	ret := <- handler.ResizeEvent
	size = &ret
	return
}

// executor 回调读取 web 端的输入
func (handler *StreamHandler) Read(p []byte) (size int, err error) {
	var (
		msg *ws.WsMessage
		xtermMsg xtermMessage
	)

	// 读web发来的输入
	if msg, err = handler.WsConn.WsRead(); err != nil {
		fmt.Println(err)
		return
	}

	// 解析客户端请求
	if err = json.Unmarshal(msg.Data, &xtermMsg); err != nil {
		fmt.Println(err)
		return
	}

	//web ssh调整了终端大小
	if xtermMsg.MsgType == "resize" {
		// 放到channel里，等remotecommand executor调用我们的Next取走
		handler.ResizeEvent <- remotecommand.TerminalSize{Width: xtermMsg.Cols, Height: xtermMsg.Rows}
	} else if xtermMsg.MsgType == "input" {	// web ssh终端输入了字符
		// copy到p数组中
		size = len(xtermMsg.Input)
		copy(p, xtermMsg.Input)
	}
	return
}

// executor 回调向 web 端输出
func (handler *StreamHandler) Write(p []byte) (size int, err error) {
	var (
		copyData []byte
	)
	copyData = make([]byte, len(p))
	copy(copyData, p)
	size = len(p)
	err = handler.WsConn.WsWrite(websocket.TextMessage, copyData)
	return
}