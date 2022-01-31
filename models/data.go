package model

type Data struct {
	Group string      `json:"group"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
