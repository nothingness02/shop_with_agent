package comment

import "time"

type Comment struct {
	// 使用单一自增主键更方便自引用与 GORM 操作
	ID uint `gorm:"primaryKey;autoIncrement"`

	UserID uint `gorm:"not null;index"` // 评论者
	ShopID uint `gorm:"not null;index"` // 所属店铺

	Content string `gorm:"type:text;not null"`

	// 自引用父评论，nullable
	ParentCommentID *uint    `gorm:"index"` // 指向父评论 ID，可为 nil
	Parent          *Comment `gorm:"foreignKey:ParentCommentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	// 子评论集合
	Children []Comment `gorm:"foreignKey:ParentCommentID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
