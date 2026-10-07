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
	users, err := s.Repo.GetUser(query, c.Request.Context())
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
	sessionID, err := s.Redis.Create(c.Request.Context(), userFound, sessionTTL)
	if err != nil {
		log.Printf("Redis error: %v\n", err)
		return models.User{}, err, http.StatusInternalServerError
	}

	// Only send the session cookie after Redis has stored the session successfully.
	c.SetCookie("session_id", sessionID, int(sessionTTL.Seconds()), "/", "", c.Request.TLS != nil, true)

	return userFound, nil, http.StatusOK
}

func (s *Service) GetListUser(c *gin.Context) ([]models.User, error, int) {
	var userQuery models.User
	if err := c.ShouldBindQuery(&userQuery); err != nil {
		return nil, err, http.StatusBadRequest
	}
	query := helper.GenerateWhereStruct(&userQuery)
	query.OrderBy = "created_at DESC , name ASC"
	query.Page = userQuery.Page
	query.PerPage = userQuery.PerPage
	fmt.Printf("QUERY JSON %v\n", query)
	
	users, err := s.Repo.GetUser(query, c.Request.Context())
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
	existingUsers, err := s.Repo.GetUser(query, c.Request.Context())
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
	existingUsers, err = s.Repo.GetUser(query, c.Request.Context())
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
	newUser := request.FromRequestToNewUser(string(hashedPassword))
	//Create User
	createdUser, err := s.Repo.CreateUser(newUser, c.Request.Context())
	if err != nil {
		return models.User{}, err, http.StatusInternalServerError
	}

	return *createdUser, nil, http.StatusOK
}

func (s *Service) UpdateUser(request *CreateOrUpdateUserRequest, c *gin.Context) (models.User, error, int) {
	//Check if user exists
	query := models.Query{
		Where: []models.WhereClause{
			{
				Key:     "id",
				Compare: "=",
				Value:   request.ID,
			},
		},
	}
	existingUsers, err := s.Repo.GetUser(query, c.Request.Context())
	if err != nil {
		return models.User{}, err, http.StatusInternalServerError
	}
	if len(existingUsers) == 0 {
		return models.User{}, fmt.Errorf("user not found"), http.StatusNotFound
	}
	userToUpdate := existingUsers[0] 

	//Not this user request 
	userActive, err := helper.GetUserActive(c)
	if err != nil {
		return models.User{}, err, http.StatusInternalServerError
	}
	if userToUpdate.ID == userActive.ID {
		return models.User{}, fmt.Errorf("cannot update the currently logged-in user"), http.StatusForbidden
	}

	//prepare password
	hashedPassword := ""
	if request.Password != "" {
		hashString, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return models.User{}, err, http.StatusInternalServerError
		}
		hashedPassword = string(hashString)
	}


	newUser := request.FromRequestToNewUser(hashedPassword)
	newUser.ID = userToUpdate.ID


	updatedUser, err := s.Repo.UpdateUser(newUser, c.Request.Context())
	if err != nil {
		return models.User{}, err, http.StatusInternalServerError
	}




	return *updatedUser, nil, http.StatusOK
}
