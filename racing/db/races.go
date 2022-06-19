package db

import (
	"database/sql"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error)
	UpdateAllRacesByColumn(races []*racing.Race)
	GetRace(race_id int64) ([]*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error
	var columnName string = "status"

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
		err = r.addColumnToRacesTable(columnName)
	})

	return err
}

/*
races: Races database in array.
columnName: Name of the column in Races database.
*/
func (r *racesRepo) UpdateAllRacesByColumn(races []*racing.Race) {
	// Column names within the array will be updated using switch statements to configure what types of update is necessary.
	columnNames := []string{
		"status",
	}

	// Can be configurable to add in different types of column names to update the database.
	for _, columnName := range columnNames {
		switch columnName {
		case "status":
			r.updateRaceStatus(races, columnName)
		}
	}
}

func (r *racesRepo) updateRaceStatus(races []*racing.Race, columnName string) {
	var currentTime = time.Now()
	const (
		raceStatusOpen   string = "OPEN"
		raceStatusClosed string = "CLOSED"
	)

	for _, race := range races {
		raceID := strconv.Itoa(int(race.GetId()))

		// Update the Races database status column after checking the GetAdvertisedStartTime().
		if race.GetAdvertisedStartTime().AsTime().Before(currentTime) {
			r.updateRacesTable(columnName, raceStatusClosed, raceID)
		} else {
			r.updateRacesTable(columnName, raceStatusOpen, raceID)
		}
	}
}

func (r *racesRepo) GetRace(race_id int64) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]
	query, args = r.applyFindRaceById(query, race_id)

	row, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(row)
}

func (r *racesRepo) List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]
	query, args = r.applyFilter(query, filter)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

// Build up query call for database search for a single race request.
func (r *racesRepo) applyFindRaceById(query string, race_id int64) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if race_id <= 0 {
		return query, args
	}

	// TODO: verify with proper gateway built.
	clauses = append(clauses, "id IN (?,?,?)")
	args = append(args, race_id)

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}

	if len(filter.RaceVisibility) > 0 {
		clauses = append(clauses, "visible IN ("+strings.Repeat("?,", len(filter.RaceVisibility)-1)+"?)")

		for _, raceVisibility := range filter.RaceVisibility {
			args = append(args, raceVisibility)
		}
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func (m *racesRepo) scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart, &race.Status); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		race.AdvertisedStartTime = ts

		races = append(races, &race)
	}

	return races, nil
}
