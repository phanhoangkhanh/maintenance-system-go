package helper

import (
	"reflect"
	"testing"
	"time"

	"maintenance-system-go/models"
)

func TestGenerateWhereStruct(t *testing.T) {
	email := "phkhanhgmail.com"
	now := time.Now()
	user := models.User{Name: "khanh", Email: &email, MobilePhone: "097394", Password: "hidden", CreatedAt: &now, UpdatedAt: &now, DeletedAt: &now}
	got := GenerateWhereStruct(&user)
	want := []models.WhereClause{
		{Key: "name", Compare: "ILIKE", Value: "%khanh%"},
		{Key: "email", Compare: "ILIKE", Value: "%" + email + "%"},
		{Key: "mobile_phone", Compare: "ILIKE", Value: "%097394%"},
	}
	if !reflect.DeepEqual(got.Where, want) {
		t.Fatalf("got %#v; want %#v", got.Where, want)
	}
}

func TestGenerateWhereStructEmptyAndInvalid(t *testing.T) {
	if got := GenerateWhereStruct[models.User](nil); len(got.Where) != 0 {
		t.Fatal(got)
	}
	if got := GenerateWhereStruct(&models.User{}); len(got.Where) != 0 {
		t.Fatal(got)
	}
	number := 1
	if got := GenerateWhereStruct(&number); len(got.Where) != 0 {
		t.Fatal(got)
	}
}

func TestGenerateWhereStructPointerZeroValues(t *testing.T) {
	empty, zero, disabled := "", 0, false
	filters := struct {
		Email   *string
		Count   *int `gorm:"column:total_count"`
		Enabled *bool
	}{&empty, &zero, &disabled}
	// Zero values are skipped, including values reached through pointers.
	var want []models.WhereClause
	if got := GenerateWhereStruct(&filters); !reflect.DeepEqual(got.Where, want) {
		t.Fatalf("got %#v; want %#v", got.Where, want)
	}
}
