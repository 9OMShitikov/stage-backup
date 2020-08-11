package roaringIndex

import (
	"time"
)

type standardCalendarIndex struct {
}

func getDigits(number uint, slice []uint) {
	for i := len(slice) - 1; i >= 0; i-- {
		slice[i] = number % 10
		number /= 10
	}
}

func (standardCalendarIndex) datetimeTransform (datetime time.Time) []uint {
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

func (standardCalendarIndex) getDigitsSizes() []uint {
	return []uint{9, 9, 9, 9, 11, 3, 9, 2, 9, 6, 9, 6, 9}
}
