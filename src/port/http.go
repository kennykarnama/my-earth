package port

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kennykarnama/my-earth/api/openapi/genapi"
	"github.com/kennykarnama/my-earth/src/app"
	"github.com/kennykarnama/my-earth/src/pkg/ptr"
)

const (
	SuccessStatus  = "success"
	SuccessMessage = "Success create weather data"
)

type HttpHandler struct {
	locSvc *app.LocationSvc
}

func NewHttpHandler(locSvc *app.LocationSvc) *HttpHandler {
	return &HttpHandler{
		locSvc: locSvc,
	}
}

func (h *HttpHandler) CreateLocation(c *gin.Context) {
	var req genapi.CreateLocationRequest
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, &genapi.ErrorResponse{
			ErrorCode:    ptr.ValueToPointer("API-400"),
			HttpCode:     ptr.ValueToPointer(http.StatusBadRequest),
			ErrorMessage: ptr.ValueToPointer(err.Error()),
		})

		return
	}

	loc, err := h.locSvc.SaveLoc(c.Request.Context(), &app.SaveLocReq{
		Name: req.Name,
		Lat:  req.Latitude,
		Lon:  req.Longitude,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, &genapi.ErrorResponse{
			ErrorMessage: ptr.ValueToPointer(err.Error()),
			ErrorCode:    ptr.ValueToPointer("API-500"),
			HttpCode:     ptr.ValueToPointer(http.StatusInternalServerError),
		})

		return
	}

	c.JSON(200, &genapi.CreateLocationResponse{
		Data: &genapi.Location{
			ID:        ptr.ValueToPointer(loc.Loc.ID),
			City:      ptr.ValueToPointer(loc.Loc.Name),
			Latitude:  ptr.ValueToPointer(loc.Loc.Lat),
			Longitude: ptr.ValueToPointer(loc.Loc.Lon),
		},
		Message: ptr.ValueToPointer(SuccessMessage),
		Status:  ptr.ValueToPointer(SuccessStatus),
	})
}

func (h *HttpHandler) UpdateLocationWeather(c *gin.Context) {
	resp, err := h.locSvc.RefreshWeather(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, &genapi.ErrorResponse{
			ErrorMessage: ptr.ValueToPointer(err.Error()),
			ErrorCode:    ptr.ValueToPointer("API-500"),
			HttpCode:     ptr.ValueToPointer(http.StatusInternalServerError),
		})
	}

	c.JSON(http.StatusOK, &genapi.UpdateLocationWeatherResponse{
		Data: &struct {
			NumUpdated *int "json:\"numUpdated,omitempty\""
		}{NumUpdated: ptr.ValueToPointer(resp.NumUpdated)},
		Status: ptr.ValueToPointer(SuccessStatus),
	})
}

func (h *HttpHandler) GetLocationsCoordinates(c *gin.Context, params genapi.GetLocationsCoordinatesParams) {
	matches, err := h.locSvc.FindLocationsCoordinates(c.Request.Context(), params.Label)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &genapi.ErrorResponse{
			HttpCode:     ptr.ValueToPointer(http.StatusInternalServerError),
			ErrorCode:    ptr.ValueToPointer("API-500"),
			ErrorMessage: ptr.ValueToPointer(err.Error()),
		})

		return
	}

	var respMatches []genapi.Location

	resp := &genapi.FindLocationsCoordinates{
		Items: &[]genapi.Location{},
	}

	for _, m := range matches.Matches.Items {
		respMatches = append(respMatches, genapi.Location{
			City:      &m.Name,
			Latitude:  &m.Lat,
			Longitude: &m.Lon,
			RefId:     &m.RefID,
		})
	}

	resp.Items = &respMatches

	c.JSON(http.StatusOK, resp)
}

func (h *HttpHandler) ListLocationsEitherProvideIdOrName(c *gin.Context, params genapi.ListLocationsEitherProvideIdOrNameParams) {
	if params.City == nil && params.Id == nil {
		c.JSON(http.StatusBadRequest, &genapi.ErrorResponse{
			ErrorCode:    ptr.ValueToPointer("API-400"),
			ErrorMessage: ptr.ValueToPointer("either provide id or city"),
			HttpCode:     ptr.ValueToPointer(http.StatusBadRequest),
		})

		return
	}

	res, err := h.locSvc.ListLocations(c.Request.Context(), app.ListLocationsQuery{
		City: params.City,
		ID:   params.Id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, &genapi.ErrorResponse{
			ErrorCode:    ptr.ValueToPointer("API-500"),
			ErrorMessage: ptr.ValueToPointer(err.Error()),
			HttpCode:     ptr.ValueToPointer(http.StatusInternalServerError),
		})
	}

	resp := &genapi.ListLocations{
		Items: &[]genapi.Location{},
	}

	var apiLocs []genapi.Location

	for _, m := range res.Locations {
		apiLocs = append(apiLocs, genapi.Location{
			ID:             &m.ID,
			City:           &m.Name,
			Latitude:       &m.Lat,
			Longitude:      &m.Lon,
			Temperature:    &m.Temperature.Value,
			WeatherSummary: &m.Weather.Summary,
			WindAngle:      &m.WindAngle,
			WindDirection:  &m.WindDirection,
			WindSpeed:      &m.WindSpeed,
		})
	}

	resp.Items = &apiLocs

	c.JSON(http.StatusOK, resp)
}
