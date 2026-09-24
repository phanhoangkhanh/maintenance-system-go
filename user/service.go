package user

import (
	"fmt"
	"log"
	models "maintenance-system-go/models"
	"net/http"
	"time"

	"maintenance-system-go/database/redis"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repo  *Repository
	Redis *redis.RedisClient
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

	//Create session in Redis
	sessionTTL := 24 * time.Hour // Set the session TTL as needed
	sessionID, err := s.Redis.Create(c, userFound, sessionTTL)
	if err != nil {
		log.Printf("Redis error: %v\n", err)
		return models.User{}, err, http.StatusInternalServerError
	}

	// Only send the session cookie after Redis has stored the session successfully.
	c.SetCookie("session_id", sessionID, int(sessionTTL.Seconds()), "/", "", c.Request.TLS != nil, true)

	return userFound, nil, http.StatusOK
}
