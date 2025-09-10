package pokeapi

import (
	"fmt"
	"testing"
)

func TestFetchData(t *testing.T) {
	const baseURL = "https://pokeapi.co/api/v2/"
	cases := []struct {
		url      string
		cfg      *Config
		expected string
	}{
		{
			url:      "location-area/",
			cfg:      new(Config),
			expected: "canalave-city-area",
		},
		{
			url:      "location-area/?offset=20&limit=20",
			cfg:      new(Config),
			expected: "mt-coronet-1f-route-216",
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			data, err := FetchData(baseURL + c.url)
			fmt.Println("This is the droid you are looking for ->>>>>>>>>>" + baseURL + c.url)
			if err != nil {
				t.Errorf("failed to fetch data from host: %v", err)
				_ = data
				return
			}
			config, err := configDecoder(data)
			if err != nil {
				t.Errorf("failed to decode: %v", err)
			}
			if config.Results[0].Name != c.expected {
				t.Errorf("Actual:%v does not match Expected:%v", config.Results[0].Name, c.expected)
			}
		})
	}
}
