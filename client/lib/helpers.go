package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/mofodox/project-live-app/api/auth"
	"github.com/mofodox/project-live-app/api/models"
)

var apiBaseURL string

func init() {
	apiBaseURL = os.Getenv("APIServerHostname") + ":" + os.Getenv("APIServerPort") + os.Getenv("APIServerBasePath")
}

func IsLoggedIn(req *http.Request) (*models.User, error) {
	myCookie, err := req.Cookie("jwt-token")
	if err != nil {
		return nil, err
	}

	token := myCookie.Value

	userID, err := auth.GetTokenID(token)
	if err != nil {
		return nil, err
	}

	// Todo: add cookie check and send JWT with request
	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodGet, apiBaseURL+"/users/"+strconv.FormatUint(uint64(userID), 10), nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)

	// handle error
	if err != nil {
		fmt.Println("error sending get user info request", err)
		return nil, err
	}
	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)

	// success
	if response.StatusCode == 200 {
		var user *models.User
		marshalErr := json.Unmarshal(data, &user)

		if marshalErr != nil {
			return nil, marshalErr
		}

		return user, nil
	}

	return nil, errors.New("unable to fetch user")
}
