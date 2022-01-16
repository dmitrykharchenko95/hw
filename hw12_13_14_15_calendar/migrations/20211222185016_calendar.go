package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCalendar, downCalendar)
}

func upCalendar(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE events(
		id serial NOT NULL PRIMARY KEY,
		title TEXT NOT NULL,
		start_date TIMESTAMP NOT NULL,
		end_date TIMESTAMP NOT NULL,
		content TEXT,
		user_id INT NOT NULL,
		send_time BIGINT NOT NULL DEFAULT 0);`)
	if err != nil {
		return err
	}
	return nil
}

func downCalendar(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE events;`)
	if err != nil {
		return err
	}

	return nil
}
