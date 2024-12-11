package models

import "fmt"

type Product struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Category string `json:"category"`
	Owner    uint   `json:"owner"`
}

func (p Product) String() string {
	return fmt.Sprintf("Product{ID: %d, Name: %s, Price: %s, Category: %s, Owner: %v}",
		p.Id, p.Name, p.Price, p.Category, p.Owner)
}
