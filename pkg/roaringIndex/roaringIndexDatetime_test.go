package roaringIndex_test

import (
	. "github.com/neganovalexey/bitmap-search/pkg/roaringIndex"
	"math/rand"
	"sort"
	"testing"
	"time"
)

//Check if slices l and r contain equal sets of timestamps and r doesn't contain duplicates
func equalTimestampSetsWithRAsInterfacesWithoutDup(l []time.Time, r []interface{}) bool {
	set := make(map[time.Time]struct{})

	for _, elementL := range l {
		set[elementL] = struct{}{}
	}

	for _, elementR := range r {
		elementRStr := elementR.(time.Time)
		if _, ok := set[elementRStr]; ok {
			delete(set, elementRStr)
		} else {
			return false
		}
	}
	if len(set) != 0 {
		return false
	}
	return true
}

func searchDatetimeInSortedSlice (slice []IndexedObjectWithDatetime, from, to time.Time) []time.Time {
	if !from.Before(to) {
		return make([]time.Time, 0)
	}

	start := sort.Search(len(slice), func(n int) bool {
		return !slice[n].Datetime.Before(from)
	})

	finish := sort.Search(len(slice), func(n int) bool {
		return !slice[n].Datetime.Before(to)
	})

	searchResult := make([]time.Time, finish - start)
	for i := range searchResult {
		searchResult[i] = slice[start + i].Datetime
	}
	return searchResult
}

func createTestObjectFromDatetime(timestamp time.Time) IndexedObjectWithDatetime {
	return CreateIndexedObjectWithDatetime(timestamp, timestamp)
}

func createDatetimeTestContainment() []IndexedObjectWithDatetime {
	testSet := make([]IndexedObjectWithDatetime, 0, 120000)
	for year := MinTimestampYear; year < MaxTimestampYear; year += 120 {
		for month := 1; month <= 12; month += 11 {
			for day := 1; day <= 31; day += 29 {
				for minute := 0; minute < 60; minute += 30 {
					for second := 0; second < 60; second += 3 {
						testSet = append(testSet, createTestObjectFromDatetime(
							time.Date(year, time.Month(month), day, 5, minute, second, 0, time.UTC)))
					}
				}
			}
		}
	}

	testSet = testSet
	rand.Shuffle(len(testSet), func(i int, j int) {
		testSet[i], testSet[j] = testSet[j], testSet[i]
	})
	return testSet
}

func createDatetimeTestQueries() []time.Time {
	testSet := make([]time.Time, 0)
	for year := MinTimestampYear; year < MaxTimestampYear; year += 60 {
		for month := 1; month <= 12; month += 11 {
			for day := 1; day <= 31; day += 29 {
				for minute := 0; minute < 60; minute += 30 {
					for second := 0; second < 60; second += 3 {
						testSet = append(testSet, time.Date(year, time.Month(month),
							day, 5, minute, second, 0, time.UTC))
					}
				}
			}
		}
	}
	return testSet
}

func TestCreateIndexWithDatetimeFrom(t *testing.T) {
	testContainment := createDatetimeTestContainment()
	index, err := CreateIndexWithDatetimeFrom(testContainment...)
	if err != nil {
		t.Error("got \"", err, "\" error while creating index")
	}
	if index == nil {
		t.Error("didn't receive index")
	}
}

func TestCreateIndexWithDatetimeFrom_error(t *testing.T) {
	index, err := CreateIndexWithDatetimeFrom(createTestObjectFromDatetime(time.Date(1000, 1,
		1, 1, 1, 1, 1, time.UTC)))
	if err == nil {
		t.Error("no errors while creating index while timestamp year is too small")
	}
	if index != nil {
		t.Error("didn't receive nil index while timestamp year is too small")
	}

	index, err = CreateIndexWithDatetimeFrom(createTestObjectFromDatetime(time.Date(4000, 1,
		1, 1, 1, 1, 1, time.UTC)))
	if err == nil {
		t.Error("no errors while creating index while timestamp year was too big")
	}
	if index != nil {
		t.Error("didn't receive nil index while timestamp year was too big")
	}
}

func TestIndex_AtInterval(t *testing.T) {
	containmentSet := createDatetimeTestContainment()
	queriesSet := createDatetimeTestQueries()

	index, err := CreateIndexWithDatetimeFrom(containmentSet...)
	if err != nil {
		t.Error("got \"", err, "\" error while creating index")
	}
	if index == nil {
		t.Error("didn't receive index")
	}

	sort.Slice(containmentSet, func(i int, j int) bool {
		return containmentSet[i].Datetime.Before(containmentSet[j].Datetime)
	})

	for _, fromStamp := range queriesSet {
		for _, toStamp := range queriesSet {
			basic := searchDatetimeInSortedSlice(containmentSet, fromStamp, toStamp)
			checked, err := index.AtInterval(fromStamp, toStamp)
			if err != nil {
				t.Error("got \"", err, "\" error while making \"from \"", fromStamp,
					"\" to \"", toStamp, "\"\" query")
			}
			if !equalTimestampSetsWithRAsInterfacesWithoutDup(basic, checked) {
				t.Error("expected ", basic, " got ", checked)
			}
		}
	}
}