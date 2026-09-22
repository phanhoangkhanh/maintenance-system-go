package user

import (
	"maintenance-system-go/helper"
	models "maintenance-system-go/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

// GetUser direct retrieves users with full conditions from whereStruct
// func (a *User) GetUser(whereStruct models.Query, c *gin.Context) ([]models.User, error) {
// 	var users []models.User

// 	if  whereStruct.Page < 0 || whereStruct.PerPage < 0 {
// 		return nil, fmt.Errorf("limit and pagination must not be negative")
// 	}

// 	query := a.Database.WithContext(c.Request.Context()).Model(&models.User{})

// 	//COnfirm parse the User model before building the query
// 	if err := query.Statement.Parse(&models.User{}); err != nil {
// 		return nil, err
// 	}

// 	for _, condition := range whereStruct.Where {
// 		col := clause.Column{Name: condition.Key}

// 		operator := strings.ToUpper(strings.TrimSpace(condition.Compare))
// 		if operator == "" {
// 			operator = "="
// 		}
// 		switch operator {
// 		case "=", "!=", "<>", ">", ">=", "<", "<=", "LIKE", "NOT LIKE", "ILIKE", "NOT ILIKE":
// 			query = query.Where(clause.Expr{
// 				SQL:  "? " + operator + " ?",
// 				Vars: []interface{}{col, condition.Value},
// 			})
// 		default:
// 			return nil, fmt.Errorf("unsupported comparison: %q", condition.Compare)
// 		}
// 	}
// 	//Order By
// 	if strings.TrimSpace(whereStruct.OrderBy) != "" {
// 		for _, order := range strings.Split(whereStruct.OrderBy, ",") {
// 			parts := strings.Fields(order)
// 			if len(parts) < 1 || len(parts) > 2 {
// 				return nil, fmt.Errorf("invalid order by: %q", order)
// 			}
// 			col := clause.Column{Name: parts[0]}

// 			direction := "ASC"
// 			if len(parts) == 2 {
// 				direction = strings.ToUpper(parts[1])
// 			}
// 			if direction != "ASC" && direction != "DESC" {
// 				return nil, fmt.Errorf("invalid order direction: %q", parts[1])
// 			}
// 			query = query.Order(clause.OrderByColumn{Column: col, Desc: direction == "DESC"})
// 		}
// 	}
// 	//Pagination:
// 	page := whereStruct.Page
// 	if page < 1 {
// 		page = 1
// 	}

// 	pageSize := whereStruct.PerPage
// 	if pageSize < 1  {
// 		pageSize = 20
// 	}
// 	offset := (page - 1) * pageSize

// 	query = query.Offset(offset).Limit(pageSize)

// 	//Execute query SQL
// 	if err := query.Find(&users).Error; err != nil {
// 		return nil, err
// 	}
// 	return users, nil
// }

func (repo *Repository) GetUser(whereStruct models.Query, c *gin.Context) ([]models.User, error) {
	var users []models.User
	query, err := helper.GrandGetAllInfo(repo.db, whereStruct, c, models.User{})
	if err != nil {
		return nil, err
	}
	// Execute the query
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
