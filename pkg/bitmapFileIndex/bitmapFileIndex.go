package bitmapFileIndex

import (
	"github.com/neganovalexey/bitmap-search/pkg/roaringIndex"
)

type FileIndex struct {
	index *roaringIndex.Index
}
