package roaringIndex

import (
	"github.com/RoaringBitmap/roaring"
	"time"
)

type datetimeIndex interface {
	peekIndices(from, to time.Time) *roaring.Bitmap
}