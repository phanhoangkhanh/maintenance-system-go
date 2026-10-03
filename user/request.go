package user

import "maintenance-system-go/models"

type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateOrUpdateUserRequest struct {
	ID       string `form:"id" json:"id" binding:"required_if=Kind update"`
	Kind     string `form:"kind" json:"kind" binding:"required,oneof=create update"`
	Name     string `form:"name" json:"name" binding:"required_if=Kind create"`
	Password string `form:"password" json:"password" binding:"required_if=Kind create"`
	Email    string `form:"email" json:"email" binding:"required_if=Kind create,omitempty,email"`
	Role     string `form:"role" json:"role" binding:"required_if=Kind create,omitempty,inRoleGroup"`
	Status   string `form:"status" json:"status" binding:"required_if=Kind create,omitempty,oneof=ok progress delete"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required_if=Kind create,omitempty,number"`
}

func (r CreateOrUpdateUserRequest) FromRequestToNewUser(passwordHash string) *models.User {
	user := &models.User{
		Name:        r.Name,
		Email:       &r.Email,
		MobilePhone: r.Mobile,
		Role:        r.Role,
		Status:      r.Status,
	}
	if passwordHash != "" {
		user.Password = passwordHash
	}
	return user
}

// FromRequestToUserUpdate leaves empty fields at their zero value so GORM
// Updates skips them, including the nullable email field.
func (r CreateOrUpdateUserRequest) FromRequestToUserUpdate(passwordHash string) *models.User {
	user := r.FromRequestToNewUser(passwordHash)
	user.ID = r.ID
	if r.Email == "" {
		user.Email = nil
	}
	return user
}
