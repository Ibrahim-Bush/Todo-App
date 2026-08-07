package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"errors"
)

var (
	ErrInvalidPassword = errors.New("Invalid password")
	ErrInvalidToken = errors.New("Invalid token")
	ErrInvalidClaims = errors.New("Invalid claims")
)

//define a secret key for signature.
var jwt_signature = "secret_key_4893586"

func Hash_password(password string) (string, error){
	//use bcrypt to hash password.
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//check the result of hashing.
	if err == nil{
	return string(hashed_password), nil
	}
	return "", ErrInvalidPassword
}

func Compare_password(password, hashed string) bool {
	//compare password with hashed value.
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	//check the result.
	if err == nil{
		return true
	}
	return false
}

func Generate_token(username, role string, id int) (string, error){
	//add required data in a map claims.
	claims := jwt.MapClaims{
		"userID": id,
		"username": username,
		"role": role,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	//create token with the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//add signature to secure token.
	signed_token, err := token.SignedString([]byte(jwt_signature))
	return signed_token, err
}

func get_secret_key(token *jwt.Token) (interface{}, error){
	return []byte(jwt_signature), nil
}

func Verify_token(input_token string) (jwt.MapClaims, error){
	//parse the token parts.
	token, err := jwt.Parse(input_token, get_secret_key)
	//check the result.
	if err != nil{
		return nil, err
	}
	//if the token is valid.
	if token.Valid == true{
		//extract claims.
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok{
			return claims, nil
		}
		return nil, ErrInvalidClaims
	}
	return nil, ErrInvalidToken
}