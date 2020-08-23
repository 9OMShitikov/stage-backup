package roaringIndex

import (
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"time"
)

func (index *Index) peekObjects(set *roaring.Bitmap) []interface{} {
	objects := make([]interface{}, set.GetCardinality())
	for it, j := set.Iterator(), 0; it.HasNext(); j++ {
		objects[j] = index.objects[it.Next()]
	}
	return objects
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

// AtInterval searches all objects in [from, to) interval (will crush if from or to is 00.00.0000 or
// isn't in [1900, 2156))
func (index* Index) AtInterval (from, to time.Time) ([]interface{}, error) {
	resultSet, err := index.peekIndices(from, to)
	if err != nil {
		return nil, fmt.Errorf("error while finding objects at interval: %w", err)
	}

	return index.peekObjects(resultSet), nil
}