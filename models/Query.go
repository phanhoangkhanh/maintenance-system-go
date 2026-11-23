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
	UpdatedAtFrom string `json:"-" gorm:"-" form:"updated_at_from" compareTime:"from,updated_at"` //no response or gorm but got from query URL request
	UpdatedAtTo   string `json:"-" gorm:"-" form:"updated_at_to" binding:"required_with=UpdatedAtFrom" compareTime:"to,updated_at"` //no response or gorm but got from query URL request
}