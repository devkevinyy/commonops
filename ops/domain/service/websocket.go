package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chujieyang/commonops/ops/opslog"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsMessage struct {
	MessageType int
	Data        []byte
}

type WsConnection struct {
	wsSocket *websocket.Conn // 底层websocket
	inChan   chan *WsMessage // 读取队列
	outChan  chan *WsMessage // 发送队列

	mutex     sync.Mutex // 避免重复关闭管道
	isClosed  bool
	closeChan chan byte // 关闭通知
}

func (wsConn *WsConnection) InitLangEnv() {
	dataMap := map[string]string{
		"type":  "input",
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
		data    []byte
		msg     *WsMessage
		err     error
	)
	for {
		if msgType, data, err = wsConn.wsSocket.ReadMessage(); err != nil {
			goto ERROR
		}

		dataMap := map[string]string{
			"type":  "input",
			"input": string(data),
		}

		inputMap := map[string]interface{}{}
		err = json.Unmarshal(data, &inputMap)
		if err != nil {
			opslog.Error().Println(err)
		}
		if inputMap["action"] == "connection_close" {
			opslog.Info().Println("auto close")
			dataMap = map[string]string{
				"type":  "input",
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
		case <-wsConn.closeChan:
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
		case msg = <-wsConn.outChan:
			if err = wsConn.wsSocket.WriteMessage(msg.MessageType, msg.Data); err != nil {
				goto ERROR
			}
		case <-wsConn.closeChan:
			goto CLOSED
		}
	}
ERROR:
	wsConn.WsClose()
CLOSED:
	wsConn.WsClose()
}

func NewWebSocketService(resp http.ResponseWriter, req *http.Request) (wsConn *WsConnection, err error) {
	var (
		wsSocket *websocket.Conn
	)
	// 应答客户端告知升级连接为websocket
	if wsSocket, err = upGrader.Upgrade(resp, req, nil); err != nil {
		return
	}
	wsConn = &WsConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan *WsMessage, 1000),
		outChan:   make(chan *WsMessage, 1000),
		closeChan: make(chan byte),
		isClosed:  false,
	}

	// 读协程
	go wsConn.wsReadLoop()
	// 写协程
	go wsConn.wsWriteLoop()

	return
}

func (wsConn *WsConnection) WsWrite(messageType int, data []byte) (err error) {
	select {
	case wsConn.outChan <- &WsMessage{messageType, data}:
	case <-wsConn.closeChan:
		err = errors.New("websocket closed")
	}
	return
}

func (wsConn *WsConnection) WsRead() (msg *WsMessage, err error) {
	select {
	case msg = <-wsConn.inChan:
		return
	case <-wsConn.closeChan:
		err = errors.New("websocket closed")
	}
	return
}

func (wsConn *WsConnection) WsClose() {
	defer func() {
		if err := wsConn.wsSocket.Close(); err != nil {
			opslog.Error().Println(err)
		}
	}()
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.closeChan)
	}
}

func (wsConn *WsConnection) GetSocket() *websocket.Conn {
	return wsConn.wsSocket
}

func NewSshClient(host string, port int, user string, passwd string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 5,
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	config.Auth = []ssh.AuthMethod{ssh.Password(passwd)}
	addr := fmt.Sprintf("%s:%d", host, port)
	c, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	return c, nil
}


type StreamHandler struct {
	WsConn *WsConnection
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
		msg *WsMessage
		xtermMsg xtermMessage
	)

	// 读web发来的输入
	if msg, err = handler.WsConn.WsRead(); err != nil {
		opslog.Error().Println(err)
		return
	}

	// 解析客户端请求
	if err = json.Unmarshal(msg.Data, &xtermMsg); err != nil {
		opslog.Error().Println(err)
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