package service

import (
	"log"
	"sort"

	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
	"golang.org/x/net/context"
)

type Racing interface {
	// ListRaces will return a collection of races.
	ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error)
}

// racingService implements the Racing interface.
type racingService struct {
	racesRepo db.RacesRepo
}

// NewRacingService instantiates and returns a new racingService.
func NewRacingService(racesRepo db.RacesRepo) Racing {
	return &racingService{racesRepo}
}

func (s *racingService) ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error) {
	// TODO: implement string array for multiple column creation
	var statusColumn string = "status"

	races, err := s.racesRepo.List(in.Filter)
	if err != nil {
		return nil, err
	}
	raceReportSortByAdvertisedTime(races)
	s.racesRepo.UpdateAllRacesByColumn(races, statusColumn)
	log.Print(races)
	return &racing.ListRacesResponse{Races: races}, nil
}

// TODO: update function to parse in user preference of sort order via filter sortByTime, 'sort-low-high' and 'sort-high-low'
func raceReportSortByAdvertisedTime(races []*racing.Race) []*racing.Race {
	// due to timestamppb package, the race[].GetAdvertisedStartTime() is converted again to UTC before sorting to
	// earliest to latest.
	sort.SliceStable(races, func(i, j int) bool {
		return races[i].GetAdvertisedStartTime().AsTime().Before(races[j].GetAdvertisedStartTime().AsTime())
	})
	return races
}
