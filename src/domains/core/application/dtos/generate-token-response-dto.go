package dtos

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

type (
	GenerateTokenResponseDTO struct {
		Token string `json:"token"`
	}
)

func (dto *GenerateTokenResponseDTO) Validate() (bool, error) {
	var validate *validator.Validate

	validate = validator.New()
	err := validate.Struct(dto)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (dto *GenerateTokenResponseDTO) ToBytes() ([]byte, error) {
	b, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	return b, nil
}
