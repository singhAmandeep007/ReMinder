package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(claims map[string]interface{}, secretKey []byte) (string, error) {
	// Add the expiration claim if it's not already present.
	if _, ok := claims["expiresAt"]; !ok {
		claims["expiresAt"] = time.Now().Add(time.Hour * 24).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	return token.SignedString(secretKey)
}

func VerifyToken(token string, secretKey []byte) (*jwt.Token, error) {
	// jwt.Parse takes a token and a function that returns the secret key
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// check if the signing method matches
		// this is a type assertion
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		// if the signing method is not valid, return an error
		if !ok {
			return nil, errors.New("invalid signing method")
		}

		// return the secret key
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, errors.New("could not parse token")
	}

	// the token will show Valid as false if the expiry time has elapsed.
	isTokenValid := parsedToken.Valid

	if !isTokenValid {
		return nil, errors.New("invalid token")
	}

	return parsedToken, nil
}

func GetClaimsFromToken(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("could not parse claims")
	}

	return claims, nil
}
