package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// For get query and model
type Plan struct {
	ID string `json:"id" gorm:"column:id;primaryKey" form:"id"`
	Customer string `json:"customer" gorm:"column:customer" form:"customer"`
	NumberVehicle string `json:"number_vehicle" gorm:"column:number_vehicle" form:"number_vehicle"`
	DriverCode string `json:"driver_code" gorm:"column:driver_code" form:"driver_code"`
	OperatorCode string `json:"operator_code" gorm:"column:operator_code" form:"operator_code"`
	Note string `json:"note" gorm:"column:note" form:"note"`
	Status string `json:"status" gorm:"column:status" form:"status"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at" form:"created_at"` // CreatedAt is auto create current timestamp
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at" form:"updated_at"` // UpdatedAt is auto update current timestamp

	GetQueryParams
}

func (Plan) TableName() string {
	return "plan"
}

//Auto create ID
func (v *Plan) BeforeCreate(tx *gorm.DB) (err error) {
  v.ID = uuid.New().String()
  return
}