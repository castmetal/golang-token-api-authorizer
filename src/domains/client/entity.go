package client

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/castmetal/golang-token-api-authorizer/src/domains/common"
	"github.com/google/uuid"
)

const SALT_DEFAULT_SIZE = 128 // 128 bits security

type Permission struct {
	ResourceName   string `json:"resource_name"`
	ResourcePath   string `json:"resource_path"`
	ResourceMethod string `json:"resource_method"`
}

type Client struct {
	common.EntityBase `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	ID                uuid.UUID            `json:"id" bson:"_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	ClientName        string               `json:"client_name" gorm:"type:varchar(60);column:client_name"`
	ScopeId           uuid.UUID            `json:"scope_id" gorm:"type:uuid;column:scoipe_id"`
	Permissions       []Permission         `json:"permissions" gorm:"type:varchar(60);column:permissions"`
	ApiId             string               `json:"api_id" gorm:"type:varchar(255);column:api_id"`
	Salt              string               `json:"salt" gorm:"type:varchar(400);column:salt"`
	ScopeName         string               `json:"scope_name"`
	KeyTimeDuration   int32                `json:"key_time_duration" gorm:"type:int(3);column:key_time_duration"`
	KeyPeriod         string               `json:"key_period" gorm:"type:varchar(60);column:key_period"`
	ClientCreatedAt   common.JsonTime      `json:"client_created_at" gorm:"column:created_at"`
	ClientUpdatedAt   common.JsonTime      `json:"client_updated_at" gorm:"column:updated_at"`
	ClientDeletedAt   *common.JsonNullTime `json:"client_deleted_at" gorm:"column:deleted_at"`
}

type ClientProps struct {
	ID              string       `json:"id"`
	ClientName      string       `json:"client_name"`
	ScopeName       string       `json:"scope_name"`
	ScopeId         uuid.UUID    `json:"scope_id"`
	Permissions     []Permission `json:"permissions"`
	ApiId           string       `json:"api_id"`
	Salt            string       `json:"salt"`
	KeyTimeDuration int32        `json:"key_time_duration"`
	KeyPeriod       string       `json:"key_period"`
}

func (e *Client) TableName() string {
	return "client"
}

func (e *Client) BeforeCreate() error {
	id := uuid.New()

	e.ID = id

	return nil
}

func NewClientEntity(props ClientProps) (*Client, error) {
	var client *Client
	var salt string
	var apiId string

	abstractEntity := common.NewAbstractEntity(props.ID)

	if common.IsNullOrEmpty(props.ClientName) {
		return nil, common.IsNullOrEmptyError("client_name")
	}

	if props.ApiId != "" {
		apiId = props.ApiId
	} else {
		apiId = CreateNewApiId()
	}

	if props.Salt != "" {
		salt = props.Salt
	} else {
		salt = CreateNewSalt(SALT_DEFAULT_SIZE)
	}

	actualDate := time.Now()

	client = &Client{
		ClientName:      props.ClientName,
		ScopeId:         props.ScopeId,
		ScopeName:       props.ScopeName,
		Permissions:     props.Permissions,
		Salt:            salt,
		ApiId:           apiId,
		KeyTimeDuration: props.KeyTimeDuration,
		KeyPeriod:       props.KeyPeriod,
		ClientCreatedAt: common.JsonTime{Value: actualDate},
		ClientUpdatedAt: common.JsonTime{Value: actualDate},
	}

	client.ID = abstractEntity.ID

	return client, nil
}

func (c *Client) GetClientHashKey() string {
	return GetHashKey(c.ApiId, c.ID.String())
}

func CreateNewApiId() string {
	newUlid := GetNewUlid()

	apiId := fmt.Sprintf("CF-%s", newUlid.String())
	return apiId
}

func CreateNewSalt(saltSize int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ@.<>,.;:?/][{}`^~]|'`")

	salt := make([]rune, saltSize)
	for i := range salt {
		salt[i] = letters[rand.Intn(len(letters))]
	}

	newUlid := GetNewUlid()

	return fmt.Sprintf("%s-@;!.-%s", string(salt), newUlid.String())
}
