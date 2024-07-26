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
		Name: ptr.ToStr(req.Name),
		Lat:  ptr.ToFloat64(req.Latitude),
		Lon:  ptr.ToFloat64(req.Longitude),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, &genapi.GeneralRequestError{
			Detail: err.Error(),
		})

		return
	}

	c.JSON(200, &genapi.CreateLocationResponse{
		Data: &genapi.Location{
			ID:        ptr.ValueToPointer(int(loc.Loc.ID)),
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
		c.JSON(http.StatusInternalServerError, &genapi.GeneralRequestError{
			Detail: err.Error(),
		})
	}

	c.JSON(http.StatusOK, &genapi.UpdateLocationWeatherResponse{
		Data: &struct {
			NumUpdated *int "json:\"numUpdated,omitempty\""
		}{NumUpdated: ptr.ValueToPointer(resp.NumUpdated)},
		Status: ptr.ValueToPointer(SuccessStatus),
	})
}
