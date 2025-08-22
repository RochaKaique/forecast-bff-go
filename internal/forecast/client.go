package forecast

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/spf13/viper"
)

func NewForecastClient(client *http.Client, config *viper.Viper) *ForecastClient {
	return &ForecastClient{
		Conf:   config,
		Client: client,
	}
}

func (fc ForecastClient) GetForecast(ctx context.Context, request *ForecastRequest) (ForecastResponse, error) {
	baseUrl := fc.Conf.GetString("forecast.uri")
	path := "/v1/forecast"
	params := request.provideUriParams()
	slog.DebugContext(ctx, "Api de forcast = "+baseUrl+path+params)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseUrl+path+params, nil)
	req.Header.Set("Accept", "application/json")
	if err != nil {
		slog.ErrorContext(ctx, "Erro ao montar requisição http", "err", err)
		return ForecastResponse{}, err
	}

	resp, err := fc.Client.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "Erro ao realizar requisição http", "err", err)
		return ForecastResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		slog.ErrorContext(ctx, "Resposta da api fora do esperado: "+resp.Status)
		return ForecastResponse{}, errors.New("Resposta da api fora do esperado")
	}

	var weather WeatherDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		slog.ErrorContext(ctx, "Erro ao serializar resposta", "err", err)
		return ForecastResponse{}, err
	}

	forecastResponse := ForecastResponse{
		CurrentTemperature:  weather.ActualTemp(),
		HighestTemperature:  weather.DayHighestTemp(),
		LowestTemperature:   weather.DayLowestTemp(),
		ApparentTemperature: weather.ApparentTemperatureNow(),
		Precipitation:       weather.PrecipitationNow(),
		NextDayForecasts:    weather.NextDaysForecast(),
	}

	return forecastResponse, nil
}

func (req ForecastRequest) provideUriParams() string {

	var validTempUnits = []string{"celsius", "fahrenheit"}
	var validPrecipitationUnits = []string{"mm", "inch"}

	now := time.Now().UTC()
	defaultDate := now.Format("2006-01-02")

	var startDt, endDt, tempUnit, precipitationUnit string

	if req.StartDate == "" {
		startDt = defaultDate
	} else {
		startDt = req.StartDate
	}

	if req.EndDate == "" { // em Go, não há diferença entre isEmpty/isBlank
		endDt = startDt
	} else {
		endDt = req.EndDate
	}

	if slices.Contains(validTempUnits, req.TempScale) {
		tempUnit = req.TempScale
	} else {
		tempUnit = "fahrenheit"
	}

	if slices.Contains(validPrecipitationUnits, req.Units) {
		if req.Units == "metrics" {
			precipitationUnit = "mm"
		} else {
			precipitationUnit = "inch"
		}
	} else {
		precipitationUnit = "inch"
	}

	params := fmt.Sprintf("?latitude=%s&longitude=%s&hourly=temperature_2m,precipitation,apparent_temperature&start_date=%s&end_date=%s&temperature_unit=%s&precipitation_unit=%s", req.Latitude, req.Longitude, startDt, endDt, tempUnit, precipitationUnit)

	return params
}
