roaringIndex
=============
 This Go package provides the ability to index objects with given properties. Package uses [roaring bitmaps](https://roaringbitmap.org) for indexing.
 
 #### Installation
 
   - go get -t github.com/neganovalexey/bitmap-search/...
 
 ### Example
 
 ```go
package main

import (
	"fmt"
	"github.com/neganovalexey/bitmap-search/pkg/roaringIndex"
)

func main() {
	
    var simpleTestObjects = []roaringIndex.IndexedObject{
	    roaringIndex.CreateIndexedObject("|a b c|", "a", "b", "c"),
	    roaringIndex.CreateIndexedObject("|a b|", "a", "b"),
	    roaringIndex.CreateIndexedObject("|a c|", "a", "c"),
	    roaringIndex.CreateIndexedObject("|b c|", "b", "c"),
	    roaringIndex.CreateIndexedObject("|a|", "a"),
	    roaringIndex.CreateIndexedObject("|b|", "b"),
	    roaringIndex.CreateIndexedObject("|c|", "c"),
	    roaringIndex.CreateIndexedObject("|e|", "e"),
    }

	index := roaringIndex.CreateIndexFrom(simpleTestObjects...)

	fmt.Println(index.WithAll("a", "b"))

    fmt.Println(index.WithAny("a", "b", "e"))
	fmt.Println(index.WithoutAny("a", "e"))
}
```