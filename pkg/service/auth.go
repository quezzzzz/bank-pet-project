package service

import (
	"bank"
	"bank/pkg/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	salt       = "fieogj23jt54o3gbklgisob[a;s"
	tokenTTL   = 30 * time.Hour
	signingKey = "asdAfkvlpSV=34;"
)

type tokenClaims struct {
	jwt.StandardClaims
	CustomerId int `json:"customer_id"`
}
type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo}
}

func (s *AuthService) CreateCustomer(customer bank.Customer) (int, error) {
	customer.Password = generateHashedPassword(customer.Password)
	customer.Balance = 0
	return s.repo.CreateCustomer(customer)
}

func (s *AuthService) GenerateToken(phone, password string) (string, error) {
	customer, err := s.repo.GetCustomer(phone, generateHashedPassword(password))
	if err != nil {
		logrus.Error("Wrong phone or password")
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		customer.Id,
	})

	return token.SignedString([]byte(signingKey))

}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}
	return claims.CustomerId, nil
}

func generateHashedPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
