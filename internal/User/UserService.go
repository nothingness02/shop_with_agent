package user

import (
	"errors"

	"github.com/myproject/shop/pkg/utils"
)

type UserService struct {
	Repo *UserRepository
}

func NewService(repo *UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) RegisterUser(username, email, password string, role uint) (*User, error) {
	if username == "" || email == "" || password == "" {
		return nil, errors.New("the input cannot be empty")
	}
	// Check if user already exists
	existingUser, _ := s.Repo.GetUserByName(username)
	if existingUser != nil {
		return nil, errors.New("User already exists")
	}
	// 等待检查是否符合 命名 ，密码 ，邮箱的规范

	user := &User{
		Username: username,
		Email:    email,
		Password: utils.HashPassword(password),
		Role:     role,
	}
	if err := s.Repo.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) DeleteUser(id uint) error {
	return s.Repo.DeleteUserByID(id)
}
