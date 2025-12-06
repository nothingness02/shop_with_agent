package comment

import (
	"context"
	"errors"
)

type CommentService struct {
	repo *CommentRepository
}

func NewCommentService(repo *CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

// CreateComment 创建评论；如果提供 ParentCommentID，则作为回复处理
func (s *CommentService) CreateComment(ctx context.Context, c *Comment) error {
	if c == nil {
		return errors.New("comment is nil")
	}
	if c.ParentCommentID != nil {
		// 检查父评论存在且属于同一 shop
		parent, err := s.repo.GetByID(ctx, *c.ParentCommentID)
		if err != nil {
			return err
		}
		if parent.ShopID != c.ShopID {
			return errors.New("parent comment belongs to different shop")
		}
		return s.repo.AddReply(ctx, *c.ParentCommentID, c)
	}
	return s.repo.CreateComment(ctx, c)
}

// ListCommentsByShopWithCount 分页获取某店铺的评论（仅顶层或指定 parent），并返回总数
func (s *CommentService) ListCommentsByShopWithCount(ctx context.Context, shopID uint, parentID *uint, page, pageSize int, selectFields []string) ([]Comment, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	items, err := s.repo.ListByShop(ctx, shopID, parentID, page, pageSize, selectFields)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.CountByShop(ctx, shopID)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// DeleteComment 删除某用户在某店铺下的评论（可选择是否级联删除子评论）
func (s *CommentService) DeleteComment(ctx context.Context, userID, shopID uint, cascade bool) error {
	return s.repo.DeleteComment(ctx, userID, shopID, cascade)
}
