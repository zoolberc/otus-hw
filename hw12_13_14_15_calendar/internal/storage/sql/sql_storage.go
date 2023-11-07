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
	dbConf DataBaseConf
}

func New(dbConnectConfigs DataBaseConf) *EventStorage {
	return &EventStorage{
		dbConf: dbConnectConfigs,
	}
}

func (s *EventStorage) Connect(ctx context.Context) error { //nolint:revive
	db, err := sqlx.Open("pgx", getPsqlInfo(s.dbConf))
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	s.db = db
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
	ctx := context.Background()
	eventDate, err := utils.ParseDate(event.Date)
	if err != nil {
		return "", err
	}
	if err = checkEventDate(ctx, eventDate, s.db); err != nil {
		return "", err
	}

	event.ID = utils.GenerateUUID()
	addEventQuery := `insert into events(id, title, date, duration, description, user_id, reminder) 
						values($1, $2, $3, $4, $5, $6, $7)`
	_, err = s.db.ExecContext(ctx, addEventQuery,
		event.ID, event.Title, event.Date, event.Duration, event.Duration, event.UserID, event.Reminder)
	if err != nil {
		return "", err
	}
	return event.ID, nil
}

func (s *EventStorage) ChangeEvent(eventID string, event storage.Event) error {
	ctx := context.Background()
	eventDate, err := utils.ParseDate(event.Date)
	if err != nil {
		return err
	}
	if eventDate.Before(time.Now()) {
		return constants.ErrDateBefore
	}
	changeEventQuery := `update events set title = $1, date = $2, duration = $3,
                  		description = $4, user_id = $5, reminder = $6 WHERE Manufacturer = $7`
	_, err = s.db.ExecContext(ctx, changeEventQuery,
		event.Title, event.Date, event.Duration, event.Duration, event.UserID, event.Reminder, eventID)
	if err != nil {
		return err
	}
	return nil
}

func (s *EventStorage) DeleteEvent(eventID string) error {
	ctx := context.Background()
	deleteEventQuery := `delete from events where id = $1`
	_, err := s.db.ExecContext(ctx, deleteEventQuery, eventID)
	if err != nil {
		return err
	}
	return nil
}

func (s *EventStorage) ListEventsForDay(date string) ([]storage.Event, error) {
	ctx := context.Background()
	searchDate, err := utils.ParseDate(date)
	if err != nil {
		return nil, err
	}
	return searchEventsOverPeriod(ctx, searchDate, SearchPeriod{days: 1}, s.db)
}

func (s *EventStorage) ListEventsForWeek(date string) ([]storage.Event, error) {
	ctx := context.Background()
	searchDate, err := utils.ParseDate(date)
	if err != nil {
		return nil, err
	}
	return searchEventsOverPeriod(ctx, searchDate, SearchPeriod{days: 7}, s.db)
}

func (s *EventStorage) ListEventsForMonth(date string) ([]storage.Event, error) {
	ctx := context.Background()
	searchDate, err := utils.ParseDate(date)
	if err != nil {
		return nil, err
	}
	return searchEventsOverPeriod(ctx, searchDate, SearchPeriod{months: 7}, s.db)
}

func searchEventsOverPeriod(
	ctx context.Context,
	startDate time.Time,
	period SearchPeriod,
	db *sqlx.DB,
) ([]storage.Event, error) {
	eventList := make([]storage.Event, 0)
	getEventsOverPeriodQuery := `select * from events where "date" >= $1 and "date" < $2`
	rows, err := db.QueryContext(ctx, getEventsOverPeriodQuery, startDate.Format("2006-01-02"),
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
	checkCountEventOnDateQuery := `select id from events where "date" == $1`
	rows, err := db.QueryContext(ctx, checkCountEventOnDateQuery, date)
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
