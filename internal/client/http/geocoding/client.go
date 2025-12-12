package geocoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type client struct {
	httpClient http.Client
}

type Response struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
}

func NewClinet(httpClient *http.Client) *client {
	return &client{
		httpClient: *httpClient,
	}
}

func (c *client) GetCoords(city string) (Response, error) {
	res, err := c.httpClient.Get(
		fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&format=json&language=ru", city),
	)

	if err != nil {
		return Response{}, err
	}

	if res.StatusCode != http.StatusOK {
		return Response{}, errors.New(res.Status)
	}

	var geoResult struct {
		Result []Response `json:"results"`
	}

	err = json.NewDecoder(res.Body).Decode(&geoResult)
	if err != nil {
		return Response{}, err
	}

	if len(geoResult.Result) == 0 {
		return Response{}, fmt.Errorf("город '%s' не найден", city)
	}

	return geoResult.Result[0], nil
}
