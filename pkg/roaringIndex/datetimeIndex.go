package roaringIndex

import (
	"github.com/RoaringBitmap/roaring"
	"time"
)

type datetimeIndex interface {
	peekIndices(from, to time.Time) *roaring.Bitmap
}

type standardCalendarIndex struct {
	times [][] *roaring.Bitmap
}

func getDigits(num uint, digits []uint) {
	for i := len(digits) - 1; i >= 0; i-- {
		digits[i] = num % 10
		num = num / 10
	}
}

func (*standardCalendarIndex) datetimeTransform (datetime time.Time) []uint {
	timeSlice := make([]uint, 13)
	timeSlice[4] = uint(datetime.Month())
	dateDigits := []uint{uint(datetime.Second()),
		uint(datetime.Minute()),
		uint(datetime.Hour()),
		uint(datetime.Day()),
		uint(datetime.Year())}
	digitsSlices := [][]uint{timeSlice[11:13],
		timeSlice[9:11],
		timeSlice[7:9],
		timeSlice[5:7],
		timeSlice[0:4]}
	for i := range dateDigits {
		getDigits(dateDigits[i], digitsSlices[i])
	}
	return timeSlice
}

func (index *standardCalendarIndex) peekIndices(from, to time.Time) *roaring.Bitmap {
	if !to.After(from) {
		return roaring.New()
	}

	fromSlice, toSlice := index.datetimeTransform(from), index.datetimeTransform(to)

	for i := range fromSlice {

	}
}
