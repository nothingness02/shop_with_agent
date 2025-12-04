package user

import "github.com/myproject/shop/pkg/database"

type UserRepository struct {
	Database *database.Database
}

func NewRepository(db *database.Database) *UserRepository {
	return &UserRepository{Database: db}
}

func (r *UserRepository) GetUserByID(id uint) (*User, error) {
	var u User
	if err := r.Database.DB.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetUserByName(username string) (*User, error) {
	var u User
	if err := r.Database.DB.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) CreateUser(u *User) error {
	if err := r.Database.DB.Create(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUserByID(id uint) error {
	if err := r.Database.DB.Delete(&User{}, id).Error; err != nil {
		return err
	}
	return nil
}

//等待完善
func (r *UserRepository) UpdateUser(u *User) error {
	if err := r.Database.DB.Save(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) ListUsers(limit, offset int) ([]User, error) {
	var users []User
	if err := r.Database.DB.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
