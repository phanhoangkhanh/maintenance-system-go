package vehicle

type CreateOrUpdateVehicleRequest struct {
	ID       string `form:"id" json:"id" binding:"required_if=Action update"`
	Action     string `form:"action" json:"action" binding:"required,oneof=create update"`

	Kind string `json:"kind" gorm:"column:kind" form:"kind" binding:"required_if=Action create,omitempty"`
	NumberVehicle string `json:"number_vehicle" gorm:"column:number_vehicle" form:"number_vehicle" binding:"required_if=Action create,omitempty"`
	Status string `json:"status" gorm:"column:status" form:"status" binding:"required_if=Action create,omitempty,oneof=ok broken waiting"`
	Note string `json:"note" gorm:"column:note" form:"note" binding:"required_if=Action create,omitempty"`
}

type AttachVehicleRequest struct {
	VehicleID  string `form:"vehicle_id" json:"vehicle_id" binding:"required"`
	OperatorID string `form:"operator_id" json:"operator_id" binding:"required"`
	DriverID   string `form:"driver_id" json:"driver_id" binding:"required"`
	Note       string `form:"note" json:"note"`
}