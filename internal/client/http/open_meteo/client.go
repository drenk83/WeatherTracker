package openmeteo

import "net/http"

type client struct {
	httpClient http.Client
}

func NewClinet(httpClient *http.Client) *client {
	return &client{
		httpClient: *httpClient,
	}
}

func (c *client) GetTemperature(lat, long float64) (float64, error) {

}
