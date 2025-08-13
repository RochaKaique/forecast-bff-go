package coordinates

import (
	"net/http"

	"github.com/spf13/viper"
)

type Coordinates struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
}

type CoordinatesClient struct {
	Conf   *viper.Viper
	Client *http.Client
}
