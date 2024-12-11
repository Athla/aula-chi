package models

type Product struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Category string `json:"category"`
	Owner    string `json:"owner"`
}
