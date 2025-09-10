package forecast

import (
	"net/http"

	"github.com/spf13/viper"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type ForecastClient struct {
	Conf   *viper.Viper
	Client Doer
}

type ForecastRequest struct {
	Latitude  string `json:"latitude" validate:"required"`
	Longitude string `json:"longitude" validate:"required"`
	StartDate string `json:"startDate" validate:"omitempty,datetime=2006-01-02"`
	EndDate   string `json:"endDate" validate:"omitempty,datetime=2006-01-02"`
	Units     string `json:"units" validate:"omitempty,oneof=metric imperial"`
	TempScale string `json:"tempScale" validate:"omitempty,oneof=C F"`
}

type ForecastResponse struct {
	CurrentTemperature  float32           `json:"current_temperature"`
	HighestTemperature  float32           `json:"highest_temperature"`
	LowestTemperature   float32           `json:"lowest_temperature"`
	ApparentTemperature float32           `json:"apparent_temperature"`
	Precipitation       float32           `json:"precipitation"`
	NextDayForecasts    []NextDayForecast `json:"next_day_forecasts"`
}

type NextDayForecast struct {
	Date              string  `json:"date"`
	MinTemperature    float32 `json:"min_temperature"`
	MaxTemperature    float32 `json:"max_temperature"`
	PrecipitationProb float32 `json:"precipitation_prob"`
}

type HourlyData struct {
	Time                []string  `json:"time"`
	Temperature2M       []float32 `json:"temperature_2m"`
	ApparentTemperature []float32 `json:"apparent_temperature"`
	Precipitation       []float32 `json:"precipitation"`
}

type WeatherDataResponse struct {
	Latitude  float32    `json:"latitude"`
	Longitude float32    `json:"longitude"`
	Elevation float32    `json:"elevation"`
	Hourly    HourlyData `json:"hourly"`
}
