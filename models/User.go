package models

import "time"

//json : response + form: input request

//this Struct prepesent a user in Database, also the query params in request
type User struct {
	ID string `json:"id" gorm:"column:id;primaryKey" form:"id"`
	Name string `json:"name" gorm:"column:name;uniqueIndex" form:"name"`
	Email *string `json:"email" gorm:"column:email" form:"email"` // *string allow null value
	MobilePhone string `json:"mobile_phone" gorm:"column:mobile_phone" form:"mobile_phone"`
	Role string `json:"role" gorm:"column:role" form:"role"`
	Status string `json:"status" gorm:"column:status" form:"status"`
	Password string `json:"-" gorm:"column:password" form:"password"`  // not response password to client
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at" form:"created_at"` // CreatedAt is auto create current timestamp
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at" form:"updated_at"` // UpdatedAt is auto update current timestamp
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at" form:"deleted_at"` 
	Page int `json:"-" form:"page"` //no response but got from query URL request
	PerPage int `json:"-" form:"per_page"` //no response but got from query URL request
}