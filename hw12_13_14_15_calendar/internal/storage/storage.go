package storage

type Storage interface {
	AddEvent(event Event) (string, error)
	ChangeEvent(eventID string, event Event) error
	DeleteEvent(eventID string) error
	ListEventsForDay(date string) ([]Event, error)
	ListEventsForWeek(date string) ([]Event, error)
	ListEventsForMonth(date string) ([]Event, error)
}
