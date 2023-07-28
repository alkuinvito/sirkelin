package utils

import "errors"

var (
	ErrMinimumParticipant        = errors.New("room must have two or more participants")
	ErrPrivateParticipantsNumber = errors.New("private room must only have two participants")
)
