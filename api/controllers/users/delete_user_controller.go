package users

import (
	"io"
	"net/http"
)

func DeleteUser(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "DELETE /api/v1/users/id")
}
