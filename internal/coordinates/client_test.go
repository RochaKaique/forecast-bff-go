package coordinates

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

// fakeHTTPClient implementa a interface Do usada no client
type fakeHTTPClient struct {
	resp *http.Response
	err  error
}

func (f *fakeHTTPClient) Do(*http.Request) (*http.Response, error) {
	return f.resp, f.err
}

func TestCoordinatesClient_GetCoordinates(t *testing.T) {
	tests := []struct {
		name    string
		client  *fakeHTTPClient
		config  *viper.Viper
		wantErr string
	}{
		{
			name:    "erro no transporte (Client.Do falha)",
			client:  &fakeHTTPClient{err: errors.New("connection refused")},
			config:  viperWithURI("http://fake-api"),
			wantErr: "connection refused",
		},
		{
			name: "JSON inválido na resposta",
			client: &fakeHTTPClient{
				resp: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader("{invalid json")),
				},
			},
			config:  viperWithURI("http://fake-api"),
			wantErr: "invalid character",
		},
		{
			name: "sucesso com JSON válido",
			client: &fakeHTTPClient{
				resp: &http.Response{
					StatusCode: 200,
					Body: io.NopCloser(strings.NewReader(`[
						{"lat": "-23.55", "lon": "-46.63"}
					]`)),
				},
			},
			config:  viperWithURI("http://fake-api"),
			wantErr: "",
		},
		{
			name: "slice vazio na resposta",
			client: &fakeHTTPClient{
				resp: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader("nil")),
				},
			},
			config:  viperWithURI("http://fake-api"),
			wantErr: "invalid character 'i' in literal null (expecting 'u')", // vai estourar no coordinates[0]
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := CoordinatesClient{
				Client: tt.client,
				Conf:   tt.config,
			}

			_, err := cc.GetCoordinates(context.Background(), "12345")
			if tt.wantErr == "" && err != nil {
				t.Fatalf("não esperava erro, mas veio: %v", err)
			}
			if tt.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tt.wantErr)) {
				t.Fatalf("esperava erro contendo %q, mas veio: %v", tt.wantErr, err)
			}
		})
	}
}


func viperWithURI(uri string) *viper.Viper {
	v := viper.New()
	v.Set("coordinates.uri", uri)
	return v
}
