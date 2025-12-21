package chat

import (
	"time"
)

const (
	ShopSender int = 5
	UserSender int = 1
)

type Message struct {
	ID         uint      `gorm:"primaryKey"`
	SessionID  string    `gorm:"index;size:50"` // 唯一标识一组对话，例如 "shop:1_user:101" (小ID在前)
	FromID     uint      `gorm:"index"`         // 发送者ID (用户ID或商家ID)
	ToID       uint      `gorm:"index"`         // 接收者ID
	SenderType int       `gorm:"size:10"`       // "user" 或 "shop"
	Content    string    `gorm:"type:text"`     // 消息内容
	IsRead     bool      `gorm:"default:false"` // 是否已读
	CreatedAt  time.Time `gorm:"index"`
}
type Conversation struct {
	ID          uint      `gorm:"primaryKey"`
	ShopID      uint      `gorm:"index"`     // 商家ID
	UserID      uint      `gorm:"index"`     // 用户ID
	LastMessage string    `gorm:"size:255"`  // 最后一条消息预览
	UnreadCount int       `gorm:"default:0"` // 商家未读数
	UpdatedAt   time.Time `gorm:"index"`     // 最后聊天时间 (用于排序)
}
