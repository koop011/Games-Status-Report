package db

import (
	"time"

	"syreclabs.com/go/faker"
)

func (r *racesRepo) seed() error {
	statement, err := r.db.Prepare(`CREATE TABLE IF NOT EXISTS races (id INTEGER PRIMARY KEY, meeting_id INTEGER, name TEXT, number INTEGER, visible INTEGER, advertised_start_time DATETIME)`)
	if err == nil {
		_, err = statement.Exec()
	}

	for i := 1; i <= 100; i++ {
		statement, err = r.db.Prepare(`INSERT OR IGNORE INTO races(id, meeting_id, name, number, visible, advertised_start_time) VALUES (?,?,?,?,?,?)`)
		if err == nil {
			_, err = statement.Exec(
				i,
				faker.Number().Between(1, 10),
				faker.Team().Name(),
				faker.Number().Between(1, 12),
				faker.Number().Between(0, 1),
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)).Format(time.RFC3339),
			)
		}
	}

	return err
}

/*
columnName: name of the column to be added to existing Races database.

Called during racesRepo.Init().
*/
func (r *racesRepo) addColumnToRacesTable(columnName string) error {

	statement, err := r.db.Prepare(`ALTER TABLE races ADD ` + columnName + ` TEXT`)
	if err == nil {
		_, err = statement.Exec()
	} else {
		/*
			This is a poor implementation of IF NOT EXIST call as it doesn't exist in sqlite3.
			The above statement.Exec() will cause an error 'duplicate column name: status', if the source has be run previously.
			This will ignore the duplicate column name, but may be vulnerable to other error statements.
		*/
		// TODO: update sql call with better error handling.
		err = nil
	}

	return err
}

/*
columnName: Name of the column in Races database.
raceStatus: "OPEN"|"CLOSED".
id: Races database id key value.

Called when receiving ListRacesRequest message to update the Races database using columnName and id key value of each races to insert the raceStatus
before ListRacesResponse is sent.
*/
func (r *racesRepo) updateRacesTable(columnName string, raceStatus string, id string) error {
	var err error

	statement, err := r.db.Prepare(`UPDATE races SET ` + columnName + ` = "` + raceStatus + `" WHERE id = ` + id)
	if err == nil {
		_, err = statement.Exec()
	}

	return err
}
