// Package roaringIndex provides ability to index objects and search objects with given properties.
// Roaring bitmaps are used as indices.
package roaringIndex

import (
	"github.com/RoaringBitmap/roaring"
)

// IndexedObject is an object prepared o be put in index
// Properties is a set of properties which object has
type IndexedObject struct {
	Object     interface{}
	Properties []string
}

// Index provides ability to find objects with desired properties
type Index struct {
	objects    []interface{}
	properties map[string]*roaring.Bitmap
	fullSet    *roaring.Bitmap
}
