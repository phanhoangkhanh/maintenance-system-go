package authenticate

import (
	"maintenance-system-go/config"
	"net/http"

	res "maintenance-system-go/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Authenticate struct {
	Database *gorm.DB
	Config *config.Config
}

func InitAuthenticate(db *gorm.DB, config *config.Config) *Authenticate {
	return &Authenticate{
		Database: db,
		Config:   config,
	}
}


func (a *Authenticate) Login(c *gin.Context) {
	var loginRequest LoginRequest
	err := c.ShouldBindJSON(&loginRequest)
	if err != nil {
		res.ResponseClient(c, http.StatusBadRequest, "Login Form not valid", err.Error())
		return
	} 
	
	user, err := a.HandleLoginForm(loginRequest, c)
	if err != nil {
		res.ResponseClient(c, http.StatusInternalServerError, "Login failed", err.Error())
		return
	}
	res.ResponseClient(c, 200, "Login successful", user)
	
}


