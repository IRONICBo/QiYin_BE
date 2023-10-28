package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// GetUserID get user uuid from context.
func GetUserUUID(c *gin.Context) (string, error) {
	if claims, ok := c.Get("claims"); ok {
		if c := claims.(*JwtClaims); c != nil {
			return c.UserUUID, nil
		}
	}

	return "", errors.Errorf("get user uuid failed")
}

//// GetCommunityUUID get community uuid from context.
//func GetCommunityUUID(c *gin.Context) (string, error) {
//	if claims, ok := c.Get("claims"); ok {
//		if c := claims.(*JwtClaims); c != nil {
//			return c.CommunityUUID, nil
//		}
//	}
//
//	return "", errors.Errorf("get community uuid failed")
//}
