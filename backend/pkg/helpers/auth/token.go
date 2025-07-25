package auth

import (
	"crypto/rand"
	"math/big"
	"net/http"

	"github.com/DragonSecSI/instancer/backend/pkg/database/models"
	"gorm.io/gorm"
)

type AuthToken struct {
	GetTeam       func(db *gorm.DB, r *http.Request) (*models.Team, error)
	GenerateToken func() (string, error)
}

func authTokenGetTeam(db *gorm.DB, r *http.Request) (*models.Team, error) {
	token := r.URL.Query().Get("token")
	if token == "" {
		token = r.Header.Get("Authorization")
		if token == "" {
			return nil, nil
		}
	}

	team, err := models.TeamGetByToken(db, token)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return team, nil
}

func authTokenGenerateToken() (string, error) {
	alphaNum := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 32
	token := make([]byte, length)
	for i := range token {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphaNum))))
		if err != nil {
			return "", err
		}
		token[i] = alphaNum[num.Int64()]
	}
	return string(token), nil
}
