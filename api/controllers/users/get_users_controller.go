package users

import (
	"io"
	"net/http"
)

func FindAllUsers(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "GET /api/v1/users")
}
