package dtos

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

type (
	GenerateTokenDTO struct {
		ApiId    string `json:"api_id" validate:"required,min=2"`
		ClientId string `json:"client_id" validate:"required,min=2"`
	}
)

func (dto *GenerateTokenDTO) Validate() (bool, error) {
	var validate *validator.Validate

	validate = validator.New()
	err := validate.Struct(dto)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (dto *GenerateTokenDTO) ToBytes() ([]byte, error) {
	b, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	return b, nil
}
