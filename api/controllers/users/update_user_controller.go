package users

import (
	"io"
	"net/http"
)

func UpdateUser(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "PUT /api/v1/users/id")
}
