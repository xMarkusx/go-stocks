package shared

import "time"

type CommandDate string

func (date CommandDate) Get() string {
	if date != "" {
		return string(date)
	}

	today := time.Now().Format("2006-01-02")
	return today
}
