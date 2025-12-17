package user

import "gorm.io/gorm"

const (
	RoleCustomer uint = 1
	RoleMerchant uint = 5
	RoleAdmin    uint = 10
)

type User struct {
	gorm.Model
	Username string `gorm:"size:100;not null;unique"` // 用户名
	Email    string `gorm:"size:100;not null;unique"` // 邮箱
	Password string `gorm:"size:255;not null"`        // 密码（哈希值）
	Role     uint   `gorm:"not null"`                 // 角色（admin, customer）
	UserImg  string `gorm:"size:500"`                 // 用户头像URL
}
