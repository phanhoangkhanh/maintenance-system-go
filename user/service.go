package user

import (
	"fmt"
	"log"
	"maintenance-system-go/helper"
	models "maintenance-system-go/models"
	"net/http"
	"time"

	"maintenance-system-go/database/redis"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (s *Service) GetListUser(c *gin.Context) ([]models.User, error, int) {
	var user models.User
	if err := c.ShouldBindQuery(&user); err != nil {
		return nil, err, http.StatusBadRequest
	}
	query := helper.GenerateWhereStruct(&user)
	query.OrderBy = "created_at DESC , name ASC"
	fmt.Printf("QUERY JSON %v\n", query)
	query.Page = user.Page
	query.PerPage = user.PerPage
	
	users, err := s.Repo.GetUser(query, c)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	return users, nil, http.StatusOK
}



func (s *Service) CreateNewUser(request *CreateOrUpdateUserRequest, c *gin.Context) (models.User, error, int) {
	//Check unique name and mobile 
	query := models.Query{
		Where: []models.WhereClause{
			{
				Key:     "name",
				Compare: "=",
				Value:   request.Name,
			},
		},
	}
	existingUsers, err := s.Repo.GetUser(query, c)
	if err != nil {
		return models.User{}, err, http.StatusInternalServerError
	}
	if len(existingUsers) > 0 {
		return models.User{}, fmt.Errorf("user with the same name already exists"), http.StatusBadRequest
	}
	query = models.Query{
		Where: []models.WhereClause{
			{
				Key:     "mobile_phone",
				Compare: "=",
				Value:   request.Mobile,
			},
		},
	}
	existingUsers, err = s.Repo.GetUser(query, c)
	if err != nil {
		return models.User{}, err, http.StatusInternalServerError
	}
	if len(existingUsers) > 0 {
		return models.User{}, fmt.Errorf("user with the same mobile phone already exists"), http.StatusBadRequest
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err, http.StatusInternalServerError
	}
	newUser := request.ToNewUser(string(hashedPassword))
	id, err := uuid.NewRandom()
	if err != nil {
		return models.User{}, err, http.StatusInternalServerError
	}
	newUser.ID = id.String()
	//Create User
	createdUser, err := s.Repo.CreateUser(newUser, c)
	if err != nil {
		return models.User{}, err, http.StatusInternalServerError
	}

	return *createdUser, nil, http.StatusOK
}

func (s *Service) UpdateUser(request *CreateOrUpdateUserRequest, c *gin.Context) (models.User, error, int) {
	// Implement the update user logic here
	return models.User{}, nil, http.StatusOK
}
