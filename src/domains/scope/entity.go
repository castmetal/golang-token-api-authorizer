package scope

import (
	"time"

	"github.com/castmetal/golang-token-api-authorizer/src/domains/common"
	"github.com/google/uuid"
)

type Scope struct {
	common.EntityBase `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	ID                uuid.UUID       `json:"id" bson:"_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	ScopeName         string          `json:"name" gorm:"type:varchar(60);column:scope_name"`
	ScopeCreatedAt    common.JsonTime `json:"scope_created_at" gorm:"column:scope_created_at"`
	ScopeUpdatedAt    common.JsonTime `json:"scope_updated_at" gorm:"column:scope_updated_at"`
}

type ScopeProps struct {
	ID        string `json:"id"`
	ScopeName string `json:"scope_name"`
}

func (e *Scope) TableName() string {
	return "scope"
}

func (e *Scope) BeforeCreate() error {
	id := uuid.New()

	e.ID = id

	return nil
}

func NewScopeEntity(props ScopeProps) (*Scope, error) {
	var scope *Scope

	abstractEntity := common.NewAbstractEntity(props.ID)

	if common.IsNullOrEmpty(props.ScopeName) {
		return nil, common.IsNullOrEmptyError("Scope_name")
	}

	actualDate := time.Now()

	scope = &Scope{
		ScopeName:      props.ScopeName,
		ScopeCreatedAt: common.JsonTime{Value: actualDate},
		ScopeUpdatedAt: common.JsonTime{Value: actualDate},
	}

	scope.ID = abstractEntity.ID

	return scope, nil
}
