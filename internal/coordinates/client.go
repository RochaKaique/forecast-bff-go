package coordinates

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

func NewCooordinatesClient(client *http.Client, config *viper.Viper) *CoordinatesClient {
	return &CoordinatesClient{
		Conf:   config,
		Client: client,
	}
}

func (cc CoordinatesClient) GetCoordinates(ctx context.Context, postalCode string) (Coordinates, error) {
	baseUrl := cc.Conf.GetString("coordinates.uri")
	path := "/search"

	params := url.Values{}
	params.Add("postalcode", postalCode)
	params.Add("country", "US")
	params.Add("format", "json")

	slog.InfoContext(ctx, baseUrl+path+"?"+params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseUrl+path+"?"+params.Encode(), nil)
	if err != nil {
		slog.ErrorContext(ctx, "Erro ao montar requisição http")
		return Coordinates{}, err
	}

	resp, err := cc.Client.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "Erro ao realizar requisição http")
		return Coordinates{}, err
	}
	defer resp.Body.Close()

	var coordinates []Coordinates
	if err := json.NewDecoder(resp.Body).Decode(&coordinates); err != nil {
		slog.ErrorContext(ctx, "Erro ao serializar resposta")
		return Coordinates{}, err
	}
	fmt.Println(coordinates)
	return coordinates[0], nil
}
