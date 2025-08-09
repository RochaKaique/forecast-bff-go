package coordinates

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

func NewCooordinatesClient(client *http.Client, config *viper.Viper) *CoordinatesClient {
	return &CoordinatesClient{
		conf:   config,
		client: client,
	}
}

func (cc CoordinatesClient) GetCoordinates(ctx context.Context, postalCode string) (Coordinates, error) {
	baseUrl := cc.conf.GetString("coordinates.url")
	path := "/search"

	params := url.Values{}
	params.Add("postalCode", postalCode)
	params.Add("country", "USA")
	params.Add("format", "JSON")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseUrl+path+"?"+params.Encode(), nil)
	if err != nil {
		slog.ErrorContext(ctx, "Erro ao montar requisição http")
		return Coordinates{}, err
	}

	resp, err := cc.client.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "Erro ao realizar requisição http")
		return Coordinates{}, err
	}
	defer resp.Body.Close()

	var coordinates Coordinates
	if err := json.NewDecoder(resp.Body).Decode(&coordinates); err != nil {
		return Coordinates{}, err
	}

	return Coordinates{}, nil
}
