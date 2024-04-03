package token

import (
	"github.com/golang-jwt/jwt/v5"
	"math"
	"time"
)

type Claims struct {
	Id   string           `json:"id"`
	Name string           `json:"name"`
	Nbf  float64          `json:"nbf"`
	Iat  float64          `json:"iat"`
	Aud  jwt.ClaimStrings `json:"aud"`
	Iss  string           `json:"iss"`
	Sub  string           `json:"sub"`
	Exp  float64          `json:"exp"`
}

func (t Claims) GetExpirationTime() (*jwt.NumericDate, error) {
	return t.parseNumericDate(t.Exp)
}

func (t Claims) GetIssuedAt() (*jwt.NumericDate, error) {
	return t.parseNumericDate(t.Iat)
}

func (t Claims) GetNotBefore() (*jwt.NumericDate, error) {
	return t.parseNumericDate(t.Nbf)
}

func (t Claims) GetIssuer() (string, error) {
	return t.Iss, nil
}

func (t Claims) GetSubject() (string, error) {
	return t.Sub, nil
}

func (t Claims) GetAudience() (jwt.ClaimStrings, error) {
	return t.Aud, nil
}

// parseNumericDate tries to parse param float64 'dt' as a number
// date. This will succeed, if the underlying type is either a [float64] or a
// [json.Number]. Otherwise, nil will be returned.
func (t Claims) parseNumericDate(dt float64) (*jwt.NumericDate, error) {
	round, frac := math.Modf(dt)
	return jwt.NewNumericDate(time.Unix(int64(round), int64(frac*1e9))), nil
}
