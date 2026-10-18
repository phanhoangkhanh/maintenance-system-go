package customer

type CustomerRequest struct {

	Name string `json:"name" binding:"required,uniqueName" form:"name"`
	Address   string    `json:"address" binding:"required,uniqueAddress" form:"address"`
	Note      string    `json:"note"  form:"note"`
	Status    string    `json:"status" binding:"required,oneof=operation hold" form:"status"`

	
}