package plan

type PlanRequest struct {
	Customer string `json:"customer" binding:"required,inCustomerList" form:"customer"`
	NumberVehicle string `json:"number_vehicle" binding:"required,inNumberVehicleList" form:"number_vehicle"`
	DriverCode string `json:"driver_code" binding:"required,inDriverList" form:"driver_code"`
	OperatorCode string `json:"operator_code"  form:"operator_code"`
	Note string `json:"note"  form:"note"`
	Status string `json:"status" binding:"required,oneof=planning collecting finished" form:"status"`
}