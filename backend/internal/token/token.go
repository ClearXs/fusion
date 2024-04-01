package token

import "github.com/golang-jwt/jwt/v5"

// CreateJwtToken wrapper jwt.NewWithClaims method, use jwt.SigningMethodHS256 signing method encrypt claims.
func CreateJwtToken(claims jwt.Claims, signedKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	singed, err := token.SignedString([]byte(signedKey))
	if err != nil {
		return "", err
	}
	return singed, nil
}

// ParseJwtToken by tokenString, parse to jwt.Claims
func ParseJwtToken(tokenString string, signedKey string) (jwt.Claims, error) {
	t, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signedKey), nil
	})
	if err != nil {
		return nil, err
	}
	return t.Claims, nil
}
