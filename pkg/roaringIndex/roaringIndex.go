// Package roaringIndex provides ability to index objects and search objects with given properties.
// Roaring bitmaps are used as indices.
package roaringIndex

import (
	"github.com/RoaringBitmap/roaring"
	"time"
)

// IndexedObject is an object prepared o be put in index
// Properties is a set of properties which object has
type IndexedObject struct {
	Object     interface{}
	Properties []string
}

// IndexedObjectWithDateTime is an object with timestamp prepared to be put in index
// Datetime is object timestamp
type IndexedObjectWithDatetime struct {
	IndexedObject
	Datetime time.Time
}

// Index provides ability to find objects with desired properties and datetime range
// Objects with zero unix time or datetime 00.00.0000 00:00:00 can be processed wrong at datetime queries
type Index struct {
	objects    []interface{}
	properties map[string] *roaring.Bitmap
	fullSet    *roaring.Bitmap
	timeIndex  datetimeIndex
}