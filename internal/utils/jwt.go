package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/IRONICBo/QiYin_BE/internal/config"
)

type JwtClaims struct {
	UserUUID      string `json:"user_uuid"`
	CommunityUUID string `json:"community_uuid"`
	jwt.RegisteredClaims
}

func GenerateJwtToken(user_uuid string, community_uuid string) (string, uint, error) {
	secret := []byte(config.Config.JWT.Secret)
	issuer := config.Config.JWT.Issuer
	expireDays := config.Config.JWT.ExpireDays

	claims := &JwtClaims{
		user_uuid,
		community_uuid,
		jwt.RegisteredClaims{
			Issuer:    issuer,
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)),
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, expireDays)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_string, err := token.SignedString(secret)
	expire_time_seconds := uint(time.Now().AddDate(0, 0, expireDays).Unix())

	return token_string, expire_time_seconds, err
}

func ParseJwtToken(tokenString string) (*JwtClaims, error) {
	secret := []byte(config.Config.JWT.Secret)
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
