package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {	
	ID string `json:"id" gorm:"column:id;primaryKey" form:"id"`
	Name string `json:"name" gorm:"column:name;uniqueIndex;not null" form:"name"`
	Address   string    `json:"address" gorm:"column:address" form:"address"`
	Note      string    `json:"note" gorm:"column:note" form:"note"`
	Status    string    `json:"status" gorm:"column:status" form:"status"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at" form:"created_at"` // CreatedAt is auto create current timestamp
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at" form:"updated_at"` // UpdatedAt is auto update current timestamp

	GetQueryParams
}

//Auto create ID
func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
  c.ID = uuid.New().String()
  return
}