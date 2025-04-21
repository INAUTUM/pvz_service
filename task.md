Разработайте backend-сервис для сотрудников ПВЗ, который позволит вносить информацию по заказам в рамках приёмки товаров.

Авторизация пользователей:
Используя ручку /dummyLogin и передав в неё желаемый тип пользователя (client, moderator), сервис в ответе вернёт токен с соответствующим уровнем доступа — обычного пользователя или модератора. Этот токен нужно передавать во все endpoints, требующие авторизации.

Регистрация и авторизация пользователей по почте и паролю:
При регистрации используется endpoint /register. В базе создаётся и сохраняется новый пользователь желаемого типа: обычный пользователь (client) или модератор (moderator). У созданного пользователя появляется токен endpoint /login. При успешной авторизации по почте и паролю возвращается токен для пользователя с соответствующим ему уровнем доступа.

Заведение ПВЗ:
Только пользователь с ролью «модератор» может завести ПВЗ в системе.
В случае успешного запроса возвращается полная информация о созданном ПВЗ. Заведение ПВЗ возможно только в трёх городах: Москва, Санкт-Петербург и Казань. В других городах ПВЗ завести на первых порах нельзя, в таком случае необходимо вернуть ошибку.
Результатом добавления ПВЗ должна стать новая запись в хранилище данных

Добавление информации о приёмке товаров:
Только авторизованный пользователь системы с ролью «сотрудник ПВЗ» может инициировать приём товара.
Результатом инициации приёма товаров должна стать новая запись в хранилище данных.
Если же предыдущая приёмка товара не была закрыта, то операция по созданию нового приёма товаров невозможна.

Добавление товаров в рамках одной приёмки:
Только авторизованный пользователь системы с ролью «сотрудник ПВЗ» может добавлять товары после его осмотра.
При этом товар должен привязываться к последнему незакрытому приёму товаров в рамках текущего ПВЗ.
Если же нет новой незакрытой приёмки товаров, то в таком случае должна возвращаться ошибка, и товар не должен добавляться в систему.
Если последняя приёмка товара все ещё не была закрыта, то результатом должна стать привязка товара к текущему ПВЗ и текущей приёмке с последующем добавлением данных в хранилище.

Удаление товаров в рамках не закрытой приёмки:
Только авторизованный пользователь системы с ролью «сотрудник ПВЗ» может удалять товары, которые были добавлены в рамках текущей приёмки на ПВЗ.
Удаление товара возможно только до закрытия приёмки, после этого уже невозможно изменить состав товаров, которые были приняты на ПВЗ.
Удаление товаров производится по принципу LIFO, т.е. возможно удалять товары только в том порядке, в котором они были добавлены в рамках текущей приёмки.

Закрытие приёмки:
Только авторизованный пользователь системы с ролью «сотрудник ПВЗ» может закрывать приём товаров.
В случае, если приёмка товаров уже была закрыта (или приёма товаров в данном ПВЗ ещё не было), то следует вернуть ошибку.
Во всех остальных случаях необходимо обновить данные в хранилище и зафиксировать товары, которые были в рамках этой приёмки.

Получение данных:
Только авторизованный пользователь системы с ролью «сотрудник ПВЗ» или «модератор» может получать эти данные.
Необходимо получить список ПВЗ и всю информацию по ним при помощи пагинации.
При этом добавить фильтр по дате приёмки товаров, т.е. выводить только те ПВЗ и всю информацию по ним, которые в указанный диапазон времени проводили приёмы товаров.
Общие вводные

У сущности «Пункт приёма заказов (ПВЗ)» есть:
Уникальный идентификатор
Дата регистрации в системе
Город

У сущности «Приёмка товара» есть:
Уникальный идентификатор
Дата и время проведения приёмки
ПВЗ, в котором была осуществлена приёмка
Товары, которые были приняты в рамках данной приёмки
Статус (in_progress, close)

У сущности «Товар» есть:
Уникальный идентификатор
Дата и время приёма товара (дата и время, когда товар был добавлен в систему в рамках приёмки товаров)
Тип (мы работаем с тремя типами товаров: электроника, одежда, обувь)

Условия:
Используйте этот API
Сервер должен быть запущен на порту 8080.
В параметрах запроса можно выбрать роль пользователя: модератор или обычный пользователь. В зависимости от роли будет сгенерирован токен с определённым уровнем доступа.

Нефункциональные требования:
RPS — 1000
SLI времени ответа — 100 мс
SLI успешности ответа — 99.99%
Код обязательно должен быть покрыт unit-тестами. Тестовое покрытие не менее 75%.

Должен быть разработан один интеграционный тест, который:
Первым делом создает новый ПВЗ
Добавляет новую приёмку заказов
Добавляет 50 товаров в рамках текущей приёмки заказов
Закрывает приёмку заказов
Дополнительные задания

Реализовать пользовательскую авторизацию по методам /register и /login (при этом метод /dummyLogin все равно должен быть реализован)
Реализовать gRPC-метод, который просто вернёт все добавленные в систему ПВЗ. Для него не требуется проверка авторизации и валидация ролей пользователей. Сервер для gRPC должен быть запущен на порту 3000. Обратите внимание, что в файле pvz.proto необходимо прописать go_package под вашу структуру проекта

Добавить в проект prometheus и собирать следующие метрики:
Технические:
Количество запросов
Время ответа
Бизнесовые:
Количество созданных ПВЗ
Количество созданных приёмок заказов
Количество добавленных товаров Сервер для prometheus должен быть поднят на порту 9000 и отдавать данные по ручке /metrics.
Настроить логирование в проекте
Настроить кодогенерацию DTO endpoint'ов по openapi схеме

Требования по стеку
язык: golang
База данных: PostgreSQL
Допустимо использовать билдеры для запросов, например, такой: https://github.com/Masterminds/squirrel
Для деплоя зависимостей и самого сервиса нужно использовать Docker или Docker & DockerCompose

pvz.proto:
syntax = "proto3";

package pvz.v1;

option go_package = "pvz/pvz_v1;pvz_v1";

import "google/protobuf/timestamp.proto";

service PVZService {
  rpc GetPVZList(GetPVZListRequest) returns (GetPVZListResponse);
}

message PVZ {
  string id = 1;
  google.protobuf.Timestamp registration_date = 2;
  string city = 3;
}

enum ReceptionStatus {
  RECEPTION_STATUS_IN_PROGRESS = 0;
  RECEPTION_STATUS_CLOSED = 1;
}

message GetPVZListRequest {}

message GetPVZListResponse {
  repeated PVZ pvzs = 1;
}

generated.go
// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Defines values for PVZCity.
const (
	Казань         PVZCity = "Казань"
	Москва         PVZCity = "Москва"
	СанктПетербург PVZCity = "Санкт-Петербург"
)

// Defines values for ProductType.
const (
	ProductTypeОбувь       ProductType = "обувь"
	ProductTypeОдежда      ProductType = "одежда"
	ProductTypeЭлектроника ProductType = "электроника"
)

// Defines values for ReceptionStatus.
const (
	Close      ReceptionStatus = "close"
	InProgress ReceptionStatus = "in_progress"
)

// Defines values for UserRole.
const (
	UserRoleEmployee  UserRole = "employee"
	UserRoleModerator UserRole = "moderator"
)

// Defines values for PostDummyLoginJSONBodyRole.
const (
	PostDummyLoginJSONBodyRoleEmployee  PostDummyLoginJSONBodyRole = "employee"
	PostDummyLoginJSONBodyRoleModerator PostDummyLoginJSONBodyRole = "moderator"
)

// Defines values for PostProductsJSONBodyType.
const (
	PostProductsJSONBodyTypeОбувь       PostProductsJSONBodyType = "обувь"
	PostProductsJSONBodyTypeОдежда      PostProductsJSONBodyType = "одежда"
	PostProductsJSONBodyTypeЭлектроника PostProductsJSONBodyType = "электроника"
)

// Defines values for PostRegisterJSONBodyRole.
const (
	Employee  PostRegisterJSONBodyRole = "employee"
	Moderator PostRegisterJSONBodyRole = "moderator"
)

// Error defines model for Error.
type Error struct {
	Message string `json:"message"`
}

// PVZ defines model for PVZ.
type PVZ struct {
	City             PVZCity             `json:"city"`
	Id               *openapi_types.UUID `json:"id,omitempty"`
	RegistrationDate *time.Time          `json:"registrationDate,omitempty"`
}

// PVZCity defines model for PVZ.City.
type PVZCity string

// Product defines model for Product.
type Product struct {
	DateTime    *time.Time          `json:"dateTime,omitempty"`
	Id          *openapi_types.UUID `json:"id,omitempty"`
	ReceptionId openapi_types.UUID  `json:"receptionId"`
	Type        ProductType         `json:"type"`
}

// ProductType defines model for Product.Type.
type ProductType string

// Reception defines model for Reception.
type Reception struct {
	DateTime time.Time           `json:"dateTime"`
	Id       *openapi_types.UUID `json:"id,omitempty"`
	PvzId    openapi_types.UUID  `json:"pvzId"`
	Status   ReceptionStatus     `json:"status"`
}

// ReceptionStatus defines model for Reception.Status.
type ReceptionStatus string

// Token defines model for Token.
type Token = string

// User defines model for User.
type User struct {
	Email openapi_types.Email `json:"email"`
	Id    *openapi_types.UUID `json:"id,omitempty"`
	Role  UserRole            `json:"role"`
}

// UserRole defines model for User.Role.
type UserRole string

// PostDummyLoginJSONBody defines parameters for PostDummyLogin.
type PostDummyLoginJSONBody struct {
	Role PostDummyLoginJSONBodyRole `json:"role"`
}

// PostDummyLoginJSONBodyRole defines parameters for PostDummyLogin.
type PostDummyLoginJSONBodyRole string

// PostLoginJSONBody defines parameters for PostLogin.
type PostLoginJSONBody struct {
	Email    openapi_types.Email `json:"email"`
	Password string              `json:"password"`
}

// PostProductsJSONBody defines parameters for PostProducts.
type PostProductsJSONBody struct {
	PvzId openapi_types.UUID       `json:"pvzId"`
	Type  PostProductsJSONBodyType `json:"type"`
}

// PostProductsJSONBodyType defines parameters for PostProducts.
type PostProductsJSONBodyType string

// GetPvzParams defines parameters for GetPvz.
type GetPvzParams struct {
	// StartDate Начальная дата диапазона
	StartDate *time.Time `form:"startDate,omitempty" json:"startDate,omitempty"`

	// EndDate Конечная дата диапазона
	EndDate *time.Time `form:"endDate,omitempty" json:"endDate,omitempty"`

	// Page Номер страницы
	Page *int `form:"page,omitempty" json:"page,omitempty"`

	// Limit Количество элементов на странице
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`
}

// PostReceptionsJSONBody defines parameters for PostReceptions.
type PostReceptionsJSONBody struct {
	PvzId openapi_types.UUID `json:"pvzId"`
}

// PostRegisterJSONBody defines parameters for PostRegister.
type PostRegisterJSONBody struct {
	Email    openapi_types.Email      `json:"email"`
	Password string                   `json:"password"`
	Role     PostRegisterJSONBodyRole `json:"role"`
}

// PostRegisterJSONBodyRole defines parameters for PostRegister.
type PostRegisterJSONBodyRole string

// PostDummyLoginJSONRequestBody defines body for PostDummyLogin for application/json ContentType.
type PostDummyLoginJSONRequestBody PostDummyLoginJSONBody

// PostLoginJSONRequestBody defines body for PostLogin for application/json ContentType.
type PostLoginJSONRequestBody PostLoginJSONBody

// PostProductsJSONRequestBody defines body for PostProducts for application/json ContentType.
type PostProductsJSONRequestBody PostProductsJSONBody

// PostPvzJSONRequestBody defines body for PostPvz for application/json ContentType.
type PostPvzJSONRequestBody = PVZ

// PostReceptionsJSONRequestBody defines body for PostReceptions for application/json ContentType.
type PostReceptionsJSONRequestBody PostReceptionsJSONBody

// PostRegisterJSONRequestBody defines body for PostRegister for application/json ContentType.
type PostRegisterJSONRequestBody PostRegisterJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Получение тестового токена
	// (POST /dummyLogin)
	PostDummyLogin(ctx echo.Context) error
	// Авторизация пользователя
	// (POST /login)
	PostLogin(ctx echo.Context) error
	// Добавление товара в текущую приемку (только для сотрудников ПВЗ)
	// (POST /products)
	PostProducts(ctx echo.Context) error
	// Получение списка ПВЗ с фильтрацией по дате приемки и пагинацией
	// (GET /pvz)
	GetPvz(ctx echo.Context, params GetPvzParams) error
	// Создание ПВЗ (только для модераторов)
	// (POST /pvz)
	PostPvz(ctx echo.Context) error
	// Закрытие последней открытой приемки товаров в рамках ПВЗ
	// (POST /pvz/{pvzId}/close_last_reception)
	PostPvzPvzIdCloseLastReception(ctx echo.Context, pvzId openapi_types.UUID) error
	// Удаление последнего добавленного товара из текущей приемки (LIFO, только для сотрудников ПВЗ)
	// (POST /pvz/{pvzId}/delete_last_product)
	PostPvzPvzIdDeleteLastProduct(ctx echo.Context, pvzId openapi_types.UUID) error
	// Создание новой приемки товаров (только для сотрудников ПВЗ)
	// (POST /receptions)
	PostReceptions(ctx echo.Context) error
	// Регистрация пользователя
	// (POST /register)
	PostRegister(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostDummyLogin converts echo context to params.
func (w *ServerInterfaceWrapper) PostDummyLogin(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostDummyLogin(ctx)
	return err
}

// PostLogin converts echo context to params.
func (w *ServerInterfaceWrapper) PostLogin(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostLogin(ctx)
	return err
}

// PostProducts converts echo context to params.
func (w *ServerInterfaceWrapper) PostProducts(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostProducts(ctx)
	return err
}

// GetPvz converts echo context to params.
func (w *ServerInterfaceWrapper) GetPvz(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPvzParams
	// ------------- Optional query parameter "startDate" -------------

	err = runtime.BindQueryParameter("form", true, false, "startDate", ctx.QueryParams(), &params.StartDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter startDate: %s", err))
	}

	// ------------- Optional query parameter "endDate" -------------

	err = runtime.BindQueryParameter("form", true, false, "endDate", ctx.QueryParams(), &params.EndDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter endDate: %s", err))
	}

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", ctx.QueryParams(), &params.Page)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter page: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPvz(ctx, params)
	return err
}

// PostPvz converts echo context to params.
func (w *ServerInterfaceWrapper) PostPvz(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPvz(ctx)
	return err
}

// PostPvzPvzIdCloseLastReception converts echo context to params.
func (w *ServerInterfaceWrapper) PostPvzPvzIdCloseLastReception(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "pvzId" -------------
	var pvzId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "pvzId", ctx.Param("pvzId"), &pvzId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter pvzId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPvzPvzIdCloseLastReception(ctx, pvzId)
	return err
}

// PostPvzPvzIdDeleteLastProduct converts echo context to params.
func (w *ServerInterfaceWrapper) PostPvzPvzIdDeleteLastProduct(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "pvzId" -------------
	var pvzId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "pvzId", ctx.Param("pvzId"), &pvzId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter pvzId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPvzPvzIdDeleteLastProduct(ctx, pvzId)
	return err
}

// PostReceptions converts echo context to params.
func (w *ServerInterfaceWrapper) PostReceptions(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostReceptions(ctx)
	return err
}

// PostRegister converts echo context to params.
func (w *ServerInterfaceWrapper) PostRegister(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostRegister(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/dummyLogin", wrapper.PostDummyLogin)
	router.POST(baseURL+"/login", wrapper.PostLogin)
	router.POST(baseURL+"/products", wrapper.PostProducts)
	router.GET(baseURL+"/pvz", wrapper.GetPvz)
	router.POST(baseURL+"/pvz", wrapper.PostPvz)
	router.POST(baseURL+"/pvz/:pvzId/close_last_reception", wrapper.PostPvzPvzIdCloseLastReception)
	router.POST(baseURL+"/pvz/:pvzId/delete_last_product", wrapper.PostPvzPvzIdDeleteLastProduct)
	router.POST(baseURL+"/receptions", wrapper.PostReceptions)
	router.POST(baseURL+"/register", wrapper.PostRegister)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}

pvz_grpc.pb.go
// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: api/pvz.proto

package pvz_v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	PVZService_GetPVZList_FullMethodName = "/pvz.v1.PVZService/GetPVZList"
)

// PVZServiceClient is the client API for PVZService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PVZServiceClient interface {
	GetPVZList(ctx context.Context, in *GetPVZListRequest, opts ...grpc.CallOption) (*GetPVZListResponse, error)
}

type pVZServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPVZServiceClient(cc grpc.ClientConnInterface) PVZServiceClient {
	return &pVZServiceClient{cc}
}

func (c *pVZServiceClient) GetPVZList(ctx context.Context, in *GetPVZListRequest, opts ...grpc.CallOption) (*GetPVZListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPVZListResponse)
	err := c.cc.Invoke(ctx, PVZService_GetPVZList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PVZServiceServer is the server API for PVZService service.
// All implementations must embed UnimplementedPVZServiceServer
// for forward compatibility.
type PVZServiceServer interface {
	GetPVZList(context.Context, *GetPVZListRequest) (*GetPVZListResponse, error)
	mustEmbedUnimplementedPVZServiceServer()
}

// UnimplementedPVZServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPVZServiceServer struct{}

func (UnimplementedPVZServiceServer) GetPVZList(context.Context, *GetPVZListRequest) (*GetPVZListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPVZList not implemented")
}
func (UnimplementedPVZServiceServer) mustEmbedUnimplementedPVZServiceServer() {}
func (UnimplementedPVZServiceServer) testEmbeddedByValue()                    {}

// UnsafePVZServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PVZServiceServer will
// result in compilation errors.
type UnsafePVZServiceServer interface {
	mustEmbedUnimplementedPVZServiceServer()
}

func RegisterPVZServiceServer(s grpc.ServiceRegistrar, srv PVZServiceServer) {
	// If the following call pancis, it indicates UnimplementedPVZServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PVZService_ServiceDesc, srv)
}

func _PVZService_GetPVZList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPVZListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PVZServiceServer).GetPVZList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PVZService_GetPVZList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PVZServiceServer).GetPVZList(ctx, req.(*GetPVZListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PVZService_ServiceDesc is the grpc.ServiceDesc for PVZService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PVZService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pvz.v1.PVZService",
	HandlerType: (*PVZServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPVZList",
			Handler:    _PVZService_GetPVZList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/pvz.proto",
}

pvz_pb.go
// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: api/pvz.proto

package pvz_v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ReceptionStatus int32

const (
	ReceptionStatus_RECEPTION_STATUS_IN_PROGRESS ReceptionStatus = 0
	ReceptionStatus_RECEPTION_STATUS_CLOSED      ReceptionStatus = 1
)

// Enum value maps for ReceptionStatus.
var (
	ReceptionStatus_name = map[int32]string{
		0: "RECEPTION_STATUS_IN_PROGRESS",
		1: "RECEPTION_STATUS_CLOSED",
	}
	ReceptionStatus_value = map[string]int32{
		"RECEPTION_STATUS_IN_PROGRESS": 0,
		"RECEPTION_STATUS_CLOSED":      1,
	}
)

func (x ReceptionStatus) Enum() *ReceptionStatus {
	p := new(ReceptionStatus)
	*p = x
	return p
}

func (x ReceptionStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ReceptionStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_api_pvz_proto_enumTypes[0].Descriptor()
}

func (ReceptionStatus) Type() protoreflect.EnumType {
	return &file_api_pvz_proto_enumTypes[0]
}

func (x ReceptionStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ReceptionStatus.Descriptor instead.
func (ReceptionStatus) EnumDescriptor() ([]byte, []int) {
	return file_api_pvz_proto_rawDescGZIP(), []int{0}
}

type PVZ struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	Id               string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	RegistrationDate *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=registration_date,json=registrationDate,proto3" json:"registration_date,omitempty"`
	City             string                 `protobuf:"bytes,3,opt,name=city,proto3" json:"city,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *PVZ) Reset() {
	*x = PVZ{}
	mi := &file_api_pvz_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PVZ) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PVZ) ProtoMessage() {}

func (x *PVZ) ProtoReflect() protoreflect.Message {
	mi := &file_api_pvz_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PVZ.ProtoReflect.Descriptor instead.
func (*PVZ) Descriptor() ([]byte, []int) {
	return file_api_pvz_proto_rawDescGZIP(), []int{0}
}

func (x *PVZ) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PVZ) GetRegistrationDate() *timestamppb.Timestamp {
	if x != nil {
		return x.RegistrationDate
	}
	return nil
}

func (x *PVZ) GetCity() string {
	if x != nil {
		return x.City
	}
	return ""
}

type GetPVZListRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPVZListRequest) Reset() {
	*x = GetPVZListRequest{}
	mi := &file_api_pvz_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPVZListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPVZListRequest) ProtoMessage() {}

func (x *GetPVZListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_pvz_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPVZListRequest.ProtoReflect.Descriptor instead.
func (*GetPVZListRequest) Descriptor() ([]byte, []int) {
	return file_api_pvz_proto_rawDescGZIP(), []int{1}
}

type GetPVZListResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Pvzs          []*PVZ                 `protobuf:"bytes,1,rep,name=pvzs,proto3" json:"pvzs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPVZListResponse) Reset() {
	*x = GetPVZListResponse{}
	mi := &file_api_pvz_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPVZListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPVZListResponse) ProtoMessage() {}

func (x *GetPVZListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_pvz_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPVZListResponse.ProtoReflect.Descriptor instead.
func (*GetPVZListResponse) Descriptor() ([]byte, []int) {
	return file_api_pvz_proto_rawDescGZIP(), []int{2}
}

func (x *GetPVZListResponse) GetPvzs() []*PVZ {
	if x != nil {
		return x.Pvzs
	}
	return nil
}

var File_api_pvz_proto protoreflect.FileDescriptor

const file_api_pvz_proto_rawDesc = "" +
	"\n" +
	"\rapi/pvz.proto\x12\x06pvz.v1\x1a\x1fgoogle/protobuf/timestamp.proto\"r\n" +
	"\x03PVZ\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12G\n" +
	"\x11registration_date\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\x10registrationDate\x12\x12\n" +
	"\x04city\x18\x03 \x01(\tR\x04city\"\x13\n" +
	"\x11GetPVZListRequest\"5\n" +
	"\x12GetPVZListResponse\x12\x1f\n" +
	"\x04pvzs\x18\x01 \x03(\v2\v.pvz.v1.PVZR\x04pvzs*P\n" +
	"\x0fReceptionStatus\x12 \n" +
	"\x1cRECEPTION_STATUS_IN_PROGRESS\x10\x00\x12\x1b\n" +
	"\x17RECEPTION_STATUS_CLOSED\x10\x012Q\n" +
	"\n" +
	"PVZService\x12C\n" +
	"\n" +
	"GetPVZList\x12\x19.pvz.v1.GetPVZListRequest\x1a\x1a.pvz.v1.GetPVZListResponseB\x13Z\x11pvz/pvz_v1;pvz_v1b\x06proto3"

var (
	file_api_pvz_proto_rawDescOnce sync.Once
	file_api_pvz_proto_rawDescData []byte
)

func file_api_pvz_proto_rawDescGZIP() []byte {
	file_api_pvz_proto_rawDescOnce.Do(func() {
		file_api_pvz_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_pvz_proto_rawDesc), len(file_api_pvz_proto_rawDesc)))
	})
	return file_api_pvz_proto_rawDescData
}

var file_api_pvz_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_pvz_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_pvz_proto_goTypes = []any{
	(ReceptionStatus)(0),          // 0: pvz.v1.ReceptionStatus
	(*PVZ)(nil),                   // 1: pvz.v1.PVZ
	(*GetPVZListRequest)(nil),     // 2: pvz.v1.GetPVZListRequest
	(*GetPVZListResponse)(nil),    // 3: pvz.v1.GetPVZListResponse
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_api_pvz_proto_depIdxs = []int32{
	4, // 0: pvz.v1.PVZ.registration_date:type_name -> google.protobuf.Timestamp
	1, // 1: pvz.v1.GetPVZListResponse.pvzs:type_name -> pvz.v1.PVZ
	2, // 2: pvz.v1.PVZService.GetPVZList:input_type -> pvz.v1.GetPVZListRequest
	3, // 3: pvz.v1.PVZService.GetPVZList:output_type -> pvz.v1.GetPVZListResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_pvz_proto_init() }
func file_api_pvz_proto_init() {
	if File_api_pvz_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_pvz_proto_rawDesc), len(file_api_pvz_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_pvz_proto_goTypes,
		DependencyIndexes: file_api_pvz_proto_depIdxs,
		EnumInfos:         file_api_pvz_proto_enumTypes,
		MessageInfos:      file_api_pvz_proto_msgTypes,
	}.Build()
	File_api_pvz_proto = out.File
	file_api_pvz_proto_goTypes = nil
	file_api_pvz_proto_depIdxs = nil
}

.
├── api
│   ├── generated.go
│   ├── openapi.yaml
│   └── pvz.proto
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── pkg
│   ├── database
│   ├── jwt
│   └── logger
├── pvz
│   └── pvz_v1
│       ├── pvz_grpc.pb.go
│       └── pvz.pb.go
├── README.md
└── task.md

реализуй проект на golang