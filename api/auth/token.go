package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(res http.ResponseWriter, userId uint32) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // Token expires after 1 hour

	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// ! OLD CODE - DON'T DELETE FIRST
	// claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
	// 	Issuer: strconv.Itoa(int(userId)),
	// 	ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // Expires in 1 hour
	// })

	token, err := tk.SignedString([]byte(os.Getenv("HBB_SECRET_KEY")))
	if err != nil {
		log.Fatalf("token error %s\n", err)
	}

	cookie := &http.Cookie{
		Name:     "jwt-token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 1),
		HttpOnly: true,
	}

	http.SetCookie(res, cookie)

	return token, nil
}

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
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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

// todo: possible to make it return user?
func GetTokenID(tokenString string) (uint32, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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
