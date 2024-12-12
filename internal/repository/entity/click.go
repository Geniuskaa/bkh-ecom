package entity

import "time"

type EventsFilter struct {
	EventType  string
	EventsFrom time.Time
	EventsTo   time.Time
}
