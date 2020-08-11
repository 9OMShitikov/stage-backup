package roaringIndex

import (
	"github.com/RoaringBitmap/roaring"
	"time"
)

type dateToDigits interface {
	getDigitsSizes () []uint
	datetimeTransform (datetime time.Time) []uint
}

type datetimeComplexDigitIndex struct {
	transformer dateToDigits
	times [][] *roaring.Bitmap
	zero *roaring.Bitmap
}

func (index *datetimeComplexDigitIndex) peekRangeFromDigit (digit, from, to uint) *roaring.Bitmap {
	digitSlice := index.times[digit]
	if (from >= to) || (from > uint(len(digitSlice))) {
		return roaring.New()
	}

	var incl *roaring.Bitmap = nil

	if from == 0 {
		incl = index.zero.Clone()
	} else {
		incl = digitSlice[from - 1].Clone()
	}

	if to > uint(len(digitSlice)) {
		return incl
	} else {
		incl.AndNot(digitSlice[to - 1])
		return incl
	}
}

func (index *datetimeComplexDigitIndex) peekDigit (digit, taken uint) *roaring.Bitmap {
	digitSlice := index.times[digit]
	if taken > uint(len(digitSlice)) {
		return roaring.New()
	}

	var incl *roaring.Bitmap = nil
	if taken == 0 {
		incl = index.zero.Clone()
	} else {
		incl = digitSlice[taken - 1].Clone()
	}

	if taken == uint(len(digitSlice)) {
		return incl
	} else {
		incl.AndNot(digitSlice[taken])
		return incl
	}
}

func (index *datetimeComplexDigitIndex) digitUpperEqual(digit, border uint) *roaring.Bitmap {
	digitSlice := index.times[digit]

	switch {
	case border == 0:
		return index.zero.Clone()
	case border > uint(len(digitSlice)):
		return roaring.New()
	default:
		return digitSlice[border - 1]
	}
}

func (index *datetimeComplexDigitIndex) digitUpper(digit, border uint) *roaring.Bitmap {
	digitSlice := index.times[digit]

	if border >= uint(len(digitSlice)) {
		return roaring.New()
	} else {
		return digitSlice[border]
	}
}

func (index *datetimeComplexDigitIndex) digitLower(digit, border uint) *roaring.Bitmap {
	digitSlice := index.times[digit]

	switch {
	case border == 0:
		return roaring.New()
	case border > uint(len(digitSlice)):
		return index.zero.Clone()
	default:
		result := index.zero.Clone()
		result.AndNot(digitSlice[border - 1])
		return result
	}
}

func (index *datetimeComplexDigitIndex) peekIndices(from, to time.Time) *roaring.Bitmap {
	if !to.After(from) {
		return roaring.New()
	}

	fromSlice, toSlice := index.transformer.datetimeTransform(from), index.transformer.datetimeTransform(to)

	areDifferent := false

	resultsToMerge := make([]*roaring.Bitmap, 0, len(index.times))

	commonDigits := make([]*roaring.Bitmap, 1, len(index.times))
	commonDigits[0] = index.zero.Clone()

	var fromBorder *roaring.Bitmap
	var toBorder *roaring.Bitmap

	for i := range fromSlice {
		if !areDifferent {
			if fromSlice[i] == toSlice[i] {
				commonDigit := index.peekDigit(uint(i), fromSlice[i])
				commonDigits = append(commonDigits, commonDigit)
			} else {
				areDifferent = true
				toBorder = roaring.FastAnd(commonDigits...)
				fromBorder = toBorder.Clone()

				addedDigits := index.peekRangeFromDigit(uint(i), fromSlice[i] + 1, toSlice[i])
				addedDigits.And(toBorder)
				resultsToMerge = append(resultsToMerge, addedDigits)

				fromBorder.And(index.peekDigit(uint(i), fromSlice[i]))
				toBorder.And(index.peekDigit(uint(i), toSlice[i]))
			}
		} else {
			toAddedSet := toBorder.Clone()
			fromAddedSet := fromBorder.Clone()

			fromAddedSet.And(index.digitUpper(uint(i), fromSlice[i]))
			toAddedSet.And(index.digitLower(uint(i), toSlice[i]))

			resultsToMerge = append(resultsToMerge, fromAddedSet, toAddedSet)

			fromBorder.And(index.peekDigit(uint(i), fromSlice[i]))
			toBorder.And(index.peekDigit(uint(i), toSlice[i]))
		}
	}
	resultsToMerge = append(resultsToMerge, fromBorder)
	return roaring.FastOr(resultsToMerge...)
}

func createComplexIndexFrom (datetimes []time.Time, zero *roaring.Bitmap, transformer dateToDigits) datetimeIndex {
	index := new(datetimeComplexDigitIndex)
	index.zero = zero
	index.transformer = transformer

	sizes := transformer.getDigitsSizes()

	index.times = make([][]*roaring.Bitmap, len(sizes))
	times := make([][][]uint32, len(sizes))
	for i := range index.times {
		index.times[i] = make([]*roaring.Bitmap, sizes[i])
		times[i] = make([][]uint32, sizes[i])
	}

	for i, datetime := range datetimes {
		datetimeAsSlice := transformer.datetimeTransform(datetime)
		for j := range times {
			for k := uint(0); k < datetimeAsSlice[j]; k++ {
				times[j][k] = append(times[j][k], uint32(i))
			}
		}
	}

	for i := range times {
		for j := range times[i] {
			index.times[i][j] = roaring.BitmapOf(times[i][j]...)
		}
	}

	return index
}