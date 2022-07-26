package utils

import "net/mail"

func CheckEmail(adress string) error {
	_, err := mail.ParseAddress(adress)
	return err
}
