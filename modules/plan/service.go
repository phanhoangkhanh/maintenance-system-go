package plan

import (
	"log"
	"maintenance-system-go/helper"
	"maintenance-system-go/models"

	"github.com/gin-gonic/gin"
)

type Service struct {
	Repo  *Repository
}

func (service *Service) CreatePlan(request *PlanRequest, ctx *gin.Context) (*models.Plan, error, int) {
	userActive, err := helper.GetUserActive(ctx)
	if err != nil {
		return nil, err, 500
	}
	request.OperatorCode = userActive.Code
	plan, err := service.Repo.CreateNewPlan(request, ctx)
	if err != nil {
		return nil, err, 500
	}
	return plan, nil, 200
}

func (service *Service) UpdatePlan(request *PlanRequest, ctx *gin.Context) (*models.Plan, error, int) {
	plan, err := service.Repo.UpdatePlan(request, ctx)
	if err != nil {
		return nil, err, 500
	}
	return plan, nil, 200
}

func (service *Service) GetListPlan(ctx *gin.Context) ([]models.Plan, int64, error, int) {

	var planQuery models.Plan 
	if err := ctx.ShouldBindQuery(&planQuery); err != nil {
		return nil, 0, err, 400
	}
	
	query := helper.GenerateWhereStruct(&planQuery)
	query.OrderBy = "updated_at DESC"
	query.Page = planQuery.Page
	query.PerPage = planQuery.PerPage

	// Query For Driver Name and Operator Name from other table
	driverName := ""
	operatorName := ""
	if planQuery.DriverName != "" {
		driverName = planQuery.DriverName
	}
	if planQuery.OperatorName != "" {
		operatorName = planQuery.OperatorName
	}
	log.Println(driverName, operatorName, planQuery)

	plans, total, err := service.Repo.GetListPlan(query, ctx, driverName, operatorName)
	if err != nil {
		return nil, 0, err, 500
	}
	return plans, total, nil, 200
}