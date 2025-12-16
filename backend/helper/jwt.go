package helper

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var SECRETKEY string = getJWTSecret()

func getJWTSecret() string {
	secret := GetEnv("JWT_SECRET")
	if secret == "" {
		secret = GetEnv("SECRETKEY")
		if secret == "" {
			secret = "your-secret-key"
		}
	}
	return secret
}

func GenerateToken(idUser string, nmUser string) (string, error) {
	claims := jwt.MapClaims{
		"id_user": idUser,
		"nm_user": nmUser,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	res, err := jwt.SignedString([]byte(SECRETKEY))

	return res, err
}

// GenerateTokenWithDomain generates JWT token with user, domain, and role context
func GenerateTokenWithDomain(userID int64, username string, email string, domainID int64, roleID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"username":  username,
		"email":     email,
		"domain_id": domainID,
		"role_id":   roleID,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // 24 hours
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	res, err := token.SignedString([]byte(SECRETKEY))

	return res, err
}

func VerifyToken(ctx *gin.Context) (jwt.MapClaims, error) {
	auth := ctx.Request.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return nil, errors.New("please login to get the token")
	}

	tokenStr := strings.TrimPrefix(auth, "Bearer ")

	// Parse token with MapClaims
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(SECRETKEY), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
