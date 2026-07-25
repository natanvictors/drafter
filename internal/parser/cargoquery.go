package parser

import (
	"encoding/json"
	"fmt"
)

type Cargoresult struct {
	Cargoresponse []*Cargoresponse `json:"cargoquery,omitempty"`
	Cargoerror    *Cargoerror      `json:"error,omitempty"`
}

type Cargoresponse struct {
	Title Title `json:"title"`
}

type Cargoerror struct {
	Code string `json:"code"`
	Info string `json:"info"`
	Star string `json:"*"`
}

type Title struct {
	Page                  string `json:"Page"`
	Patch                 string `json:"Patch"`
	Team1Role1            string `json:"Team1Role1"`
	Team1Role2            string `json:"Team1Role2"`
	Team1Role3            string `json:"Team1Role3"`
	Team1Role4            string `json:"Team1Role4"`
	Team1Role5            string `json:"Team1Role5"`
	Team2Role1            string `json:"Team2Role1"`
	Team2Role2            string `json:"Team2Role2"`
	Team2Role3            string `json:"Team2Role3"`
	Team2Role4            string `json:"Team2Role4"`
	Team2Role5            string `json:"Team2Role5"`
	Team1Ban1             string `json:"Team1Ban1"`
	Team1Ban2             string `json:"Team1Ban2"`
	Team1Ban3             string `json:"Team1Ban3"`
	Team1Ban4             string `json:"Team1Ban4"`
	Team1Ban5             string `json:"Team1Ban5"`
	Team1Pick1            string `json:"Team1Pick1"`
	Team1Pick2            string `json:"Team1Pick2"`
	Team1Pick3            string `json:"Team1Pick3"`
	Team1Pick4            string `json:"Team1Pick4"`
	Team1Pick5            string `json:"Team1Pick5"`
	Team2Ban1             string `json:"Team2Ban1"`
	Team2Ban2             string `json:"Team2Ban2"`
	Team2Ban3             string `json:"Team2Ban3"`
	Team2Ban4             string `json:"Team2Ban4"`
	Team2Ban5             string `json:"Team2Ban5"`
	Team2Pick1            string `json:"Team2Pick1"`
	Team2Pick2            string `json:"Team2Pick2"`
	Team2Pick3            string `json:"Team2Pick3"`
	Team2Pick4            string `json:"Team2Pick4"`
	Team2Pick5            string `json:"Team2Pick5"`
	Team1                 string `json:"Team1"`
	Team2                 string `json:"Team2"`
	Winner                string `json:"Winner"`
	Team1PicksByRoleOrder string `json:"Team1PicksByRoleOrder"`
	Team2PicksByRoleOrder string `json:"Team2PicksByRoleOrder"`
	GameID                string `json:"GameId"`
	MatchID               string `json:"MatchId"`
	DateTimeUTC           string `json:"DateTime UTC"`
	DateTimeUTCPrecision  string `json:"DateTime UTC__precision"`
}

func ParseCargo(data []byte) (*Cargoresult, error) {
	var cargoquery Cargoresult
	err := json.Unmarshal(data, &cargoquery)
	if err != nil {
		return nil, err
	}
	if cargoquery.Cargoerror != nil {
		return nil, fmt.Errorf("cargoquery error: %x", cargoquery.Cargoerror.Code)
	}
	return &cargoquery, nil
}
