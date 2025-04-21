package http

import (
	"pvz_service/internal/domain"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
)

// Типы для запросов/ответов API
type (
	Error struct {
		Message string `json:"message"`
	}

	Token string

	User struct {
		Id    *types.UUID 		`json:"id,omitempty"`
		Email string      		`json:"email"`
		Role  domain.UserRole 	`json:"role"`
	}

	PVZ struct {
		City string `json:"city"`
	}

	PostDummyLoginJSONRequestBody struct {
		Role  domain.UserRole `json:"role"`
	}

	PostRegisterJSONRequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role  domain.UserRole `json:"role"`
	}

	PostReceptionsJSONRequestBody struct {
		PvzId types.UUID `json:"pvzId"`
	}

	PostProductsJSONRequestBody struct {
		PvzId types.UUID `json:"pvzId"`
		Type  string     `json:"type"`
	}
)

// Временная переменная для примера
var userID = types.UUID(uuid.New())

