package main

import (
	"time"
)

type CheckInfoFromBot struct {
	DateTimeString string `schema:"t"`
	Sum            string `schema:"s"`
	FN             string `schema:"fn"`
	FD             string `schema:"i"`
	FP             string `schema:"fp"`
	N              string `schema:"n"`
}

func (c *CheckInfoFromBot) GetTime() (time.Time, error) {
	dateTime, err := time.Parse("20060102T1504", c.DateTimeString)
	if err != nil {
		return time.Time{}, err
	}

	return dateTime, nil
}
