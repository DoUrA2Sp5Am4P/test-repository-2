package main

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type TestRegisterType struct {
	user *User
	err  bool
}

var testsReg = []TestRegisterType{
	{&User{Name: randTest, Password: "password"}, false},
	{&User{Name: randTest, Password: "password"}, true},
}

var randTest = ""

func CompareHashAndPassword_(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true
	}
	return false
}
func TestRegisterFunc(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randTest = strconv.Itoa(r.Int())
	for _, test := range testsReg {
		res, err := register(test.user.Name, test.user.Name)
		if err != nil && test.err == true {
			return
		}
		if res.Name != test.user.Name || CompareHashAndPassword_(test.user.Password, res.Password) {
			t.Error("Error in TestRegisterFunc")
		}
	}
}
