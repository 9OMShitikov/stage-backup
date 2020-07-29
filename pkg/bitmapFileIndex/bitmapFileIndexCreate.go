package bitmapFileIndex

import (
	"fmt"
	"github.com/neganovalexey/bitmap-search/pkg/roaringIndex"
	"io/ioutil"
	"strings"
	"unicode"
)

func nonWordSymbol(r rune) bool {
	return !(unicode.IsLetter(r) || unicode.IsNumber(r))
}

func CreateFileIndex (filenames ...string) (FileIndex, error) {
//	filesWords := make([][]string, len(filenames))
	filesWithWords := make([]roaringIndex.IndexedObject, len(filenames))

	for i, name := range filenames {
		contents, err := ioutil.ReadFile(name)
		if err != nil {
			return FileIndex{}, fmt.Errorf("error while reading %s: %w", name, err)
		}
		wordsSet := make(map[string] struct{})
		contentsAsString := string(contents)
		contentsWords := strings.FieldsFunc(contentsAsString, nonWordSymbol)
		for _, word := range contentsWords {
			wordsSet[word] = struct{}{}
		}
		words := make([]string, 0, len(wordsSet))
		for word, _ := range wordsSet {
			words = append(words, word)
		}

		filesWithWords[i] = roaringIndex.IndexedObject{name, words}
	}

	return FileIndex{roaringIndex.CreateIndexFrom(filesWithWords...)}, nil
}