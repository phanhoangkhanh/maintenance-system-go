package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VehicleUser struct {
	ID         string `json:"id" gorm:"column:id;primaryKey" form:"id"`
	VehicleID  string `json:"vehicle_id" gorm:"column:vehicle_id" form:"vehicle_id"`
	OperatorID string `json:"operator_id" gorm:"column:operator_id" form:"operator_id"`
	DriverID   string `json:"driver_id" gorm:"column:driver_id" form:"driver_id"`

	Note      string     `json:"note" gorm:"column:note" form:"note"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at" form:"created_at"` // CreatedAt is auto create current timestamp
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at" form:"updated_at"` // UpdatedAt is auto update current timestamp

}

func (VehicleUser) TableName() string {
	return "vehicle_user"
}

// Auto create ID
func (v *VehicleUser) BeforeCreate(tx *gorm.DB) (err error) {
	v.ID = uuid.New().String()
	return
}
