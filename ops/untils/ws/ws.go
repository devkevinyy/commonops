package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func (r *http.Request) bool {
		return true
	},
}

type WsMessage struct {
	MessageType int
	Data []byte
}

type WsConnection struct {
	wsSocket *websocket.Conn // 底层websocket
	inChan chan *WsMessage	// 读取队列
	outChan chan *WsMessage // 发送队列

	mutex sync.Mutex	// 避免重复关闭管道
	isClosed bool
	closeChan chan byte  // 关闭通知
}

func (wsConn *WsConnection) InitLangEnv() {
	dataMap := map[string]string{
		"type": "input",
		"input": "export LANG=en_US.UTF-8 \r\n",
	}
	dataBytes, _ := json.Marshal(dataMap)
	msg := &WsMessage{
		1,
		dataBytes,
	}
	wsConn.inChan <- msg
}

func (wsConn *WsConnection) wsReadLoop() {
	var (
		msgType int
		data []byte
		msg *WsMessage
		err error
	)
	for {
		if msgType, data, err = wsConn.wsSocket.ReadMessage(); err != nil {
			goto ERROR
		}

		dataMap := map[string]string{
			"type": "input",
			"input": string(data),
		}

		inputMap := map[string]interface{}{}
		err = json.Unmarshal(data, &inputMap)
		if err != nil {
			//fmt.Println(err.Error())
		}
		if inputMap["action"]=="connection_close" {
			fmt.Println("auto close")
			dataMap = map[string]string{
				"type": "input",
				"input": "exit \r\n",
			}
		}
		dataBytes, _ := json.Marshal(dataMap)
		msg = &WsMessage{
			msgType,
			dataBytes,
		}
		select {
		case wsConn.inChan <- msg:
		case <- wsConn.closeChan:
			goto CLOSED
		}
	}
ERROR:
	wsConn.WsClose()
CLOSED:
	wsConn.WsClose()
}

func (wsConn *WsConnection) wsWriteLoop() {
	var (
		msg *WsMessage
		err error
	)
	for {
		select {
		case msg = <- wsConn.outChan:
			if err = wsConn.wsSocket.WriteMessage(msg.MessageType, msg.Data); err != nil {
				goto ERROR
			}
		case <- wsConn.closeChan:
			goto CLOSED
		}
	}
ERROR:
	wsConn.WsClose()
CLOSED:
	wsConn.WsClose()
}

func InitWebsocket(resp http.ResponseWriter, req *http.Request) (wsConn *WsConnection, err error) {
	var (
		wsSocket *websocket.Conn
	)
	// 应答客户端告知升级连接为websocket
	if wsSocket, err = upGrader.Upgrade(resp, req, nil); err != nil {
		return
	}
	wsConn = &WsConnection{
		wsSocket: wsSocket,
		inChan: make(chan *WsMessage, 1000),
		outChan: make(chan *WsMessage, 1000),
		closeChan: make(chan byte),
		isClosed: false,
	}

	// 读协程
	go wsConn.wsReadLoop()
	// 写协程
	go wsConn.wsWriteLoop()

	return
}

func (wsConn *WsConnection) WsWrite(messageType int, data []byte) (err error) {
	select {
	case wsConn.outChan <- &WsMessage{messageType, data,}:
	case <- wsConn.closeChan:
		err = errors.New("websocket closed")
	}
	return
}

func (wsConn *WsConnection) WsRead() (msg *WsMessage, err error) {
	select {
	case msg = <- wsConn.inChan:
		return
	case <- wsConn.closeChan:
		err = errors.New("websocket closed")
	}
	return
}

func (wsConn *WsConnection) WsClose() {
	wsConn.wsSocket.Close()
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.closeChan)
	}
}
