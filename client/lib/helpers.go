package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
	"unicode"

	"github.com/joho/godotenv"
	"github.com/mofodox/project-live-app/api/auth"
	"github.com/mofodox/project-live-app/api/models"
)

var Tpl *template.Template
var ApiBaseURL string

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error getting env values %v\n", err)
	} else {
		log.Println("Successfully loaded the env values")
	}

	ApiBaseURL = os.Getenv("APIServerHostname") + ":" + os.Getenv("APIServerPort") + os.Getenv("APIServerBasePath")

	funcMap := template.FuncMap{
		"add": func(a int, b int) int {
			return a + b
		},

		"ucFirst": func(str string) string {
			if len(str) == 0 {
				return ""
			}
			tmp := []rune(str)
			tmp[0] = unicode.ToUpper(tmp[0])
			return string(tmp)
		},

		"formatDistance": func(distance float64) string {
			if distance > 1 {
				return fmt.Sprint(math.Floor(distance*100)/100) + "KM"
			} else {
				return fmt.Sprint(math.Floor(distance*100)) + "M"
			}
		},

		"formatCommentDate": func(t time.Time) string {
			return t.Format("2 Jan 2006 3:04 PM")
		},
	}

	Tpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*"))
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

	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodGet, ApiBaseURL+"/users/"+strconv.FormatUint(uint64(userID), 10), nil)
	//request.Header.Set("Content-Type", "application/json")
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

func GetJWT(req *http.Request) string {

	myCookie, err := req.Cookie("jwt-token")
	if err == nil {
		return myCookie.Value
	}

	return ""
}
