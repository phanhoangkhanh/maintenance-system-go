package helper

import (
	"fmt"
	"maintenance-system-go/models"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Helper struct {
}

func InitHelper() *Helper {
	return &Helper{}
}

func  GrandGetAllInfo[T any](db *gorm.DB, whereStruct models.Query, c *gin.Context, modelName T) ( *gorm.DB,error) {

	if  whereStruct.Page < 0 || whereStruct.PerPage < 0 {
		return nil, fmt.Errorf("limit and pagination must not be negative")
	}
	query := db.WithContext(c.Request.Context()).Model(&modelName)
	//1 - COnfirm parse the T model before building the query
	if err := query.Statement.Parse(&modelName); err != nil {
		return nil, err
	}
	//struture of whereClause need build :
	// whereClause := clause.Where{
	// Exprs: []clause.Expression{
	// 		clause.Eq{
	// 			Column: clause.Column{Name: "status"},
	// 			Value:  "active",
	// 		},
	// 		clause.Gte{
	// 			Column: clause.Column{Name: "age"},
	// 			Value:  18,
	// 		},
	// 	},
	// }

	//Request like this
	// query := models.Query{
	// 	Where: []models.WhereClause{
	// 		{
	// 			Key: "name",
	// 			Compare: "ILIKE",
	// 			Value: "%"+loginRequest.Name+"%",
	// 		},
	// 	},
	// 	OrderBy: "created_at DESC , name ASC",
	// 	Page: 1,
	// 	PerPage: 2,
	// }

	//2 - Build the WHERE clauses based on the provided conditions
	for _, condition := range whereStruct.Where {
		col := clause.Column{Name: condition.Key}
		
		operator := strings.ToUpper(strings.TrimSpace(condition.Compare))
		if operator == "" {
			operator = "="
		}
		switch operator {
		case "=", "!=", "<>", ">", ">=", "<", "<=", "LIKE", "NOT LIKE", "ILIKE", "NOT ILIKE":
			query = query.Where(clause.Expr{
				SQL:  "? " + operator + " ?",
				Vars: []interface{}{col, condition.Value},
			})
		default:
			return nil, fmt.Errorf("unsupported comparison: %q", condition.Compare)
		}
	}
	//3 - Order By
	if strings.TrimSpace(whereStruct.OrderBy) != "" {
		for _, order := range strings.Split(whereStruct.OrderBy, ",") {
			parts := strings.Fields(order)
			if len(parts) < 1 || len(parts) > 2 {
				return nil, fmt.Errorf("invalid order by: %q", order)
			}
			col := clause.Column{Name: parts[0]}
			
			direction := "ASC"
			if len(parts) == 2 {
				direction = strings.ToUpper(parts[1])
			}
			if direction != "ASC" && direction != "DESC" {
				return nil, fmt.Errorf("invalid order direction: %q", parts[1])
			}
			query = query.Order(clause.OrderByColumn{Column: col, Desc: direction == "DESC"})
		}
	}

	//4 - Pagination
	page := whereStruct.Page
	if page < 1 {
		page = 1
	}

	pageSize := whereStruct.PerPage
	if pageSize >= 1 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}
	return query, nil
}