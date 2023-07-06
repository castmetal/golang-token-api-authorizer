package dtos

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

type (
	CreateClientResponseDTO struct {
		ID              string       `json:"id"`
		ClientName      string       `json:"client_name"`
		ScopeId         string       `json:"scope_id"`
		ScopeName       string       `json:"scope_name"`
		ApiId           string       `json:"api_id"`
		Permissions     []Permission `json:"permissions"`
		KeyTimeDuration int32        `json:"key_time_duration"`
		KeyPeriod       string       `json:"key_period"`
	}
)

func (dto *CreateClientResponseDTO) Validate() (bool, error) {
	var validate *validator.Validate

	validate = validator.New()
	err := validate.Struct(dto)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (dto *CreateClientResponseDTO) ToBytes() ([]byte, error) {
	b, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	return b, nil
}
