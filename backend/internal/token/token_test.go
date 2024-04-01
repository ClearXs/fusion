package token

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestCreateJwtToken(t *testing.T) {
	a := assert.New(t)
	exp := time.Now().Add(math.MaxInt)
	token, err := CreateJwtToken(&Claims{Id: 0, Name: "123", Exp: float64(exp.Unix())}, "turbo")
	a.Nil(err)

	a.NotNil(token)
}

func TestParseJwtToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MCwibmFtZSI6IjEyMyIsIm5iZiI6MCwiaWF0IjowLCJhdWQiOm51bGwsImlzcyI6IiIsInN1YiI6IiIsImV4cCI6MTA5MzUzMTU4MjJ9.Ai8vK1eo-7QAGNIHCOl4xrP0cTCrq1pF60rpFrrbd7w"
	a := assert.New(t)
	claims, err := ParseJwtToken(tokenString, "turbo")
	a.Nil(err)
	a.NotNil(claims)
}
