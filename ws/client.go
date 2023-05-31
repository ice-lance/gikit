package ws

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/ice-lance/gikit/log"
)

type WSClient struct {
	IP        string
	Port      int
	WSChannel string
	WSID      string
	conn      *websocket.Conn
	codec     ICodec
	handler   IHandler
}

func NewWSClient(ip string, port int, wschannel string, wsid string, codec ICodec, handler IHandler) *WSClient {
	return &WSClient{
		IP:        ip,
		Port:      port,
		WSChannel: wschannel,
		WSID:      wsid,
		conn:      nil,
		codec:     codec,
		handler:   handler,
	}
}

func (client *WSClient) Connect() error {

	address := fmt.Sprintf("ws://%s:%d/ws/%s/%s", client.IP, client.Port, client.WSChannel, client.WSID)
	var err error
	// 连接WebSocket服务器
	client.conn, _, err = websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		return err
	}
	go func() {
		defer recover()
		for {
			// 读取消息
			messageType, p, err := client.conn.ReadMessage()
			if err != nil {
				log.Error("%s\n", err)
				continue
			}
			if messageType == websocket.BinaryMessage {
				client.handler.Router(client.codec.Encode(p))
			}
		}
	}()

	return nil

}

func (client *WSClient) Send(msg_t byte, v any) error {
	// 发送消息
	err := client.conn.WriteMessage(websocket.BinaryMessage, client.codec.Decode(int(msg_t), v))
	if err != nil {
		return err
	}
	return nil
}
