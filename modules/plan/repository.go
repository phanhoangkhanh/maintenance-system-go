package plan

import (
	"context"
	"maintenance-system-go/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func InitRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}
func (repo *Repository) CreateNewPlan(request *PlanRequest, ctx context.Context) (*models.Plan, error) {
	plan := &models.Plan{
		Customer:      request.Customer,
		NumberVehicle: request.NumberVehicle,
		DriverCode:        request.DriverCode,
		Status:        request.Status,
		Note:          request.Note,
		OperatorCode:  request.OperatorCode,
	
	}
	result := gorm.WithResult()
	err := gorm.G[models.Plan](repo.db, result).Create(ctx, plan)
	if err != nil {
		return nil, err
	}
	return plan, nil
}
