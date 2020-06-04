package main

import "strconv"

type Response struct {
	Data Data `json:"data"`
}

type Data struct {
	Json Check `json:"json"`
}

type Check struct {
	OrganizationAddress string    `json:"retailPlaceAddress"`
	OrganizationName    string    `json:"user"`
	Products            []Product `json:"items"`
	Sum                 float64   `json:"totalSum"`
}

type Product struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Sum      float64 `json:"sum"`
}

func (c *Check) String() string {
	return "Адрес организации: " + c.OrganizationAddress +
		"\nИмя организации: " + c.OrganizationName +
		"\nОбщая сумма чека: " + strconv.FormatFloat(c.Sum, 'f', -1, 64)
}
