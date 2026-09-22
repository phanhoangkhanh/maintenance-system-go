package helper

import "net/http"

var ErrorList = map[int]string{
	http.StatusNotFound: "USER_NOT_FOUND",
	http.StatusUnauthorized: "INVALID_CREDENTIALS",
	http.StatusInternalServerError: "INTERNAL_SERVER_ERROR",
}