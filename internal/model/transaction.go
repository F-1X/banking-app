package model

import "encoding/json"

type Transaction struct {
	From   json.Number `json:"from"`
	To     json.Number `json:"to"`
	Amount json.Number `json:"amount"`
}
