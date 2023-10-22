package sqlstorage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/zoolberc/otus-hw/hw12_13_14_15_calendar/internal/constants"
	"github.com/zoolberc/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/zoolberc/otus-hw/hw12_13_14_15_calendar/utils"
)

type SearchPeriod struct {
	years  int
	months int
	days   int
}

type DataBaseConf struct {
	DBHost     string `yaml:"host"`
	DBPort     string `yaml:"port"`
	DBName     string `yaml:"name"`
	DBUser     string `yaml:"user"`
	DBPassword string `yaml:"password"`
}

type EventStorage struct {
	db     *sqlx.DB
	ctx    context.Context
	dbConf DataBaseConf
}

func New(dbConnectConfigs DataBaseConf) *EventStorage {
	return &EventStorage{
		dbConf: dbConnectConfigs,
	}
}

func (s *EventStorage) Connect(ctx context.Context) error {
	db, err := sqlx.Open("pgx", getPsqlInfo(s.dbConf))
	if err != nil {
		return constants.ErrConnectionToDB
	}
	s.db = db
	s.ctx = ctx
	return nil
}

func (s *EventStorage) Close(ctx context.Context) error { //nolint:revive
	err := s.db.Close()
	if err != nil {
		return constants.ErrDisconnectionToDB
	}
	return nil
}

func (s *EventStorage) AddEvent(event storage.Event) (string, error) {
	eventDate, err := utils.ParseDate(event.Date)
	if err != nil {
		return "", err
	}
	if err = checkEventDate(s.ctx, eventDate, s.db); err != nil {
		return "", err
	}

	event.ID = utils.GenerateUUID()
	_, err = s.db.ExecContext(s.ctx, constants.AddEventQuery,
		event.ID, event.Title, event.Date, event.Duration, event.Duration, event.UserID, event.Reminder)
	if err != nil {
		return "", err
	}
	return event.ID, nil
}

func (s *EventStorage) ChangeEvent(eventID string, event storage.Event) error {
	eventDate, err := utils.ParseDate(event.Date)
	if err != nil {
		return err
	}
	if eventDate.Before(time.Now()) {
		return constants.ErrDateBefore
	}
	_, err = s.db.ExecContext(s.ctx, constants.ChangeEventQuery,
		event.Title, event.Date, event.Duration, event.Duration, event.UserID, event.Reminder, eventID)
	if err != nil {
		return err
	}
	return nil
}

func (s *EventStorage) DeleteEvent(eventID string) error {
	_, err := s.db.ExecContext(s.ctx, constants.DeleteEventQuery, eventID)
	if err != nil {
		return err
	}
	return nil
}

func (s *EventStorage) ListEventsForDay(date string) ([]storage.Event, error) {
	searchDate, err := utils.ParseDate(date)
	if err != nil {
		return nil, err
	}
	return searchEventsOverPeriod(s.ctx, searchDate, SearchPeriod{days: 1}, s.db)
}

func (s *EventStorage) ListEventsForWeek(date string) ([]storage.Event, error) {
	searchDate, err := utils.ParseDate(date)
	if err != nil {
		return nil, err
	}
	return searchEventsOverPeriod(s.ctx, searchDate, SearchPeriod{days: 7}, s.db)
}

func (s *EventStorage) ListEventsForMonth(date string) ([]storage.Event, error) {
	searchDate, err := utils.ParseDate(date)
	if err != nil {
		return nil, err
	}
	return searchEventsOverPeriod(s.ctx, searchDate, SearchPeriod{months: 7}, s.db)
}

func searchEventsOverPeriod(
	ctx context.Context,
	startDate time.Time,
	period SearchPeriod,
	db *sqlx.DB,
) ([]storage.Event, error) {
	eventList := make([]storage.Event, 0)
	rows, err := db.QueryContext(ctx, constants.GetEventsOverPeriodQuery, startDate.Format("2006-01-02"),
		startDate.AddDate(period.years, period.months, period.days).Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ev storage.Event
		if err := rows.Scan(&ev.ID, &ev.Title, &ev.Date, &ev.Duration, &ev.Duration, &ev.UserID, &ev.Reminder); err != nil {
			return nil, err
		}
		eventList = append(eventList, ev)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return eventList, nil
}

func checkEventDate(ctx context.Context, date time.Time, db *sqlx.DB) error {
	if date.Before(time.Now()) {
		return constants.ErrDateBefore
	}
	rows, err := db.QueryContext(ctx, constants.CheckCountEventOnDateQuery, date)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return constants.ErrDateBusy
	}
	return nil
}

func getPsqlInfo(conf DataBaseConf) string {
	host, isSet := os.LookupEnv(constants.DatabaseHostEnv)
	if !isSet {
		host = conf.DBHost
	}
	port, isSet := os.LookupEnv(constants.DatabasePortEnv)
	if !isSet {
		port = conf.DBPort
	}
	dbName, isSet := os.LookupEnv(constants.DatabaseNameEnv)
	if !isSet {
		dbName = conf.DBName
	}
	user, isSet := os.LookupEnv(constants.DatabaseUserEnv)
	if !isSet {
		user = conf.DBUser
	}
	password, isSet := os.LookupEnv(constants.DatabasePasswordEnv)
	if !isSet {
		password = conf.DBPassword
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
}
