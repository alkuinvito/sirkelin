package utils

import (
	b64 "encoding/base64"
	"errors"
	"os"
)

func VerifyClientID(id string) error {
	decoded, err := b64.StdEncoding.DecodeString(id)
	if err != nil {
		return err
	}

	if string(decoded) != os.Getenv("CLIENT_ID_NEXT") {
		return errors.New("invalid client id")
	}
	return nil
}
