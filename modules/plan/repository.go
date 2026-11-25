package plan

import (
	"context"
	"maintenance-system-go/helper"
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

func (repo *Repository) GetListPlan(whereStruct models.Query, c context.Context, driverName string, operatorName string) ([]models.Plan, int64, error) {
	var plans []models.Plan
	query, _, err := helper.GrandGetAllInfo(repo.db, whereStruct, c, models.Plan{})
	if err != nil {
		return nil, 0, err
	}
	
	if driverName != ""  {
		query = query.Where(`
			EXISTS (
				SELECT 1
				FROM users u
				WHERE u.code = plan.driver_code
				AND u.name LIKE ?
			)`, "%"+driverName+"%")
	}

	if operatorName != ""  {
		query = query.Where(`
			EXISTS (
				SELECT 1
				FROM users u
				WHERE u.code = plan.operator_code
				AND u.name LIKE ?
				)`, "%"+operatorName+"%")
	}
	// Execute the query (Preload must call first before Find)
	result := query.Preload("DriverAttach").Preload("OperatorAttach").Find(&plans)
	if err := result.Error; err != nil {
		return nil, 0, err
	}
	return plans, int64(result.RowsAffected), nil
}
