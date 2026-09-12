package models


type User struct {
	Name string `json:"name" gorm:"column:name"`
	Email string `json:"email" gorm:"column:email"`
	MobilePhone string `json:"mobile_phone" gorm:"column:mobile_phone"`
	Role string `json:"role" gorm:"column:role"`
	Status string `json:"status" gorm:"column:status"`
	Password string `json:"password" gorm:"column:password"`
	CreatedAt string `json:"created_at" gorm:"column:created_at"`
	UpdatedAt string `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt string `json:"deleted_at" gorm:"column:deleted_at"`
}