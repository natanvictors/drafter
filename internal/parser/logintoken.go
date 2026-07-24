package parser

import (
	"encoding/json"
)

type LoginTokenResp struct {
	Query struct {
		Tokens struct {
			LoginToken string `json:"logintoken"`
		} `json:"tokens"`
	} `json:"query"`
}

func ParseToken(data []byte) (LoginTokenResp, error) {
	var token LoginTokenResp

	err := json.Unmarshal(data, &token)

	return token, err
}
