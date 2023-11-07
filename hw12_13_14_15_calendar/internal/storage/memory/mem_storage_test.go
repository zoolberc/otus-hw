package memorystorage

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zoolberc/otus-hw/hw12_13_14_15_calendar/internal/constants"
	"github.com/zoolberc/otus-hw/hw12_13_14_15_calendar/internal/storage"
)

func TestMemStorage(t *testing.T) {
	event := storage.Event{
		Title:       "Визит к стаматологу",
		Date:        "2024-11-02T13:00:00Z",
		Duration:    0,
		Description: "Приём у стоматолога на 15:00",
		UserID:      "1",
	}

	t.Run("Success add event to storage", func(t *testing.T) {
		st := New()
		eventID, err := st.AddEvent(event)
		require.NoError(t, err)
		event.ID = eventID
		require.Equal(t, event, st.events[eventID])
	})

	t.Run("Success change event", func(t *testing.T) {
		st := New()
		eventID, err := st.AddEvent(event)
		require.NoError(t, err)
		event.ID = eventID
		event.Date = "2024-11-03T13:00:00Z"
		err = st.ChangeEvent(eventID, event)
		require.NoError(t, err)
		require.Equal(t, event, st.events[eventID])
	})

	t.Run("Success delete event", func(t *testing.T) {
		st := New()
		eventID, err := st.AddEvent(event)
		require.NoError(t, err)
		err = st.DeleteEvent(eventID)
		require.NoError(t, err)
		_, ok := st.events[eventID]
		require.False(t, ok)
	})

	t.Run("Get list events for one day", func(t *testing.T) {
		st := New()
		for i := 0; i < 5; i++ {
			event.Date = fmt.Sprintf("2024-11-03T13:0%s:00Z", strconv.Itoa(i))
			evID, _ := st.AddEvent(event)
			_ = evID
		}
		event.Date = "2024-11-04T13:00:00Z"
		_, err := st.AddEvent(event)
		require.NoError(t, err)
		events, err := st.ListEventsForDay("2024-11-03T00:00:00Z")
		require.NoError(t, err)
		require.Equal(t, len(events), 5)
	})

	t.Run("Get list events for week", func(t *testing.T) {
		st := New()
		for i := 0; i < 5; i++ {
			event.Date = fmt.Sprintf("2024-11-1%sT13:00:00Z", strconv.Itoa(i))
			evID, _ := st.AddEvent(event)
			_ = evID
		}
		event.Date = "2024-12-04T13:00:00Z"
		_, err := st.AddEvent(event)
		require.NoError(t, err)
		events, err := st.ListEventsForWeek("2024-11-06T13:00:00Z")
		require.NoError(t, err)
		require.Equal(t, len(events), 3)
	})

	t.Run("Get list events for month", func(t *testing.T) {
		st := New()
		for i := 0; i < 5; i++ {
			event.Date = fmt.Sprintf("2024-11-1%sT13:00:00Z", strconv.Itoa(i))
			evID, _ := st.AddEvent(event)
			_ = evID
		}
		event.Date = "2024-12-04T13:00:00Z"
		_, err := st.AddEvent(event)
		require.NoError(t, err)
		events, err := st.ListEventsForMonth("2024-11-01T13:00:00Z")
		require.NoError(t, err)
		require.Equal(t, len(events), 5)
	})

	t.Run("Add event. Get error busy date", func(t *testing.T) {
		st := New()
		_, err := st.AddEvent(event)
		require.NoError(t, err)
		_, err = st.AddEvent(event)
		require.ErrorIs(t, err, constants.ErrDateBusy)
	})

	t.Run("Add event. Get error early date", func(t *testing.T) {
		st := New()
		event.Date = "2020-12-04T13:00:00Z"
		_, err := st.AddEvent(event)
		require.ErrorIs(t, err, constants.ErrDateBefore)
	})

	t.Run("Change event. Get error busy date", func(t *testing.T) {
		st := New()
		event := storage.Event{
			Title:       "Визит к стаматологу",
			Date:        "2024-11-02T13:00:00Z",
			Duration:    0,
			Description: "Приём у стоматолога на 15:00",
			UserID:      "1",
		}
		eventID, err := st.AddEvent(event)
		require.NoError(t, err)
		event.ID = eventID
		err = st.ChangeEvent(eventID, event)
		require.ErrorIs(t, err, constants.ErrDateBusy)
	})

	t.Run("Change event. Get error early date", func(t *testing.T) {
		st := New()
		event := storage.Event{
			Title:       "Визит к стаматологу",
			Date:        "2024-11-02T13:00:00Z",
			Duration:    0,
			Description: "Приём у стоматолога на 15:00",
			UserID:      "1",
		}
		eventID, err := st.AddEvent(event)
		require.NoError(t, err)
		event.ID = eventID
		event.Date = "2020-12-04T13:00:00Z"
		err = st.ChangeEvent(eventID, event)
		require.ErrorIs(t, err, constants.ErrDateBefore)
	})

	t.Run("Change event. Get error event not found", func(t *testing.T) {
		st := New()
		event.ID = "123"
		err := st.ChangeEvent("123", event)
		require.ErrorIs(t, err, constants.ErrEventNotFound)
	})

	t.Run("Delete event. Get error event not found", func(t *testing.T) {
		st := New()
		event.ID = "123"
		err := st.DeleteEvent("123")
		require.ErrorIs(t, err, constants.ErrEventNotFound)
	})
}
