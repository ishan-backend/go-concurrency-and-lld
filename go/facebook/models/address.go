package models

// Address represents the address of a user
type Address struct {
	ZipCode    int
	StreetName string
	City       string
	Country    string
}

// NewAddress creates a new Address instance
func NewAddress(zipCode int, streetName, city, country string) *Address {
	return &Address{
		ZipCode:    zipCode,
		StreetName: streetName,
		City:       city,
		Country:    country,
	}
}
