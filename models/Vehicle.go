package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Vehicle struct {
	ID string `json:"id" gorm:"column:id;primaryKey" form:"id"`
	Kind string `json:"kind" gorm:"column:kind" form:"kind"`
	NumberVehicle string `json:"number_vehicle" gorm:"column:number_vehicle" form:"number_vehicle"`
	Status string `json:"status" gorm:"column:status" form:"status"`
	Note string `json:"note" gorm:"column:note" form:"note"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at" form:"created_at"` // CreatedAt is auto create current timestamp
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at" form:"updated_at"` // UpdatedAt is auto update current timestamp
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at" form:"deleted_at"` 

	GetQueryParams
}

//Auto create ID
func (v *Vehicle) BeforeCreate(tx *gorm.DB) (err error) {
  v.ID = uuid.New().String()
  return
}