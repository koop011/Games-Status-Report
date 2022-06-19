package service

import (
	"sort"

	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
	"golang.org/x/net/context"
)

type Racing interface {
	// ListRaces will return a collection of races.
	ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error)
	GetRace(ctx context.Context, in *racing.GetRaceRequest) (*racing.GetRaceResponse, error)
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
	races, err := s.racesRepo.List(in.Filter)
	if err != nil {
		return nil, err
	}
	raceReportSort(races, in)
	s.racesRepo.UpdateAllRacesByColumn(races)

	return &racing.ListRacesResponse{Races: races}, nil
}

// Allow user to sort based on SortOrder message.
func raceReportSort(races []*racing.Race, in *racing.ListRacesRequest) []*racing.Race {
	orderedItem := in.GetFilter().GetSortOrder().GetOrderedItem()
	lowToHigh := in.GetFilter().GetSortOrder().GetLowToHigh()
	// TODO: implement alphabetical order sort
	// alphabetical := in.GetFilter().GetSortOrder().GetAlphabetical()
	earliest_to_latest := in.GetFilter().GetSortOrder().GetEarliestToLatest()
	const forced_default_earliest_to_latest = true

	switch orderedItem {
	case "start time":
		sortByAdvertisedStartTime(races, earliest_to_latest)
	case "id":
		sortById(races, lowToHigh)
	case "meetingid":
		sortByMeetingId(races, lowToHigh)
	case "number":
		sortByNumber(races, lowToHigh)
	// TODO: implement alphabetical order sort
	// case "name":
	// 	sortByName(races, alphabetical)
	default:
		sortByAdvertisedStartTime(races, forced_default_earliest_to_latest)
	}

	return races
}

func sortByNumber(races []*racing.Race, lowToHigh bool) []*racing.Race {
	if lowToHigh {
		sort.SliceStable(races, func(i, j int) bool {
			return races[i].GetNumber() < races[j].GetNumber()
		})
	} else {
		sort.SliceStable(races, func(i, j int) bool {
			return races[i].GetNumber() > races[j].GetNumber()
		})
	}

	return races
}

func sortByMeetingId(races []*racing.Race, lowToHigh bool) []*racing.Race {
	if lowToHigh {
		sort.SliceStable(races, func(i, j int) bool {
			return races[i].GetMeetingId() < races[j].GetMeetingId()
		})
	} else {
		sort.SliceStable(races, func(i, j int) bool {
			return races[i].GetMeetingId() > races[j].GetMeetingId()
		})
	}

	return races
}

func sortById(races []*racing.Race, lowToHigh bool) []*racing.Race {
	if lowToHigh {
		sort.SliceStable(races, func(i, j int) bool {
			return races[i].GetId() < races[j].GetId()
		})
	} else {
		sort.SliceStable(races, func(i, j int) bool {
			return races[i].GetId() > races[j].GetId()
		})
	}

	return races
}

func sortByAdvertisedStartTime(races []*racing.Race, earliest_to_latest bool) []*racing.Race {
	// due to timestamppb package, the race[].GetAdvertisedStartTime() is converted again to UTC before sorting to
	// earliest to latest.
	if earliest_to_latest {
		sort.SliceStable(races, func(i, j int) bool {
			return races[i].GetAdvertisedStartTime().AsTime().Before(races[j].GetAdvertisedStartTime().AsTime())
		})
	} else {
		sort.SliceStable(races, func(i, j int) bool {
			return races[j].GetAdvertisedStartTime().AsTime().Before(races[i].GetAdvertisedStartTime().AsTime())
		})
	}

	return races
}

func (s *racingService) GetRace(ctx context.Context, in *racing.GetRaceRequest) (*racing.GetRaceResponse, error) {
	race, err := s.racesRepo.GetRace(in.RaceId)

	if err != nil {
		return nil, err
	}

	return &racing.GetRaceResponse{Race: race}, nil
}
