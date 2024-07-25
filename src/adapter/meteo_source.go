package adapter

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kennykarnama/my-earth/api/openapi/genapi"
	"github.com/kennykarnama/my-earth/src/domain"
	"github.com/kennykarnama/my-earth/src/pkg/ptr"
)

type ErrGeneralMeteoSource struct {
	genapi.GeneralRequestError
	HttpCode int
}

func (e *ErrGeneralMeteoSource) Error() string {
	return e.Detail
}

type ErrHttpValidationMeteoSource struct {
	genapi.HTTPValidationError
	HttpCode int
}

func (e *ErrHttpValidationMeteoSource) Error() string {
	var errMessage []string
	if e.Detail != nil {
		for _, emsg := range *e.Detail {
			errMessage = append(errMessage, fmt.Sprintf("t:%v-v:%v", emsg.Type, emsg.Msg))
		}
	}

	return strings.Join(errMessage, "\n")
}

type MeteoSource struct {
	cli    genapi.ClientWithResponsesInterface
	apiKey string
}

func NewMeteoSource(cli genapi.ClientWithResponsesInterface, apiKey string) *MeteoSource {
	return &MeteoSource{
		cli:    cli,
		apiKey: apiKey,
	}
}

func (m *MeteoSource) FindByPoint(ctx context.Context, lat, lon float64) (*domain.WeatherByPoint, error) {
	resp, err := m.cli.PointPointGetWithResponse(ctx, &genapi.PointPointGetParams{
		Lat:      ptr.ValueToPointer(fmt.Sprintf("%v", lat)),
		Lon:      ptr.ValueToPointer(fmt.Sprintf("%v", lon)),
		Sections: ptr.ValueToPointer("current"),
		Key:      ptr.ValueToPointer(m.apiKey),
		Timezone: ptr.ValueToPointer("utc"),
	})
	if err != nil {
		return nil, fmt.Errorf("findByPoint err: %w", err)
	}

	if resp.JSON400 != nil {
		return nil, &ErrGeneralMeteoSource{*resp.JSON400, http.StatusBadRequest}
	}

	if resp.JSON402 != nil {
		return nil, &ErrGeneralMeteoSource{*resp.JSON402, 402}
	}

	if resp.JSON403 != nil {
		return nil, &ErrGeneralMeteoSource{*resp.JSON403, http.StatusForbidden}
	}

	if resp.JSON422 != nil {
		return nil, &ErrHttpValidationMeteoSource{*resp.JSON422, http.StatusUnprocessableEntity}
	}

	if resp.JSON429 != nil {
		return nil, &ErrGeneralMeteoSource{*resp.JSON429, 429}
	}

	if resp.JSON200 == nil || resp.JSON200.Current == nil {
		return nil, fmt.Errorf("empty response")
	}

	// use  current
	current := resp.JSON200.Current

	w := domain.WeatherByPoint{
		Lat: lat,
		Lon: lon,
		Weather: domain.Weather{
			Summary: ptr.ToStr(current.Summary),
			Temperature: domain.Temperature{
				Value: ptr.ToFloat32(current.Temperature),
				Unit:  resp.JSON200.Units,
			},
			WindSpeed:     ptr.ToFloat32(current.Wind.Speed),
			WindAngle:     ptr.ToFloat32(current.Wind.Angle),
			WindDirection: ptr.ToStr(current.Wind.Dir),
		},
	}

	tz, err := time.LoadLocation(strings.ToUpper(ptr.ToStr(resp.JSON200.Timezone)))
	if err != nil {
		return nil, fmt.Errorf("findByPoint.parseTimeLocation err: %w", err)
	}

	headers := resp.HTTPResponse.Header
	for hk, hv := range headers {
		if hk == "Expires" {
			if len(hv) == 0 {
				return nil, fmt.Errorf("findByPoint err: %w", fmt.Errorf("missing expires header"))
			}

			log.Println("found expires header with value", hv[0])

			err = w.SetExpiredAt(hv[0], tz, time.RFC1123)
			if err != nil {
				return nil, fmt.Errorf("findByPoint err: %w", err)
			}
		}
	}

	return &w, nil
}
