package users

import (
	"io"
	"net/http"
)

func CreateUserHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "POST /api/v1/users")
}

func GetAllUsersHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "GET /api/v1/users")
}

func GetUserHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "GET /api/v1/users/id")
}

func UpdateUserHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "PUT /api/v1/users/id")
}

func DeleteUserHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "DELETE /api/v1/users/id")
}

