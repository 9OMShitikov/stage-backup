package roaringIndex

import (
	"github.com/RoaringBitmap/roaring"
	"time"
)

func (index *Index) peekQueries(queries []string) []*roaring.Bitmap {
	peekSets := make([]*roaring.Bitmap, 0, len(queries))
	for _, query := range queries {
		querySet, ok := index.properties[query]
		if ok {
			peekSets = append(peekSets, querySet)
		} else {
			peekSets = append(peekSets, roaring.New())
		}
	}
	return peekSets
}

func (index *Index) peekObjects(set *roaring.Bitmap) []interface{} {
	objects := make([]interface{}, set.GetCardinality())
	for it, j := set.Iterator(), 0; it.HasNext(); j++ {
		objects[j] = index.objects[it.Next()]
	}
	return objects
}

func (index *Index) withAllBitmap(queries []string) *roaring.Bitmap {
	toIntersect := index.peekQueries(queries)
	return roaring.FastAnd(toIntersect...)
}

func (index *Index) withAnyBitmap(queries []string) *roaring.Bitmap {
	toUnite := index.peekQueries(queries)
	return roaring.FastOr(toUnite...)
}

func (index *Index) complexQueryBitmap(withAll []string, withAny[]string, withoutAny[]string) *roaring.Bitmap {
	resultSets := make([]*roaring.Bitmap, 1, 4)
	resultSets[0] = index.fullSet.Clone()
	if withAll != nil {
		resultSets = append(resultSets, index.withAllBitmap(withAll))
	}
	if withAny != nil {
		resultSets = append(resultSets, index.withAnyBitmap(withAny))
	}
	if withoutAny != nil {
		resultSets = append(resultSets, index.withoutAnyBitmap(withoutAny))
	}
	return roaring.FastAnd(resultSets...)
}

func (index *Index) withoutAnyBitmap(queries []string) *roaring.Bitmap {
	toUnite := index.peekQueries(queries)
	return roaring.AndNot(index.fullSet, roaring.FastOr(toUnite...))
}

// WithAll searches all objects in index with all of specified properties
func (index *Index) WithAll(queries ...string) []interface{} {
	resultSet := index.withAllBitmap(queries)
	return index.peekObjects(resultSet)
}

// WithAny searches all objects in index with at least one of specified properties
func (index *Index) WithAny(queries ...string) []interface{} {
	resultSet := index.withAnyBitmap(queries)
	return index.peekObjects(resultSet)
}

// WithoutAny searches all objects in index without any of specified properties
func (index *Index) WithoutAny(queries ...string) []interface{} {
	resultSet := index.withoutAnyBitmap(queries)
	return index.peekObjects(resultSet)
}

// FromTo searches all objects in index with datetime in [from, to) interval
func (index *Index) FromTo(from, to time.Time) []interface{} {
	resultSet := index.timeIndex.peekIndices(from, to)
	return index.peekObjects(resultSet)
}

// ComplexQuery searches all objects in index with all, with any and without any of
// specified properties (argument should be nil if you want to exclude that part of query)
func (index *Index) ComplexQuery(withAll []string, withAny []string,
	withoutAny []string) []interface{} {
	return index.peekObjects(index.complexQueryBitmap(withAll, withAny, withoutAny))
}

// ComplexQuery searches all objects in index with all, with any and without any of
// specified properties (argument should be nil if you want to exclude that part of query) with datetime in [from, to)
func (index *Index) ComplexQueryDatetime(withAll []string, withAny []string,
	withoutAny []string, from, to time.Time) []interface{} {
	resultSet := index.complexQueryBitmap(withAll, withAny, withoutAny)
	resultSet.And(index.timeIndex.peekIndices(from, to))
	return index.peekObjects(resultSet)
}