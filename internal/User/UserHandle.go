package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/myproject/shop/pkg/utils"
)

type UserHandle struct {
	servive *UserService
}

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     uint   `json:"role" binding:"required,oneof=1 5 10"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
	UserImg  string `json:"user_img" binding:"omitempty,url"`
	Phone    string `json:"phonenums" binding:"omitempty,phone"`
}

func NewUserHandle(service *UserService) *UserHandle {
	return &UserHandle{servive: service}
}

func (h *UserHandle) GetUserByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.servive.Repo.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandle) GetUserByName(c *gin.Context) {
	username := c.Query("username")
	user, err := h.servive.Repo.GetUserByName(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandle) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.servive.RegisterUser(req.Username, req.Email, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//等待规范返回的数据
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandle) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.servive.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandle) UpdateUser(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Username == "" && req.Email == "" && req.Password == "" && req.UserImg == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.servive.Repo.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		user.Password = utils.HashPassword(req.Password)
	}
	if req.UserImg != "" {
		user.UserImg = req.UserImg
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if err := h.servive.Repo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
