package storage

type Event struct {
	ID          string `db:"id"`
	Title       string
	Date        string
	Duration    int
	Description string
	UserID      string `db:"user_id"`
	Reminder    string
}
