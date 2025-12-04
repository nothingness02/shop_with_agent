package comment

import (
	"context"

	"github.com/myproject/shop/pkg/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentRepository struct {
	Database *database.Database
}

func NewRepository(db *database.Database) *CommentRepository {
	return &CommentRepository{Database: db}
}

func (r *CommentRepository) CreateComment(ctx context.Context, comment *Comment) error {
	if err := r.Database.DB.WithContext(ctx).Create(comment).Error; err != nil {
		return err
	}
	return nil
}

func (r *CommentRepository) ListByShop(ctx context.Context, shopID uint, parentID *uint, page, pageSize int, selectFields []string) ([]Comment, error) {
	var comments []Comment
	query := r.Database.DB.WithContext(ctx).Model(&Comment{}).Where("shop_id = ?", shopID)
	if parentID == nil {
		query = query.Where("parent_comment_id IS NULL")
	} else {
		query = query.Where("parent_comment_id = ?", *parentID)
	}
	if len(selectFields) > 0 {
		query = query.Select(selectFields)
	}
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepository) DeleteComment(ctx context.Context, userID, shopID uint, cascadeChildren bool) error {
	return r.Database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先收集将被删除的评论 ID（限定在 user_id & shop_id 范围内）
		var ids []uint
		if err := tx.Model(&Comment{}).
			Where("user_id = ? AND shop_id = ?", userID, shopID).
			Pluck("id", &ids).Error; err != nil {
			return err
		}
		if len(ids) == 0 {
			// 没有要删除的记录，直接返回 nil
			return nil
		}

		if !cascadeChildren {
			// 仅把属于将被删除的评论的子评论 parent_comment_id 置空
			if err := tx.Model(&Comment{}).
				Where("parent_comment_id IN ?", ids).
				Update("parent_comment_id", nil).Error; err != nil {
				return err
			}
		} else {
			// 级联删除子评论（递归或同层级），这里直接删除 parent_comment_id IN ids 的子评论
			if err := tx.Where("parent_comment_id IN ?", ids).Delete(&Comment{}).Error; err != nil {
				return err
			}
		}

		// 删除目标评论
		if err := tx.Where("id IN ?", ids).Delete(&Comment{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *CommentRepository) AddReply(ctx context.Context, parentID uint, reply *Comment) error {
	return r.Database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var parent Comment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&parent, "id=?", parentID).Error; err != nil {
			return err
		}
		reply.ParentCommentID = &parentID
		reply.ShopID = parent.ShopID
		if err := tx.Create(reply).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *CommentRepository) CountByCommentID(ctx context.Context, shopID uint, userID uint) (int64, error) {
	var count int64
	if err := r.Database.DB.WithContext(ctx).Model(&Comment{}).Where("shop_id = ? AND user_id = ?", shopID, userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByShop 返回某个店铺的评论总数
func (r *CommentRepository) CountByShop(ctx context.Context, shopID uint) (int64, error) {
	var count int64
	if err := r.Database.DB.WithContext(ctx).Model(&Comment{}).Where("shop_id = ?", shopID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetByID 获取单条评论，用于存在性校验及检索父评论
func (r *CommentRepository) GetByID(ctx context.Context, id uint) (*Comment, error) {
	var c Comment
	if err := r.Database.DB.WithContext(ctx).First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
