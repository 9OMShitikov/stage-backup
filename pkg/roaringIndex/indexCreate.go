package roaringIndex

import (
	"github.com/RoaringBitmap/roaring"
	"time"
)

// CreateIndexedObject creates IndexedObject from specified object and properties
func CreateIndexedObject(object interface{}, properties ...string) IndexedObject {
	return IndexedObject{object, properties}
}

// CreateIndexFrom creates Index from IndexedObjects
// Datetime queries on these objects don't return correct results
func CreateIndexFrom(objects ...IndexedObject) *Index {
	index := new(Index)

	index.objects = make([]interface{}, len(objects))
	index.fullSet = roaring.New()
	index.properties = make(map[string]*roaring.Bitmap)

	objectSets := make(map[string][]uint32)
	for i, object := range objects {
		index.objects[i] = object.Object

		for _, property := range object.Properties {
			objectSets[property] = append(objectSets[property], uint32(i))
		}
	}

	for property, objectSet := range objectSets {
		index.properties[property] = roaring.BitmapOf(objectSet...)
	}

	index.fullSet.AddRange(uint64(0), uint64(len(objects)))
	index.timeIndex = nil
	
	return index
}

func createDigitIndexFrom(transformer dateToDigits, objects ...IndexedObjectWithDatetime) *Index {
	objectsTable := make([]IndexedObject, len(objects))
	timeTable := make([]time.Time, len(objects))

	for i, object := range objects {
		objectsTable[i] = CreateIndexedObject(object.Object, object.Properties...)
		timeTable[i] = object.Datetime
	}
	index := CreateIndexFrom(objectsTable...)
	index.timeIndex = createComplexIndexFrom(timeTable, index.fullSet, transformer)
	return index
}

func CreateCalendarIndexFrom(objects ...IndexedObjectWithDatetime) *Index {
	return createDigitIndexFrom(standardCalendarIndex{}, objects...)
}