package constants

import "errors"

var (
	ErrEventNotFound = errors.New("event not found")
	ErrDateBusy      = errors.New("this time is busy")
	ErrDateBefore    = errors.New("event date is earlier than current date")
)

var (
	ErrIncorrectDateFormat = errors.New("incorrect date format")
	ErrConnectionToDB      = errors.New("error connecting to database")
	ErrDisconnectionToDB   = errors.New("error disconnecting from database")
)
