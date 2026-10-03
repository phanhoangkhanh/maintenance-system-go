package models

type WhereClause struct {
	Key     string
	Compare string
	Value   any
}

type Query struct {
	Where   []WhereClause
	OrderBy string
	Page    int
	PerPage int
}

type GetQueryParams struct {
	Page int `json:"-" gorm:"-" form:"page"` //no response or gorm but got from query URL request
	PerPage int `json:"-" gorm:"-" form:"per_page"` //no response or gorm but got from query URL request
	CreatedAtFrom string `json:"-" gorm:"-" form:"created_at_from" compareTime:"from,created_at"` //no response or gorm but got from query URL request
	CreatedAtTo   string `json:"-" gorm:"-" form:"created_at_to" binding:"required_with=CreatedAtFrom" compareTime:"to,created_at"` //no response or gorm but got from query URL request
}