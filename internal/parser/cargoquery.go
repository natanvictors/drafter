package parser

import (
	"encoding/json"
)

type Cargoresult struct {
	Cargoresponse []*Cargoresponse `json:"cargoquery,omitempty"`
	Cargoerror    *Cargoerror      `json:"error,omitempty"`
}

type Cargoresponse struct {
	Game Game `json:"title"`
}

type Cargoerror struct {
	Code string `json:"code"`
	Info string `json:"info"`
	Star string `json:"*"`
}

type Game struct {
	Championship string `json:"Page"`
	Patch        string `json:"Patch"`
	Winner       string `json:"Winner"`
	GameID       string `json:"GameId"`
	MatchID      string `json:"MatchId"`
	Team1Role1   string `json:"Team1Role1"`
	Team1Role2   string `json:"Team1Role2"`
	Team1Role3   string `json:"Team1Role3"`
	Team1Role4   string `json:"Team1Role4"`
	Team1Role5   string `json:"Team1Role5"`
	Team2Role1   string `json:"Team2Role1"`
	Team2Role2   string `json:"Team2Role2"`
	Team2Role3   string `json:"Team2Role3"`
	Team2Role4   string `json:"Team2Role4"`
	Team2Role5   string `json:"Team2Role5"`
	Team1Pick1   string `json:"Team1Pick1"`
	Team1Pick2   string `json:"Team1Pick2"`
	Team1Pick3   string `json:"Team1Pick3"`
	Team1Pick4   string `json:"Team1Pick4"`
	Team1Pick5   string `json:"Team1Pick5"`
	Team2Pick1   string `json:"Team2Pick1"`
	Team2Pick2   string `json:"Team2Pick2"`
	Team2Pick3   string `json:"Team2Pick3"`
	Team2Pick4   string `json:"Team2Pick4"`
	Team2Pick5   string `json:"Team2Pick5"`
	Team1Ban1    string `json:"Team1Ban1"`
	Team1Ban2    string `json:"Team1Ban2"`
	Team1Ban3    string `json:"Team1Ban3"`
	Team1Ban4    string `json:"Team1Ban4"`
	Team1Ban5    string `json:"Team1Ban5"`
	Team2Ban1    string `json:"Team2Ban1"`
	Team2Ban2    string `json:"Team2Ban2"`
	Team2Ban3    string `json:"Team2Ban3"`
	Team2Ban4    string `json:"Team2Ban4"`
	Team2Ban5    string `json:"Team2Ban5"`
}

func (g Game) Teams() (team1, team2 Team) {
	team1 = Team{
		Side:  1,
		Roles: []string{g.Team1Role1, g.Team1Role2, g.Team1Role3, g.Team1Role4, g.Team1Role5},
		Picks: []string{g.Team1Pick1, g.Team1Pick2, g.Team1Pick3, g.Team1Pick4, g.Team1Pick5},
		Bans:  []string{g.Team1Ban1, g.Team1Ban2, g.Team1Ban3, g.Team1Ban4, g.Team1Ban5},
	}
	team2 = Team{
		Side:  2,
		Roles: []string{g.Team2Role1, g.Team2Role2, g.Team2Role3, g.Team2Role4, g.Team2Role5},
		Picks: []string{g.Team2Pick1, g.Team2Pick2, g.Team2Pick3, g.Team2Pick4, g.Team2Pick5},
		Bans:  []string{g.Team2Ban1, g.Team2Ban2, g.Team2Ban3, g.Team2Ban4, g.Team2Ban5},
	}

	return team1, team2
}

type Team struct {
	Name  string
	Side  int
	Roles []string
	Picks []string
	Bans  []string
}

func ParseCargo(data []byte) (*Cargoresult, error) {
	var cargoquery Cargoresult
	//log.Printf("%s\n", data)
	err := json.Unmarshal(data, &cargoquery)
	if err != nil {
		return nil, err
	}
	return &cargoquery, nil
}
