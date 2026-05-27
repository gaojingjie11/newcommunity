package repository

import "time"

func nowPtr() *time.Time {
	now := time.Now()
	return &now
}
