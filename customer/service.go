package customer

import (
	"context"
	"maintenance-system-go/models"
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