package bitmapFileIndex

func resultsToStrings (results []interface{}) []string {
	resultsAsStrings := make([]string, len(results))
	for i, result := range results {
		resultsAsStrings[i] = result.(string)
	}
	return resultsAsStrings
}

// WithAll searches all files in index with all of specified words
func (index *FileIndex) WithAll(words ...string) []string {
	return resultsToStrings(index.index.WithAll(words...))
}

// WithAny searches all files in index with at least one of specified words
func (index *FileIndex) WithAny(words ...string) []string {
	return resultsToStrings(index.index.WithAny(words...))
}

// WithoutAny searches all files in index without any of specified words
func (index *FileIndex) WithoutAny(words ...string) []string {
	return resultsToStrings(index.index.WithoutAny(words...))
}