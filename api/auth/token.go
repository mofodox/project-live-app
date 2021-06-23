package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func TokenValid(req *http.Request) error {
	tokenString := ExtractToken(req)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		
		return []byte(os.Getenv("HBB_SECRET_KEY")), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		b, err := json.MarshalIndent(claims, "", " ")
		if err != nil {
			log.Println(err)
			return err
		}

		fmt.Println(string(b))
	}

	return nil
}

func ExtractToken(req *http.Request) string {
	keys := req.URL.Query()
	token := keys.Get("jwt-token")
	if token != "" {
		return token
	}

	bearerToken := req.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}

func ExtractTokenID(req *http.Request) (uint32, error) {
	tokenString := ExtractToken(req)
	token, err := jwt.Parse(tokenString, func (token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("HBB_SECRET_KEY")), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}

		fmt.Println("jwt-token")
		return uint32(uid), nil
	}

	return 0, nil
}