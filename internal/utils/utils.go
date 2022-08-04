package utils

import "net/mail"

func CheckEmail(address string) error {
	_, err := mail.ParseAddress(address)
	return err
}
