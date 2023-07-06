package resource

import (
	"strings"
	"time"

	"github.com/castmetal/golang-token-api-authorizer/src/domains/common"
	"github.com/google/uuid"
)

type Method int64

const (
	Get Method = iota + 1
	Put
	Post
	Delete
	Unknown
)

func GetMethodByString(s string) Method {
	sLower := strings.ToLower(s)

	switch sLower {
	case "get":
		return Get
	case "put":
		return Put
	case "post":
		return Post
	case "delete":
		return Delete
	}

	return Unknown
}

func (m Method) String() string {
	switch m {
	case Get:
		return "GET"
	case Put:
		return "PUT"
	case Post:
		return "POST"
	case Delete:
		return "DELETE"
	}

	return "GET"
}

type Resource struct {
	common.EntityBase `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	ID                uuid.UUID       `json:"id" bson:"_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	ResourceName      string          `json:"resource_name" gorm:"type:varchar(60);column:resource_name"`
	ResourcePath      string          `json:"resource_path" gorm:"type:varchar(255);column:resource_path"`
	ResourceMethod    Method          `json:"resource_method" gorm:"type:resource_method;column:resource_method"`
	ResourceCreatedAt common.JsonTime `json:"resource_created_at" gorm:"column:resource_created_at"`
	ResourceUpdatedAt common.JsonTime `json:"resource_updated_at" gorm:"column:resource_updated_at"`
}

type ResourceProps struct {
	ID             string `json:"id"`
	ResourceName   string `json:"resource_name"`
	ResourcePath   string `json:"resource_path"`
	ResourceMethod string `json:"resource_method"`
}

func (e *Resource) TableName() string {
	return "resource"
}

func (e *Resource) BeforeCreate() error {
	id := uuid.New()

	e.ID = id

	return nil
}

func NewResourceEntity(props ResourceProps) (*Resource, error) {
	var resource *Resource

	abstractEntity := common.NewAbstractEntity(props.ID)

	if common.IsNullOrEmpty(props.ResourceName) {
		return nil, common.IsNullOrEmptyError("resource_name")
	}
	if common.IsNullOrEmpty(props.ResourcePath) {
		return nil, common.IsNullOrEmptyError("resource_path")
	}
	if common.IsNullOrEmpty(props.ResourceMethod) {
		return nil, common.IsNullOrEmptyError("resource_method")
	}

	actualDate := time.Now()

	resource = &Resource{
		ResourceName:      props.ResourceName,
		ResourcePath:      props.ResourcePath,
		ResourceMethod:    GetMethodByString(props.ResourceMethod),
		ResourceCreatedAt: common.JsonTime{Value: actualDate},
		ResourceUpdatedAt: common.JsonTime{Value: actualDate},
	}

	resource.ID = abstractEntity.ID

	return resource, nil
}
