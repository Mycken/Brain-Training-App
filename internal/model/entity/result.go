package entity

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

type Result struct {
	UserID      int           `json:"user_id,omitempty"`
	TestID      int           `json:"test_id"`
	DateTest    time.Time     `json:"date_test"`
	ResultInter time.Duration `json:"result_inter"`
	ResultOne   int           `json:"result_one"`
	ResultTwo   int           `json:"result_two"`
	Results     string        `json:"results"`
	TestSet     string        `json:"test_set"`
}

func (res *Result) Validate() error {
	return validation.ValidateStruct(
		res)
}
