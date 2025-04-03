package auth

import "testing"

func TestPasswords(t *testing.T) {
	password1 := "snakeeater"
	password2 := "keptyouwaiting"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	testList := []Tests{
		{
			Name:     "password1 and hash1 match",
			Password: password1,
			Hash:     hash1,
			Result:   true,
		},
		{
			Name:     "password2 and hash1 don't match",
			Password: password2,
			Hash:     hash1,
			Result:   false,
		},
		{
			Name:     "password2 and hash2 match",
			Password: password2,
			Hash:     hash2,
			Result:   true,
		},
		{
			Name:     "password1 and hash2 don't match",
			Password: password1,
			Hash:     hash2,
			Result:   false,
		},
	}

	for _, test := range testList {
		err := CheckPasswordHash(test.Password, test.Hash)
		if !(err != nil) != test.Result {
			t.Errorf("CheckPasswordHash error = %v, test.Result = %v", err, test.Result)
		}
	}
}

type Tests struct {
	Name     string
	Password string
	Hash     string
	Result   bool
}
