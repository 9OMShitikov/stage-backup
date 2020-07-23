package roaringIndex

import (
	"github.com/RoaringBitmap/roaring"
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

// WithAll searches all objects in index with all of specified properties
func (index *Index) WithAll(queries ...string) []interface{} {
	toIntersect := index.peekQueries(queries)
	resultSet := roaring.FastAnd(toIntersect...)
	return index.peekObjects(resultSet)
}

// WithAny searches all objects in index with at least one of specified properties
func (index *Index) WithAny(queries ...string) []interface{} {
	toUnite := index.peekQueries(queries)
	resultSet := roaring.FastOr(toUnite...)
	return index.peekObjects(resultSet)
}

// WithoutAny searches all objects in index without any of specified properties
func (index *Index) WithoutAny(queries ...string) []interface{} {
	toUnite := index.peekQueries(queries)
	resultSet := roaring.AndNot(index.fullSet, roaring.FastOr(toUnite...))
	return index.peekObjects(resultSet)
}
