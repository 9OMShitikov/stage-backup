package roaringIndex

import (
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"time"
)

// CreateIndexedObject creates IndexedObject from specified object and properties
func CreateIndexedObject(object interface{}, properties ...string) IndexedObject {
	return IndexedObject{object, properties}
}

// CreateIndexedObjectWithDatetime creates IndexedObjectWithDatetime from specified object, properties and datetime
func CreateIndexedObjectWithDatetime(object interface{},
                                     timestamp time.Time, properties ...string) IndexedObjectWithDatetime {
	return IndexedObjectWithDatetime{object, properties, timestamp}
}

// CreateIndexFrom creates Index from IndexedObjects
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
	index.datetime = nil
	return index
}

// CreateIndexWithDatetimeFrom creates Index from IndexedObjectWithDatetimes (will crush if any timestamp is 00.00.0000
// or isn't in [1900, 2156))
func CreateIndexWithDatetimeFrom(objects ...IndexedObjectWithDatetime) (*Index, error) {
	timestamps := make([][]uint32, datetimeSliceLen)
	objectsWithProperties := make([]IndexedObject, len(objects))

	for i := range objects {
		timestamp := objects[i].Datetime
		if checkTimestampRange(timestamp) {
			return nil, fmt.Errorf("incorrect year %d at object \"%s\"", objects[i].Datetime.Year(),
				objects[i].Object)
		}
	}

	for i := range objects {
		datetimeAsSlice := datetimeToBitSlice(objects[i].Datetime)
		for j, bit := range datetimeAsSlice {
			if bit {
				timestamps[j] = append(timestamps[j], uint32(i))
			}
		}

		objectsWithProperties[i] = CreateIndexedObject(objects[i].Object, objects[i].Properties...)
	}

	index := CreateIndexFrom(objectsWithProperties...)
	index.datetime = make([]*roaring.Bitmap, datetimeSliceLen)

	for i := range index.datetime {
		index.datetime[i] = roaring.BitmapOf(timestamps[i]...)
	}
	return index, nil
}