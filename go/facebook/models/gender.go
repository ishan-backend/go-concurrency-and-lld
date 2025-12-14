package models

// Gender represents the gender of a user
type Gender string

const (
	GenderFemale      Gender = "FEMALE"
	GenderMale        Gender = "MALE"
	GenderTransgender Gender = "TRANSGENDER"
)
