package coordinates

import (
	"net/http"

	"github.com/spf13/viper"
)

type Coordinates struct {
	latitude  string `json:"lat"`
	longitude string `json:"lon"`
}

type CoordinatesClient struct {
	conf   *viper.Viper
	client *http.Client
}
