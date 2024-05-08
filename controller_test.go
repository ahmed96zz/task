package main

import (
	"testing"
)

type createTest struct {
	user      ProfileReq
	expectErr bool
}

func TestUserTransaction(t *testing.T) {
	initDB()
	testNumber1, _ := generateOTP("1234")
	testNumber2, _ := generateOTP("2345")
	var tests = []createTest{
		{user: ProfileReq{FirstName: "First test user", LastName: "055555" + testNumber1},
			expectErr: false,
		},
		// Too long phone number that exceed the length in DB
		{user: ProfileReq{FirstName: "Failed test User", LastName: "01111111111111111" + testNumber2},
			expectErr: true,
		},
	}
	for _, test := range tests {
		err := UserTransaction(test.user)
		check := err != nil
		if check != test.expectErr {
			t.Fatal(err)
		}
	}
}
