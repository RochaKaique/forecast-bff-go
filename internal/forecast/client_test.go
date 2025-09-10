package forecast

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

type fakeHTTPClient struct {
	resp *http.Response
	err  error
}

func (f *fakeHTTPClient) Do(*http.Request) (*http.Response, error) {
	return f.resp, f.err
}

func TestForecastClient_GetForecast(t *testing.T) {
	testCases := []struct {
		name    string
		client  *fakeHTTPClient
		config  *viper.Viper
		request *ForecastRequest
		wantErr string
	}{
		{
			name:   "erro no transporte (Client.Do falha)",
			client: &fakeHTTPClient{err: errors.New("connection refused")},
			config: viperWithURI("http://fake-api"),
			request: &ForecastRequest{
				Latitude:  "",
				Longitude: "",
				StartDate: "",
				EndDate:   "",
				Units:     "",
				TempScale: "",
			},
			wantErr: "connection refused",
		},
		{
			name:   "sucesso na bahia",
			client: &fakeHTTPClient{
				resp: &http.Response{
					StatusCode: 200,
					Body: io.NopCloser(strings.NewReader("{invalid json")),
				},
			},
			config: viperWithURI("http://fake-api"),
			request: &ForecastRequest{
				Latitude:  "",
				Longitude: "",
				StartDate: "",
				EndDate:   "",
				Units:     "",
				TempScale: "",
			},
			wantErr: "invalid character 'i' looking for beginning of object key string",
		},
		{
			name:   "diferente de 200",
			client: &fakeHTTPClient{
				resp: &http.Response{
					StatusCode: 500,
					Body: io.NopCloser(strings.NewReader("{invalid json")),
				},
			},
			config: viperWithURI("http://fake-api"),
			request: &ForecastRequest{
				Latitude:  "",
				Longitude: "",
				StartDate: "",
				EndDate:   "",
				Units:     "",
				TempScale: "",
			},
			wantErr: "Resposta da api fora do esperado",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fc := ForecastClient{
				Client: tc.client,
				Conf:   tc.config,
			}

			_, err := fc.GetForecast(context.Background(), tc.request)
			if tc.wantErr == "" && err != nil {
				t.Fatalf("não esperava erro, mas veio: %v", err)
			}
			if tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)) {
				t.Fatalf("esperava erro contendo %q, mas veio: %v", tc.wantErr, err)
			}
		})
	}

}

func viperWithURI(uri string) *viper.Viper {
	v := viper.New()
	v.Set("coordinates.uri", uri)
	return v
}
