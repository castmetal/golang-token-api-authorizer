package dtos

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

type (
	AllowAccessTokenDTO struct {
		ApiId          string `json:"api_id" validate:"required,min=2"`
		ClientId       string `json:"client_id" validate:"required,min=2"`
		Token          string `json:"token" validate:"required,min=2"`
		ResourcePath   string `json:"resource_path" validate:"required,min=2"`
		ResourceMethod string `json:"resource_method" validate:"required,min=2"`
	}
)

func (dto *AllowAccessTokenDTO) Validate() (bool, error) {
	var validate *validator.Validate

	validate = validator.New()
	err := validate.Struct(dto)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (dto *AllowAccessTokenDTO) ToBytes() ([]byte, error) {
	b, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	return b, nil
}
