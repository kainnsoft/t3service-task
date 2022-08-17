package utils

import (
	"fmt"
	"net/mail"
)

func CheckEmail(address string) error {
	_, err := mail.ParseAddress(address)
	if err != nil {
		return fmt.Errorf("wrong email address: %w", err)
	}

	return nil
}
