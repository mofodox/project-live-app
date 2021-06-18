package users

import (
	"io"
	"net/http"
)

func CreateUser(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "POST /api/v1/users")
}
