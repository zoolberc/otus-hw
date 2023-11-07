package memorystorage

import (
	"sync"
	"time"

	"github.com/zoolberc/otus-hw/hw12_13_14_15_calendar/internal/constants"
	"github.com/zoolberc/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/zoolberc/otus-hw/hw12_13_14_15_calendar/utils"
)

type SearchPeriod struct {
	years  int
	months int
	days   int
}

type EventStorage struct {
	events map[string]storage.Event
	mu     sync.RWMutex
}

func New() *EventStorage {
	return &EventStorage{
		events: make(map[string]storage.Event),
	}
}

func (e *EventStorage) AddEvent(event storage.Event) (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	eventDate, err := utils.ParseDate(event.Date)
	if err != nil {
		return "", err
	}
	if eventDate.Before(time.Now()) {
		return "", constants.ErrDateBefore
	}
	for _, ev := range e.events {
		if ev.Date == event.Date {
			return "", constants.ErrDateBusy
		}
	}
	event.ID = utils.GenerateUUID()
	e.events[event.ID] = event

	return event.ID, nil
}

func (e *EventStorage) ChangeEvent(eventID string, event storage.Event) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, ok := e.events[eventID]
	if !ok {
		return constants.ErrEventNotFound
	}

	eventDate, err := utils.ParseDate(event.Date)
	if err != nil {
		return err
	}

	if eventDate.Before(time.Now()) {
		return constants.ErrDateBefore
	}
	for _, ev := range e.events {
		if ev.Date == event.Date {
			return constants.ErrDateBusy
		}
	}

	event.ID = eventID
	e.events[eventID] = event

	return nil
}

func (e *EventStorage) DeleteEvent(eventID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, ok := e.events[eventID]
	if !ok {
		return constants.ErrEventNotFound
	}
	delete(e.events, eventID)

	return nil
}

func (e *EventStorage) ListEventsForDay(date string) ([]storage.Event, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	searchDate, err := utils.ParseDate(date)
	if err != nil {
		return nil, err
	}
	return searchOverCurrentDay(searchDate, e.events)
}

func (e *EventStorage) ListEventsForWeek(date string) ([]storage.Event, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	startDate, err := utils.ParseDate(date)
	if err != nil {
		return nil, err
	}
	return searchEventsOverPeriod(startDate, SearchPeriod{days: 7}, e.events)
}

func (e *EventStorage) ListEventsForMonth(date string) ([]storage.Event, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	startDate, err := utils.ParseDate(date)
	if err != nil {
		return nil, err
	}
	return searchEventsOverPeriod(startDate, SearchPeriod{months: 1}, e.events)
}

func searchEventsOverPeriod(
	startDate time.Time,
	period SearchPeriod,
	events map[string]storage.Event,
) ([]storage.Event, error) {
	endDate := startDate.AddDate(period.years, period.months, period.days)
	eventList := make([]storage.Event, 0)
	for _, event := range events {
		eventDate, err := utils.ParseDate(event.Date)
		if err != nil {
			return nil, err
		}
		if eventDate.Before(startDate) {
			continue
		}
		if eventDate.Before(endDate) {
			eventList = append(eventList, event)
		}
	}
	return eventList, nil
}

func searchOverCurrentDay(date time.Time, events map[string]storage.Event) ([]storage.Event, error) {
	eventList := make([]storage.Event, 0)
	date = date.Truncate(24 * time.Hour)
	for _, event := range events {
		eventDate, err := utils.ParseDate(event.Date)
		if err != nil {
			return nil, err
		}
		eventDate = eventDate.Truncate(24 * time.Hour)
		timeDiff := int(date.Sub(eventDate).Hours())
		if timeDiff == 0 {
			eventList = append(eventList, event)
		}
	}
	return eventList, nil
}
