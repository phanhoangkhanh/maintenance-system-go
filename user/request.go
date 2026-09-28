package user

type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateOrUpdateUserRequest struct {	
	Name     string `form:"name" json:"name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required,noduplicateEmail"`
	Role     string `form:"role" json:"role" binding:"required,inRoleGroup"`
	Status   string `form:"status" json:"status" binding:"required,in(ok|progress|delete)"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required,noduplicateMobile"`
}


