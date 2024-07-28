// Package genapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by unknown module path version unknown version DO NOT EDIT.
package genapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// CreateLocationRequest defines model for CreateLocationRequest.
type CreateLocationRequest struct {
	// Latitude latitude of the location
	Latitude float64 `json:"latitude"`

	// Longitude longitude of the location
	Longitude float64 `json:"longitude"`

	// Name location name
	Name string `json:"name"`
}

// CreateLocationResponse defines model for CreateLocationResponse.
type CreateLocationResponse struct {
	Data *Location `json:"data,omitempty"`

	// Message message of this operation
	Message *string `json:"message,omitempty"`

	// Status status of location creation
	Status *string `json:"status,omitempty"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	// ErrorCode error code
	ErrorCode *string `json:"errorCode,omitempty"`

	// ErrorMessage error message
	ErrorMessage *string `json:"errorMessage,omitempty"`

	// HttpCode http status error code
	HttpCode *int `json:"httpCode,omitempty"`
}

// FindLocationsCoordinates defines model for FindLocationsCoordinates.
type FindLocationsCoordinates struct {
	Items *[]Location `json:"items,omitempty"`
}

// ListLocations defines model for ListLocations.
type ListLocations struct {
	Items *[]Location `json:"items,omitempty"`
}

// Location defines model for Location.
type Location struct {
	// ID ID of location
	ID *int32 `json:"ID,omitempty"`

	// City location name
	City *string `json:"city,omitempty"`

	// Latitude latitude of the location
	Latitude *float64 `json:"latitude,omitempty"`

	// Longitude longitude of the location
	Longitude *float64 `json:"longitude,omitempty"`

	// RefId original reference id if any
	RefId *string `json:"ref_id,omitempty"`

	// Temperature temperatur of the location
	Temperature *float32 `json:"temperature,omitempty"`

	// WeatherSummary weather summary
	WeatherSummary *string `json:"weather_summary,omitempty"`

	// WindAngle angle of wind
	WindAngle *float32 `json:"wind_angle,omitempty"`

	// WindDirection direction of wind
	WindDirection *string `json:"wind_direction,omitempty"`

	// WindSpeed speed of wind in the location
	WindSpeed *float32 `json:"wind_speed,omitempty"`
}

// UpdateLocationWeatherResponse defines model for UpdateLocationWeatherResponse.
type UpdateLocationWeatherResponse struct {
	Data *struct {
		NumUpdated *int `json:"numUpdated,omitempty"`
	} `json:"data,omitempty"`

	// Message message of this operation
	Message *string `json:"message,omitempty"`

	// Status status of location creation
	Status *string `json:"status,omitempty"`
}

// ListLocationsEitherProvideIdOrNameParams defines parameters for ListLocationsEitherProvideIdOrName.
type ListLocationsEitherProvideIdOrNameParams struct {
	// Id id of the location
	Id *int32 `form:"id,omitempty" json:"id,omitempty"`

	// City name of the city
	City *string `form:"city,omitempty" json:"city,omitempty"`
}

// GetLocationsCoordinatesParams defines parameters for GetLocationsCoordinates.
type GetLocationsCoordinatesParams struct {
	// Label find by labels. Depends on the implementation
	Label string `form:"label" json:"label"`
}

// UpdateLocationWeatherJSONBody defines parameters for UpdateLocationWeather.
type UpdateLocationWeatherJSONBody interface{}

// CreateLocationJSONRequestBody defines body for CreateLocation for application/json ContentType.
type CreateLocationJSONRequestBody = CreateLocationRequest

// UpdateLocationWeatherJSONRequestBody defines body for UpdateLocationWeather for application/json ContentType.
type UpdateLocationWeatherJSONRequestBody UpdateLocationWeatherJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// list locations
	// (GET /locations)
	ListLocationsEitherProvideIdOrName(c *gin.Context, params ListLocationsEitherProvideIdOrNameParams)
	// Create a new location
	// (POST /locations)
	CreateLocation(c *gin.Context)
	// Find location coordinates
	// (GET /locations/coordinates)
	GetLocationsCoordinates(c *gin.Context, params GetLocationsCoordinatesParams)
	// Update location weathers
	// (POST /locations/weathers)
	UpdateLocationWeather(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// ListLocationsEitherProvideIdOrName operation middleware
func (siw *ServerInterfaceWrapper) ListLocationsEitherProvideIdOrName(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params ListLocationsEitherProvideIdOrNameParams

	// ------------- Optional query parameter "id" -------------

	err = runtime.BindQueryParameter("form", true, false, "id", c.Request.URL.Query(), &params.Id)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "city" -------------

	err = runtime.BindQueryParameter("form", true, false, "city", c.Request.URL.Query(), &params.City)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter city: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.ListLocationsEitherProvideIdOrName(c, params)
}

// CreateLocation operation middleware
func (siw *ServerInterfaceWrapper) CreateLocation(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.CreateLocation(c)
}

// GetLocationsCoordinates operation middleware
func (siw *ServerInterfaceWrapper) GetLocationsCoordinates(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetLocationsCoordinatesParams

	// ------------- Required query parameter "label" -------------

	if paramValue := c.Query("label"); paramValue != "" {

	} else {
		siw.ErrorHandler(c, fmt.Errorf("Query argument label is required, but not found"), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "label", c.Request.URL.Query(), &params.Label)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter label: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetLocationsCoordinates(c, params)
}

// UpdateLocationWeather operation middleware
func (siw *ServerInterfaceWrapper) UpdateLocationWeather(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.UpdateLocationWeather(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/locations", wrapper.ListLocationsEitherProvideIdOrName)
	router.POST(options.BaseURL+"/locations", wrapper.CreateLocation)
	router.GET(options.BaseURL+"/locations/coordinates", wrapper.GetLocationsCoordinates)
	router.POST(options.BaseURL+"/locations/weathers", wrapper.UpdateLocationWeather)
}
