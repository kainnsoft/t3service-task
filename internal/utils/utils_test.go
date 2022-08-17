package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckEmail(t *testing.T) {
	req := require.New(t)

	var err error

	t.Run("correct email test", func(t *testing.T) {
		err = CheckEmail("addr@gmail.com")
		req.NoError(err)
	})

	tests := map[string]struct {
		address string
		want    string
	}{
		"email without@":          {"addrgmail.com", "wrong email address"},
		"email with emptybefore@": {"@gmailcom", "wrong email address"},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			err = CheckEmail(testCase.address)
			req.ErrorContains(err, testCase.want)
		})
	}
}
