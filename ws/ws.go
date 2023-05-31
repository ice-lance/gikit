package ws

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ice-lance/gikit/log"
)

// Manager 所有 websocket 信息
type Manager struct {
	Group                   map[string]map[string]*Client
	groupCount, clientCount uint
	Lock                    sync.Mutex
	Register, UnRegister    chan *Client
	Message                 chan *MessageData
	GroupMessage            chan *GroupMessageData
	BroadCastMessage        chan *BroadCastMessageData
	codec                   ICodec
	handler                 IHandler
}

// Client 单个 websocket 信息
type Client struct {
	Id, Group string
	Socket    *websocket.Conn
	Message   chan []byte
	Codec     ICodec
	Router    IHandler
}

// messageData 单个发送数据信息
type MessageData struct {
	Id, Group string
	Message   []byte
}

// groupMessageData 组广播数据信息
type GroupMessageData struct {
	Group   string
	Message []byte
}

// 广播发送数据信息
type BroadCastMessageData struct {
	Message []byte
}

// 读信息，从 websocket 连接直接读取数据
func (c *Client) Read() {
	defer func() {
		WebsocketManager.UnRegister <- c
		log.Info("client [%s] disconnect", c.Id)
		if err := c.Socket.Close(); err != nil {
			log.Info("client [%s] disconnect err: %s", c.Id, err)
		}
	}()

	for {
		messageType, message, err := c.Socket.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			break
		}

		log.Info("client [%s] receive message: %s", c.Id, string(message))
		// c.Message <- message
		c.Router.Router(c.Codec.Encode(message))

	}
}

// 写信息，从 channel 变量 Send 中读取数据写入 websocket 连接
func (c *Client) Write() {
	defer func() {
		log.Info("client [%s] disconnect", c.Id)
		if err := c.Socket.Close(); err != nil {
			log.Info("client [%s] disconnect err: %s", c.Id, err)
		}
	}()

	for message := range c.Message {
		log.Info("client [%s] write message: %s", c.Id, string(message))
		err := c.Socket.WriteMessage(websocket.BinaryMessage, message)
		if err != nil {
			log.Info("client [%s] writemessage err: %s", c.Id, err)
		}
	}
	_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
}

// 启动 websocket 管理器
func (manager *Manager) Start() {
	log.Info("websocket manage start")
	for {
		select {
		// 注册
		case client := <-manager.Register:
			log.Info("client [%s] connect", client.Id)
			log.Info("register client [%s] to group [%s]", client.Id, client.Group)

			manager.Lock.Lock()
			if manager.Group[client.Group] == nil {
				manager.Group[client.Group] = make(map[string]*Client)
				manager.groupCount += 1
			}
			manager.Group[client.Group][client.Id] = client
			manager.clientCount += 1
			manager.Lock.Unlock()

		// 注销
		case client := <-manager.UnRegister:
			log.Info("unregister client [%s] from group [%s]", client.Id, client.Group)
			manager.Lock.Lock()
			if _, ok := manager.Group[client.Group]; ok {
				if _, ok := manager.Group[client.Group][client.Id]; ok {
					close(client.Message)
					delete(manager.Group[client.Group], client.Id)
					manager.clientCount -= 1
					if len(manager.Group[client.Group]) == 0 {
						//log.Printf("delete empty group [%s]", client.Group)
						delete(manager.Group, client.Group)
						manager.groupCount -= 1
					}
				}
			}
			manager.Lock.Unlock()

			// 发送广播数据到某个组的 channel 变量 Send 中
			// case data := <-manager.BroadCastMessage:
			// 	if groupMap, ok := manager.Group[data.GroupId]; ok {
			// 		for _, conn := range groupMap {
			// 			conn.Message <- data.Data
			// 		}
			// 	}
		}
	}
}

// 处理单个 client 发送数据
func (manager *Manager) SendService() {
	for data := range manager.Message {
		if groupMap, ok := manager.Group[data.Group]; ok {
			if conn, ok := groupMap[data.Id]; ok {
				conn.Message <- data.Message
			}
		}
	}
}

// 处理 group 广播数据
func (manager *Manager) SendGroupService() {
	for data := range manager.GroupMessage {
		// 发送广播数据到某个组的 channel 变量 Send 中
		if groupMap, ok := manager.Group[data.Group]; ok {
			for _, conn := range groupMap {
				conn.Message <- data.Message
			}
		}
	}

}

// 处理广播数据
func (manager *Manager) SendAllService() {
	for data := range manager.BroadCastMessage {
		for _, v := range manager.Group {
			for _, conn := range v {
				conn.Message <- data.Message
			}
		}
	}
}

// 向指定的 client 发送数据
func (manager *Manager) Send(id string, group string, message []byte) {
	data := &MessageData{
		Id:      id,
		Group:   group,
		Message: message,
	}
	manager.Message <- data
}

// 向指定的 Group 广播
func (manager *Manager) SendGroup(group string, message []byte) {
	data := &GroupMessageData{
		Group:   group,
		Message: message,
	}
	manager.GroupMessage <- data
}

// 广播
func (manager *Manager) SendAll(message []byte) {
	data := &BroadCastMessageData{
		Message: message,
	}
	manager.BroadCastMessage <- data
}

// 注册
func (manager *Manager) RegisterClient(client *Client) {
	manager.Register <- client
}

// 注销
func (manager *Manager) UnRegisterClient(client *Client) {
	manager.UnRegister <- client
}

// 当前组个数
func (manager *Manager) LenGroup() uint {
	return manager.groupCount
}

// 当前连接个数
func (manager *Manager) LenClient() uint {
	return manager.clientCount
}

// 获取 wsManager 管理器信息
func (manager *Manager) Info() map[string]interface{} {
	managerInfo := make(map[string]interface{})
	managerInfo["groupLen"] = manager.LenGroup()
	managerInfo["clientLen"] = manager.LenClient()
	managerInfo["chanRegisterLen"] = len(manager.Register)
	managerInfo["chanUnregisterLen"] = len(manager.UnRegister)
	managerInfo["chanMessageLen"] = len(manager.Message)
	managerInfo["chanGroupMessageLen"] = len(manager.GroupMessage)
	managerInfo["chanBroadCastMessageLen"] = len(manager.BroadCastMessage)
	return managerInfo
}

var WebsocketManager *Manager

// 初始化 wsManager 管理器
func New(c ICodec, h IHandler) *Manager {
	return &Manager{
		Group:            make(map[string]map[string]*Client),
		groupCount:       0,
		clientCount:      0,
		Lock:             sync.Mutex{},
		Register:         make(chan *Client, 128),
		UnRegister:       make(chan *Client, 128),
		Message:          make(chan *MessageData, 128),
		GroupMessage:     make(chan *GroupMessageData, 128),
		BroadCastMessage: make(chan *BroadCastMessageData, 128),
		codec:            c,
		handler:          h,
	}
}

// gin 处理 websocket handler
func (manager *Manager) WsClient(ctx *gin.Context) {
	upGrader := websocket.Upgrader{
		// cross origin domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// 处理 Sec-WebSocket-Protocol Header
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}

	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Error("websocket connect error: %s", ctx.Param("channel"))
		return
	}

	client := &Client{
		Id:      ctx.Param("id"),
		Group:   ctx.Param("channel"),
		Socket:  conn,
		Message: make(chan []byte, 2048),
		Codec:   manager.codec,
		Router:  manager.handler,
	}

	manager.RegisterClient(client)
	go client.Read()
	go client.Write()
	time.Sleep(time.Second * 2)
	// 发送欢迎数据
	manager.Send(client.Id, client.Group, []byte("Send message ----"+time.Now().Format("2006-01-02 15:04:05")))
}

// 测试组广播
func TestSendGroup() {
	for {
		time.Sleep(time.Second * 20)
		WebsocketManager.SendGroup("leffss", []byte("SendGroup message ----"+time.Now().Format("2006-01-02 15:04:05")))
	}
}

// 测试广播
func TestSendAll() {
	for {
		time.Sleep(time.Second * 25)
		WebsocketManager.SendAll([]byte("SendAll message ----" + time.Now().Format("2006-01-02 15:04:05")))
		fmt.Println(WebsocketManager.Info())
	}
}
