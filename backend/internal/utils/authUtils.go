package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type ParsedToken struct {
	UserId  int
	Email   string
	IsAdmin bool
}

func CreateToken(user_id, email, is_admin string) (string, error) {
	now := time.Now()
	token_duration, err := strconv.Atoi(os.Getenv("TOKEN_EXPIRY_DURATION"))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id":    user_id,
		"email":      email,
		"is_admin":   is_admin,
		"created_at": now.Unix(),
		"expires_at": now.Add(time.Duration(token_duration) * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Parsing and validating the token
func ParseAndValidateToken(tokenString string) (*ParsedToken, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	// Parse and validate claim
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("Invalid token signature")
		} else {
			return nil, errors.New("Invalid token")
		}
	}

	// Extracting and validating the claims.
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, errors.New("Invalid token claims")
	}

	// Checking token expiration
	expiry_time, _ := claims["expires_at"].(float64)

	if int64(expiry_time) < time.Now().Unix() {
		return nil, errors.New("Token expired")
	}
	parsed_token := ParsedToken{
		UserId:  claims["user_id"].(int),
		Email:   claims["email"].(string),
		IsAdmin: claims["is_admin"].(bool),
	}
	return &parsed_token, nil
}

// Hash password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// Compare hash with password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Validate the password
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("Password length can not be less than 8!")
	}
	return nil
}
