package plan

import (
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