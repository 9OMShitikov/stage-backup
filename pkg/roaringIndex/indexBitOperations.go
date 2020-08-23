package roaringIndex

import "github.com/RoaringBitmap/roaring"

func (index *Index) complement (set *roaring.Bitmap) *roaring.Bitmap {
	result := index.fullSet.Clone()
	result.AndNot(set)
	return result
}
