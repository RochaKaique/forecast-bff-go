package coordinates

import (
	"net/http"

	"github.com/spf13/viper"
)

type Doer interface {
    Do(req *http.Request) (*http.Response, error)
}

type Coordinates struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
}

type CoordinatesClient struct {
	Conf   *viper.Viper
	Client Doer
}
