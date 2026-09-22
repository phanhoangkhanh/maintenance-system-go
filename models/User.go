package models

import "time"


type User struct {
	ID string `json:"id" gorm:"column:id;primaryKey"`
	Name string `json:"name" gorm:"column:name;uniqueIndex"`
	Email *string `json:"email" gorm:"column:email"` // *string allow null value
	MobilePhone string `json:"mobile_phone" gorm:"column:mobile_phone"`
	Role string `json:"role" gorm:"column:role"`
	Status string `json:"status" gorm:"column:status"`
	Password string `json:"-" gorm:"column:password"`  // not response password to client
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at"` // CreatedAt is auto create current timestamp
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at"` // UpdatedAt is auto update current timestamp
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at"` 
}