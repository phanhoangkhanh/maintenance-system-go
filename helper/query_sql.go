package helper

import (
	"context"
	"fmt"
	"maintenance-system-go/models"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// GenerateWhereStruct builds text and scalar filters, --- Skipping timestamps and hidden fields.
// Strings, including dereferenced string pointers, use substring matching.
//Return WhereClause :
// 	Where: []models.WhereClause{
// 		{
// 			Key: "name",
// 			Compare: "ILIKE",
// 			Value: "%"+loginRequest.Name+"%",
// 		},
// 	},
func GenerateWhereStruct[T any](model *T) models.Query {
	query := models.Query{}
	if model == nil {
		return query
	}
	//As Generic func -> using reflect to handle any struct type , numfiled, and their tags dynamically
	value := reflect.ValueOf(model).Elem()
	// ex: model := &User{
		//     Name: "An",
		//     Age:  25,
		// }
		// ==> Value := reflect.ValueOf(model).Elem() ==> will be : User{Name: "An", Age: 25}
		// ==> value.Field(0).Interface() // "An"
		// 	value.Field(1).Interface() // 25

	if value.Kind() != reflect.Struct {
		return query
	}
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		condition := models.WhereClause{}
		key := typ.Field(i) //key: {Name  string json:"name" gorm:"column:name;uniqueIndex" form:"name" 0 [0] false}
		field := value.Field(i) // field: khanh
		//Handle compare Time FROM_TO:
		if key.Tag.Get("compareTime") != "" {
			if field.Kind() != reflect.String || field.String() == "" {
				continue
			}
			//compareTime:"to,created_at"
			parts := strings.Split(key.Tag.Get("compareTime"), ",")
			if len(parts) == 2 {
				if parts[0] == "from" {
					condition.Compare = ">="
					condition.Value = field.String() + " 00:00:00"
				} else if parts[0] == "to" {
					condition.Value = field.String() + " 23:59:59"
					condition.Compare = "<="
				}
				condition.Key = parts[1]
				query.Where = append(query.Where, condition)

				// log.Fatalf("CONDITIO %v", condition)
				continue
			}
		}
		if !field.CanInterface() || key.Tag.Get("json") == "-" || key.Tag.Get("gorm") == "-" {
			continue
		}
		column := schema.NamingStrategy{}.ColumnName("", key.Name)
		for _, setting := range strings.Split(key.Tag.Get("gorm"), ";") {
			if name, ok := strings.CutPrefix(setting, "column:"); ok {
				column = name
			}
		}
		if column == "created_at" || column == "updated_at" || column == "deleted_at" {
			continue
		}
		isPointer := field.Kind() == reflect.Ptr //reflect Pointer
		//filed is pointer and not nil
		for field.Kind() == reflect.Ptr && !field.IsNil() {
			field = field.Elem() // parse pointer *string → string
		}
		//not pointer and zero value
		if !isPointer && field.IsZero() {
			continue
		}
		//pointer and zero value
		if isPointer && field.IsZero() {
			continue
		}
		
		switch field.Kind() {
		case reflect.String:
			condition.Key = column
			condition.Compare = "ILIKE"
			condition.Value = "%" + field.String() + "%"
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			condition.Key = column
			condition.Compare = "="
			condition.Value = field.Interface()
		default:
			continue
		}
		query.Where = append(query.Where, condition)
	}
	return query
}
// Using Clause to more dynamic control all conditions
func GrandGetAllInfo[T any](db *gorm.DB, whereStruct models.Query, c context.Context, modelName T) (*gorm.DB, error) {

	if whereStruct.Page < 0 || whereStruct.PerPage < 0 {
		return nil, fmt.Errorf("limit and pagination must not be negative")
	}
	query := db.WithContext(c).Model(&modelName)
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
		//install wheres 
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
