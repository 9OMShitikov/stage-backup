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

// IndexedObject is an object prepared o be put in index
// Properties is a set of properties which object has
// Datetime is an object datetime
type IndexedObjectWithDatetime struct {
	Object     interface{}
	Properties []string
	Datetime time.Time
}

// Index provides ability to find objects with desired properties
type Index struct {
	objects    []interface{}
	properties map[string]*roaring.Bitmap
	datetime []*roaring.Bitmap
	fullSet    *roaring.Bitmap
}
