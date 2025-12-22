package cart

type CartService struct {
	repo *CartRepository
}

func NewCartService(repo *CartRepository) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) List(userID uint) ([]CartItem, error) {
	return s.repo.ListByUser(userID)
}

func (s *CartService) Add(userID, productID uint, quantity int) (*CartItem, error) {
	return s.repo.AddOrUpdate(userID, productID, quantity)
}

func (s *CartService) UpdateQuantity(userID, itemID uint, quantity int) error {
	return s.repo.UpdateQuantity(userID, itemID, quantity)
}

func (s *CartService) Delete(userID, itemID uint) error {
	return s.repo.DeleteItem(userID, itemID)
}

func (s *CartService) Clear(userID uint) error {
	return s.repo.Clear(userID)
}