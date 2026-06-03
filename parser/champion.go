package parser

import "encoding/json"

type Champion struct {
	Name string `json:"name"`
}

type ChampionResponse struct {
	Data map[string]json.RawMessage `json:"data"`
}

func ParseChampion(data []byte) (ChampionResponse, error) {
	var resp ChampionResponse

	if err := json.Unmarshal(data, &resp); err != nil {
		return ChampionResponse{}, err
	}

	return resp, nil
}
