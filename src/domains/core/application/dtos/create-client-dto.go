package dtos

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/castmetal/golang-token-api-authorizer/src/domains/resource"
	"github.com/go-playground/validator/v10"
)

type (
	Permission struct {
		ResourceName   string `json:"resource_name"`
		ResourcePath   string `json:"resource_path"`
		ResourceMethod string `json:"resource_method"`
	}

	CreateClientDTO struct {
		ID              string       `json:"id"`
		ClientName      string       `json:"client_name" validate:"required,min=2"`
		ScopeName       string       `json:"scope_name" validate:"required,min=2"`
		Permissions     []Permission `json:"permissions"`
		KeyTimeDuration int32        `json:"key_time_duration" validate:"required,gte=1"`
		KeyPeriod       string       `json:"key_period" validate:"required,min=2"`
	}
)

func (dto *CreateClientDTO) Validate() (bool, error) {
	var validate *validator.Validate

	validate = validator.New()
	err := validate.Struct(dto)
	if err != nil {
		return false, err
	}

	if validate := validatePermissions(dto.Permissions); len(dto.Permissions) > 0 && validate == false {
		return false, fmt.Errorf("Error:field permissions has been invalid\n")
	}

	return true, nil
}

func (dto *CreateClientDTO) ToBytes() ([]byte, error) {
	b, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func validatePermissions(permissions []Permission) bool {
	for _, val := range permissions {
		if val.ResourceName == "" || val.ResourceMethod == "" || val.ResourcePath == "" {
			return false
		}

		if !strings.Contains(val.ResourcePath, "/") {
			return false
		}

		if resource.GetMethodByString(val.ResourceMethod) == resource.Unknown {
			return false
		}
	}

	return true
}
