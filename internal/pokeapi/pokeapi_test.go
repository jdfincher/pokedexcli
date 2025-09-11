package pokeapi

import (
	"fmt"
	"testing"
)

func TestFetchData(t *testing.T) {
	const baseURL = "https://pokeapi.co/api/v2/"
	cases := []struct {
		url      string
		loc      *Locations
		expected string
	}{
		{
			url:      "location-area/",
			loc:      new(Locations),
			expected: "canalave-city-area",
		},
		{
			url:      "location-area/?offset=20&limit=20",
			loc:      new(Locations),
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
			loc, err := locationsDecoder(data)
			if err != nil {
				t.Errorf("failed to decode: %v", err)
			}
			if loc.Results[0].Name != c.expected {
				t.Errorf("Actual:%v does not match Expected:%v", loc.Results[0].Name, c.expected)
			}
		})
	}
}
