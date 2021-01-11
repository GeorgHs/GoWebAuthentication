package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func main() {
	http.HandleFunc("/", foo)
	http.HandleFunc("/submit", bar)
	http.ListenAndServe(":8080", nil)
}

const myKey = "I love thursday when it rains 8723 inches"

type myClaims struct {
	jwt.StandardClaims
	Email string
}

func getJWT(msg string) (string, error) {

	claims := myClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
		Email: msg,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(myKey))
	if err != nil {
		return "", fmt.Errorf("couldn't SignedString in NewWithClaims %w", err)
	}
	return ss, nil
}

func bar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ss, err := getJWT(email)
	if err != nil {
		http.Error(w, "couldn't getJWT", http.StatusInternalServerError)
		return
	}

	// "hash / message digest / digest / hash value" | "what we stored"
	c := http.Cookie{
		Name:  "session",
		Value: ss,
	}

	http.SetCookie(w, &c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func foo(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		c = &http.Cookie{}
	}

	ss := c.Value
	afterVerificationToken, err := jwt.ParseWithClaims(ss, &myClaims{}, func(beforeVerificationToken *jwt.Token) (interface{}, error) {
		if beforeVerificationToken.Method.Alg() != jwt.SigningMethodES256.Alg() {
			return nil, fmt.Errorf("SOMEONE TRIED TO HACK changed signing method %w")
		}

		return []byte(myKey), nil
	})

	// StandardClaims has the ...
	// Valid() error
	// ... method which means it implements the Claims interface ...
	//
	//

	isEqual := err == nil && afterVerificationToken.Valid

	message := "Not logged in"
	if isEqual {
		message = "Logged in"
		claims := afterVerificationToken.Claims.(*myClaims)
		fmt.Println(claims.Email)
		fmt.Println(claims.ExpiresAt)
	}

	html := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>HMAC Example</title>
	</head>
	<body>
		<p>Cookie value: ` + c.Value + `</p>
		<p>` + message + `</p>
		<form action="/submit" method="post">
			<input type="email" name="email" />
			<input type="submit" />
		</form>
	</body>
	</html>`
	io.WriteString(w, html)
}
