package repository

import "errors"

var (
	ErrParkingSpaceUnavailable = errors.New("parking space unavailable")
	ErrDuplicateParkingNo      = errors.New("duplicate parking number")
	ErrParkingUserNotFound     = errors.New("parking user not found")
	ErrPropertyFeePaid         = errors.New("property fee already paid")
)
