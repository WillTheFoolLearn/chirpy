package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()
	tokenSecret1 := "cheese"
	tokenSecret2 := "BillieHoliday"
	expires1 := 10 * time.Second
	expires2 := 1 * time.Minute
	expires3 := -5 * time.Second

	testList := []JWTTests{
		{
			TestName:   "Test 1",
			User:       user1,
			Secret:     tokenSecret1,
			TestSecret: tokenSecret1,
			Expires:    expires1,
			Result:     true,
		},
		{
			TestName:   "Test 2",
			User:       user2,
			Secret:     tokenSecret2,
			TestSecret: tokenSecret2,
			Expires:    expires2,
			Result:     true,
		},
		{
			TestName:   "Test 3",
			User:       user1,
			Secret:     tokenSecret1,
			TestSecret: tokenSecret2,
			Expires:    expires2,
			Result:     false,
		},
		{
			TestName:   "Test 4",
			User:       user2,
			Secret:     tokenSecret2,
			TestSecret: tokenSecret1,
			Expires:    expires1,
			Result:     false,
		},
		{
			TestName:   "Negative Time",
			User:       user1,
			Secret:     tokenSecret1,
			TestSecret: tokenSecret1,
			Expires:    expires3,
			Result:     false,
		},
	}

	for _, test := range testList {
		fmt.Printf("Testing %s\n", test.TestName)
		token, err := MakeJWT(test.User, test.Secret, test.Expires)
		if err != nil {
			t.Errorf("Unable to make JWT: %v", err)
		}

		validatedUser, err := ValidateJWT(token, test.TestSecret)
		if !(err != nil) != test.Result {
			t.Errorf("Validated JWT error = %v, test.Result = %v", err != nil, test.Result)
		}
		if (validatedUser == test.User) != test.Result {
			t.Errorf("validatedUser = %v, test.User = %v, test.Result = %v", validatedUser, test.User, test.Result)
		}
	}
}

type JWTTests struct {
	TestName   string
	User       uuid.UUID
	Secret     string
	TestSecret string
	Expires    time.Duration
	Result     bool
}

func TestBearerToken(t *testing.T) {
	header := make(http.Header)
	header.Add("Authorization", "Bearer cheese")

	tokenString, err := GetBearerToken(header)
	if tokenString == "" || err != nil {
		t.Errorf("TokenString is %s and err is %v", tokenString, err)
	}
}
