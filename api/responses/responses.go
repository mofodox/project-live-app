package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func JSON(res http.ResponseWriter, statusCode int, data interface{}) {
	res.WriteHeader(statusCode)
	err := json.NewEncoder(res).Encode(data)
	if err != nil {
		fmt.Fprintf(res, "%s", err.Error())
	}
}

func ERROR (res http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(res, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})

		return
	}

	JSON(res, http.StatusBadRequest, nil)
}