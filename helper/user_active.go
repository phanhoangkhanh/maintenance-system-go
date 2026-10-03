package helper

import (
	"encoding/json"
	"maintenance-system-go/models"

	"github.com/gin-gonic/gin"
)

func GetUserActive(c *gin.Context) (models.User, error) {
	userByte := c.MustGet("current_user").([]byte)
	var user models.User
	err := json.Unmarshal(userByte, &user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}