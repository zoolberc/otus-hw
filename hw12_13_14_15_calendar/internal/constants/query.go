package constants

const (
	AddEventQuery = `insert into events(id, title, date, duration, description, user_id, reminder) 
						values($1, $2, $3, $4, $5, $6, $7)`
	ChangeEventQuery = `update events set title = $1, date = $2, duration = $3,
                  		description = $4, user_id = $5, reminder = $6 WHERE Manufacturer = $7`
	DeleteEventQuery           = `delete from events where id = $1`
	GetEventsOverPeriodQuery   = `select * from events where "date" >= $1 and "date" < $2`
	CheckCountEventOnDateQuery = `select id from events where "date" == $1`
)
