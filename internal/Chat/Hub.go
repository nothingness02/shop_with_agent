package chat

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	ID       string
	UserType int
}
type Hub struct {
	mu         sync.RWMutex // 读写锁，防止并发崩溃
	Clients    map[string]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

// 升级协议为websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// 循环进行监听来创建链接和断开链接
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) SendTo(targetID string, msg []byte) bool {
	h.mu.RLock()
	client, ok := h.Clients[targetID]
	h.mu.RUnlock()
	if ok {
		select {
		case client.Send <- msg:
			return true
		default:
			return false
		}
	}
	return false
}
