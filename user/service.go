package user

import (
	"fmt"
	models "maintenance-system-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repo *Repository
}

func (s *Service) HandleLoginForm(loginRequest LoginRequest, c *gin.Context) (models.User, error, int) {
	// Implement the login logic here
	// fmt.Println(loginRequest)
	query := models.Query{
		Where: []models.WhereClause{
			{
				Key:     "name",
				Compare: "ILIKE",
				Value:   "%" + loginRequest.Name + "%",
			},
		},
		OrderBy: "created_at DESC , name ASC",
	}
	users, err := s.Repo.GetUser(query, c)
	if err != nil {
		return models.User{}, err, http.StatusInternalServerError
	}
	// log.Fatalf("users: %+v\n", users)

	if len(users) == 0 {
		return models.User{}, fmt.Errorf("user not found"), http.StatusNotFound
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
		return models.User{}, fmt.Errorf("Not Found User"), http.StatusNotFound
	}

	return userFound, nil, http.StatusOK
}
