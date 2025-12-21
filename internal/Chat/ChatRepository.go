package chat

import (
	"github.com/myproject/shop/pkg/database"
	"gorm.io/gorm"
)

type ChatRepository struct {
	Database *database.Database
}

func NewChatRepository(db *database.Database) *ChatRepository {
	return &ChatRepository{Database: db}
}

func (r *ChatRepository) ListMessageBySessionID(sessionID string, limit, offset int) ([]Message, error) {
	var msg []Message
	if err := r.Database.DB.Where("session_id = ?", sessionID).Order("created_at desc").Limit(limit).Offset(offset).Find(&msg).Error; err != nil {
		return nil, err
	}
	return msg, nil
}

func (r *ChatRepository) CreateMessage(msg *Message) error {
	return r.Database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(msg).Error; err != nil {
			return err
		}
		return nil
	})
}
