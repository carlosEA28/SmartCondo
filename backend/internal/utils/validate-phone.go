package utils

import (
	"fmt"

	"github.com/nyaruka/phonenumbers/v2"
)

func ValidatePhoneNumber(phone string) (string, error) {

	num, err := phonenumbers.Parse(phone, "BR")

	if err != nil {
		return "", fmt.Errorf("Error parsing phone number")
	}

	valid := phonenumbers.IsValidNumber(num)

	if valid != true {
		return "", fmt.Errorf("The provided number is not valid")
	}

	formated := phonenumbers.Format(num, phonenumbers.NATIONAL)

	return formated, nil
}
