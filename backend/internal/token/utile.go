package token

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const TOKEN_LITERAL = "token"

// GetTokenByRequest with http.Header by header name TOKEN_LITERAL get token string
func GetTokenByRequest(c *gin.Context) string {
	header := c.Request.Header
	return header.Get(TOKEN_LITERAL)
}

// GetClaimsByRequest with http.Header by header name TOKEN_LITERAL get Claims
func GetClaimsByRequest(c *gin.Context, signedKey string) (jwt.Claims, error) {
	tokenString := GetTokenByRequest(c)

	claims, err := ParseJwtToken(tokenString, signedKey)
	if err != nil {
		return nil, err
	}
	return claims, err
}
