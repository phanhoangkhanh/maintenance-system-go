package customer

import (
	"context"
	"maintenance-system-go/helper"
	"maintenance-system-go/models"

	"github.com/gin-gonic/gin"
)

type Service struct {
	Repo  *Repository
}


func (service *Service) CreateCustomer(request *CustomerRequest, ctx context.Context) (*models.Customer, error, int) {
	// Implement the logic to create a customer using the repository
	customer, err := service.Repo.Create(request, ctx)
	if err != nil {
		return nil, err, 500
	}
	return customer, nil, 200
}

func (service *Service) GetListCustomer(c *gin.Context) ([]models.Customer, int64, error, int) {
	var customerQuery models.Customer
	if err := c.ShouldBindQuery(&customerQuery); err != nil {
		return nil, 0, err, 400
	}
	query := helper.GenerateWhereStruct(&customerQuery)
	query.OrderBy = "created_at DESC"
	query.Page = customerQuery.Page
	query.PerPage = customerQuery.PerPage

	customers, total, err := service.Repo.GetCustomer(query, c.Request.Context())
	if err != nil {
		return nil, 0, err, 500
	}
	return customers, total, nil, 200
}