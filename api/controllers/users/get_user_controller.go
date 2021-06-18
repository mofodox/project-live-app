package users

import (
	"io"
	"net/http"
)

func FindUser(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "GET /api/v1/users/id")
}
