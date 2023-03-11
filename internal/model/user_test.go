package model_test

import (
	"BrainTApp/internal/model/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *entity.User
		isValid bool
	}{
		{
			name: "valid",
			u: func() *entity.User {
				return entity.TestUser(t)
			},
			isValid: true,
		}, //valid
		{
			name: "with encrypt password",
			u: func() *entity.User {
				u := entity.TestUser(t)
				u.Password = ""
				u.EncryptedPassword = "encrypt password"
				return u
			},
			isValid: true,
		}, //with_encrypt_password
		{
			name: "empty email",
			u: func() *entity.User {
				u := entity.TestUser(t)
				u.Email = ""
				return u
			},
			isValid: false,
		}, //empty_email
		{
			name: "invalid email",
			u: func() *entity.User {
				u := entity.TestUser(t)
				u.Email = "invalid"
				return u
			},
			isValid: false,
		}, //invalid_email
		{
			name: "empty password",
			u: func() *entity.User {
				u := entity.TestUser(t)
				u.Password = ""
				return u
			},
			isValid: false,
		}, //empty_password
		{
			name: "short password",
			u: func() *entity.User {
				u := entity.TestUser(t)
				u.Password = "123"
				return u
			},
			isValid: false,
		}, //short_password
		{
			name: "empty username",
			u: func() *entity.User {
				u := entity.TestUser(t)
				u.Username = ""
				return u
			},
			isValid: false,
		}, //empty_username
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}

func TestUser_BeforeCreate(t *testing.T) {
	u := entity.TestUser(t)
	assert.NoError(t, u.BeforeCreate())
	assert.NotEmpty(t, u.EncryptedPassword)
}
