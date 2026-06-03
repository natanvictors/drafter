package parser

import "encoding/json"

type ChampionInfo struct {
	Attack     string `json:"attack"`
	Defense    string `json:"defense"`
	Magic      string `json:"magic"`
	Difficulty string `json:"difficulty"`
}

type ChampionImage struct {
	Full   string `json:"full"`
	Sprite string `json:"sprite"`
	Group  string `json:"group"`
	X      string `json:"x"`
	Y      string `json:"y"`
	W      string `json:"w"`
	H      string `json:"h"`
}

type ChampionStats struct {
	Hp                   int     `json:"hp"`
	Hpperlevel           int     `json:"hpperlevel"`
	Mp                   int     `json:"mp"`
	Mpperlevel           int     `json:"mpperlevel"`
	Movespeed            int     `json:"movespeed"`
	Armor                int     `json:"armor"`
	Armorperlevel        float64 `json:"armorperlevel"`
	Spellblock           int     `json:"spellblock"`
	Spellblockperlevel   float64 `json:"spellblockperlevel"`
	Attackrange          int     `json:"attackrange"`
	Hpregen              int     `json:"hpregen"`
	Hpregenperlevel      float64 `json:"hpregenperlevel"`
	Mpregen              int     `json:"mpregen"`
	Mpregenperlevel      int     `json:"mpregenperlevel"`
	Crit                 int     `json:"crit"`
	Critperlevel         int     `json:"critperlevel"`
	Attackdamage         int     `json:"attackdamage"`
	Attackdamageperlevel int     `json:"attackdamageperlevel"`
	Attackspeedperlevel  float64 `json:"attackspeedperlevel"`
	Attackspeed          float64 `json:"attackspeed"`
}

type Champion struct {
	Version string        `json:"version"`
	Id      string        `json:"id"`
	Key     string        `json:"key"`
	Name    string        `json:"name"`
	Title   string        `json:"title"`
	Blurb   string        `json:"blurb"`
	Info    ChampionInfo  `json:"info"`
	Image   ChampionImage `json:"image"`
	Tags    []string      `json:"tags"`
	Partype string        `json:"partype"`
	Stats   ChampionStats `json:"stats"`
}

type ChampionResponse struct {
	Type    string              `json:"type"`
	Format  string              `json:"format"`
	Version string              `json:"version"`
	Data    map[string]Champion `json:"data"`
}

func ParseChampion(data []byte) (ChampionResponse, error) {
	var resp ChampionResponse

	if err := json.Unmarshal(data, &resp); err != nil {
		return ChampionResponse{}, err
	}

	return resp, nil
}
