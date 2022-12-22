package controllers

import "time"

type SystemClock struct{}

func (_ SystemClock) Now() time.Time { return time.Now() }

type Clock interface {
	Now() time.Time
}
