package roaringIndex

import (
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"time"
)

func getBits(num uint, bits []bool) {
	for i := len(bits) - 1; i >= 0; i-- {
		bits[i] = num & 1 != 0
		num >>= 1
	}
}

const (
	datetimeSliceLen = 8 + 4 + 5 + 6 + 6 + 6
	MinTimestampYear = 1900
	MaxTimestampYear = MinTimestampYear + 256
)

func datetimeToBitSlice (datetime time.Time) []bool {
	timeSlice := make([]bool, datetimeSliceLen)
	dateDigits := []uint{uint(datetime.Second()),
		uint(datetime.Minute()),
		uint(datetime.Hour()),
		uint(datetime.Day()),
		uint(datetime.Month()),
		uint(datetime.Year() - MinTimestampYear)}
	digitsSlices := [][]bool{
		timeSlice[29:35],
		timeSlice[23:29],
		timeSlice[17:23],
		timeSlice[12:17],
		timeSlice[8:12],
		timeSlice[0:8]}
	for i := range dateDigits {
		getBits(dateDigits[i], digitsSlices[i])
	}
	return timeSlice
}

func checkTimestampRange (timestamp time.Time) bool {
	return timestamp.Year() < MinTimestampYear ||
		timestamp.Year() >= MaxTimestampYear
}

func checkZeroTimestamp (timestamp time.Time) bool {
	return timestamp.Year() == 0 && timestamp.Month() == 0 && timestamp.Day() == 0 &&
		timestamp.Hour() == 0 && timestamp.Minute() == 0 && timestamp.Second() == 0
}

func (index *Index) peekIndices(from, to time.Time) (*roaring.Bitmap, error) {
	if checkTimestampRange(from) {
		return nil, fmt.Errorf("incorrect from year: %d", from.Year())
	}
	if checkTimestampRange(to) {
		return nil, fmt.Errorf("incorrect to year: %d", from.Year())
	}
	if checkZeroTimestamp(from) {
		return nil, fmt.Errorf("zero from timestamp")
	}
	if checkZeroTimestamp(to) {
		return nil, fmt.Errorf("zero to timestamp")
	}

	if !to.After(from) {
		return roaring.New(), nil
	}

	fromSlice, toSlice := datetimeToBitSlice(from), datetimeToBitSlice(to)
	datetimesDiffer := false
	fromBorder := index.fullSet.Clone()
	var toBorder *roaring.Bitmap = nil
	resultSets := make([]*roaring.Bitmap, 0, 2 * datetimeSliceLen)

	for i := range fromSlice {
		if datetimesDiffer {
			if fromSlice[i] {
				fromBorder.And(index.datetime[i])
			} else {
				addedSlice := fromBorder.Clone()
				addedSlice.And(index.datetime[i])
				resultSets = append(resultSets, addedSlice)

				fromBorder.AndNot(index.datetime[i])
			}

			if toSlice[i] {
				addedSlice := toBorder.Clone()
				addedSlice.AndNot(index.datetime[i])
				resultSets = append(resultSets, addedSlice)

				toBorder.And(index.datetime[i])
			} else {
				toBorder.AndNot(index.datetime[i])
			}
		} else {
			if fromSlice[i] == toSlice[i] {
				if fromSlice[i] {
					fromBorder.And(index.datetime[i])
				} else {
					fromBorder.AndNot(index.datetime[i])
				}
			} else {
				datetimesDiffer = true
				toBorder = fromBorder.Clone()

				fromBorder.AndNot(index.datetime[i])
				toBorder.And(index.datetime[i])
			}
		}
	}

	resultSets = append(resultSets, fromBorder)
	return roaring.FastOr(resultSets...), nil
}
