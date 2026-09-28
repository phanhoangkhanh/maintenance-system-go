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
