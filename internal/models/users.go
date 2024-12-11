package models

import (
	"fmt"
	"strings"
)

type User struct {
	Id       uint      `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Products []Product `json:"products"`
}

func (u User) String() string {
	var products string
	if len(u.Products) > 0 {
		productStrings := make([]string, len(u.Products))
		for i, p := range u.Products {
			productStrings[i] = p.String()
		}
		products = strings.Join(productStrings, ", ")
	}

	return fmt.Sprintf("User{ID: %d, Name: %s, Email: %s, Products: [%s]}",
		u.Id, u.Name, u.Email, products)
}
