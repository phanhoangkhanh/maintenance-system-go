package authenticate

import (
	"fmt"
	models "maintenance-system-go/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


func (a *Authenticate) HandleLoginForm(loginRequest LoginRequest, c *gin.Context) (models.User, error) {
	// Implement the login logic here
	// fmt.Println(loginRequest)
	query := models.Query{
		Where: []models.WhereClause{
			{
				Key: "name",
				Compare: "ILIKE",
				Value: "%"+loginRequest.Name+"%",
			},
		},
		OrderBy: "created_at DESC , name ASC",
		Page: 1,
		PerPage: 2,
	}
	users, err := a.GetUser(query, c)
	if err != nil {
		return models.User{}, err
	}
	// log.Fatalf("users: %+v\n", users)


	if len(users) == 0 {
		return models.User{}, fmt.Errorf("user not found")
	}
	var userFound models.User
	for _, user := range users {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
		if err != nil {
		    continue
		}
		// Successful login
		userFound = user
		break
	}
	if userFound.ID == "" {
		return models.User{}, fmt.Errorf("Not Found User")
	}

	return userFound , nil
}