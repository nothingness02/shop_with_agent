package chat

import (
	"encoding/json"
	"fmt"

	"github.com/myproject/shop/pkg/middleware"
)

type ChatService struct {
	repo  *ChatRepository
	redis *middleware.RedisStore
	hub   *Hub
}

func NewChatService(repo *ChatRepository, redis *middleware.RedisStore, hub *Hub) *ChatService {
	return &ChatService{
		repo:  repo,
		redis: redis,
		hub:   hub,
	}
}

func (s *ChatService) SendMessage(shopID, userID uint, SendType int, content string) error {
	//生成对话ID
	sessionID := fmt.Sprintf("shop:%d_user:%d", shopID, userID)
	msg := Message{
		SessionID:  sessionID,
		FromID:     userID, // 注意：如果是商家发，这里要做相应调整逻辑
		ToID:       shopID,
		SenderType: SendType, // "shop" 或 "user"
		Content:    content,
	}
	if SendType == ShopSender {
		msg.FromID = shopID
		msg.ToID = userID
	}
	//持久化
	if err := s.repo.CreateMessage(&msg); err != nil {
		return err
	}
	//推送列表
	targetKey := ""
	if SendType == ShopSender {
		targetKey = fmt.Sprintf("user:%d", userID)
	} else {
		targetKey = fmt.Sprintf("shop:%d", shopID)
	}
	msgBytes, _ := json.Marshal(msg)
	// 尝试直接通过 WebSocket 推送
	// 如果用户在线，他会立即收到。
	// 如果用户【断开连接】，SendTo 返回 false，什么都不做。
	s.hub.SendTo(targetKey, msgBytes)
	return nil
}

func (s *ChatService) GetHistory(shopID, userID uint) ([]Message, error) {
	var msgs []Message
	sessionID := fmt.Sprintf("shop:%d_user:%d", shopID, userID)
	msgs, err := s.repo.ListMessageBySessionID(sessionID, 50, 0)
	// 按时间倒序拉取最近 50 条
	return msgs, err
}
