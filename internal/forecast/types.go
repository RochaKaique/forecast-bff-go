package forecast

import (
	"net/http"

	"github.com/spf13/viper"
)

type ForecastClient struct {
	Conf   *viper.Viper
	Client *http.Client
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
	CurrentTemperature  float64           `json:"current_temperature"`
	HighestTemperature  float64           `json:"highest_temperature"`
	LowestTemperature   float64           `json:"lowest_temperature"`
	ApparentTemperature float64           `json:"apparent_temperature"`
	Precipitation       float64           `json:"precipitation"`
	NextDayForecasts    []NextDayForecast `json:"next_day_forecasts,omitempty"`
}

type NextDayForecast struct {
	Date              string  `json:"date"`
	MinTemperature    float64 `json:"min_temperature"`
	MaxTemperature    float64 `json:"max_temperature"`
	PrecipitationProb float64 `json:"precipitation_prob,omitempty"`
}
