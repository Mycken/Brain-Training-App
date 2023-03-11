package entity

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Email:    "user@example.org",
		Password: "123456",
		Username: "username",
	}
}
