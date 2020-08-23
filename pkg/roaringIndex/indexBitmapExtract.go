package roaringIndex

import "github.com/RoaringBitmap/roaring"

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

func (index *Index) withAllBitmap(queries []string) *roaring.Bitmap {
	toIntersect := index.peekQueries(queries)
	return roaring.FastAnd(toIntersect...)
}

func (index *Index) withAnyBitmap(queries []string) *roaring.Bitmap {
	toUnite := index.peekQueries(queries)
	return roaring.FastOr(toUnite...)
}

func (index *Index) withoutAnyBitmap(queries []string) *roaring.Bitmap {
	toUnite := index.peekQueries(queries)
	return roaring.AndNot(index.fullSet, roaring.FastOr(toUnite...))
}