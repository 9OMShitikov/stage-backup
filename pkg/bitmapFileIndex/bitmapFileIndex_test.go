package bitmapFileIndex

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func removeFiles (names []string) error {
	for _, filename := range names {
		err := os.Remove(filename)
		if err != nil {
			return fmt.Errorf("error while removing %s: %w", filename, err)
		}
	}
	return nil
}

func createTestFiles (contents ...string) (map[string] string, []string, error) {
	namesContents := make(map[string] string)
	names := make([]string, 0, len(contents))
	var err error = nil
	for _, content := range contents {
		var testFile *os.File
		testFile, err = ioutil.TempFile("", "tmp_test_file")
		if err != nil {
			break
		}

		namesContents[testFile.Name()] = content
		names = append(names, testFile.Name())

		_, err = testFile.WriteString(content)
		if err != nil {
			break
		}

		err = testFile.Close()
		if err != nil {
			break
		}
	}
	if err != nil {
		for _, filename := range names {
			os.Remove(filename)
		}
		return nil, nil, fmt.Errorf("error while creating temp files: %w" , err)
	}
	return namesContents, names, nil
}

//Check if slices l and contents[r] contain equal sets of strings and contains[r] doesn't contain duplicates
func equalStringSetsWithRWithoutDup(l []string, r []string, contents map[string] string) bool {
	set := make(map[string]struct{})

	for _, elementL := range l {
		set[elementL] = struct{}{}
	}

	for _, elementR := range r {
		contentR, ok := contents[elementR]

		if !ok {
			return false
		}

		if _, ok := set[contentR]; ok {
			delete(set, contentR)
		} else {
			return false
		}
	}
	if len(set) != 0 {
		return false
	}
	return true
}

var simpleTestContents = [] string {
	"|a b c|",
	"|a b|",
	"|a c|",
	"|b c|",
	"|a|",
	"|b|",
	"|c|",
	"|e|",
}

func TestFileIndex_WithAll_Simple (t *testing.T) {
	var tests = []struct {
		query   []string
		results []string
	}{
		{[]string{"a"}, []string{"|a|", "|a c|", "|a b|", "|a b c|"}},
		{[]string{"a", "b"}, []string{"|a b|", "|a b c|"}},
		{[]string{"a", "b", "c"}, []string{"|a b c|"}},
		{[]string{"d"}, []string{}},
		{[]string{"a", "d"}, []string{}},
		{[]string{"e"}, []string{"|e|"}},
		{[]string{"a", "e"}, []string{}},
		{[]string{}, []string{}},
	}

	filesMap, filesNames, err := createTestFiles(simpleTestContents...)
	if err != nil {
		t.Error("test failed while creating files: ", err)
	}

	filesIndex, err := CreateFileIndex(filesNames...)
	if err != nil {
		t.Error("test failed while creating index: ", err)
	}

	for _, test := range tests {
		testName := fmt.Sprint("with all from ", test.query, ":")
		t.Run(testName, func(t *testing.T) {
			ans := filesIndex.WithAll(test.query...)
			if !equalStringSetsWithRWithoutDup(test.results, ans, filesMap) {
				t.Error("expected ", test.results, ", got ", ans, " at map:\n", filesMap)
			}
		})
	}

	err = removeFiles(filesNames)
	if err != nil {
		t.Error("test failed while removing files: ", err)
	}
}

func TestFileIndex_WithAny_Simple (t *testing.T) {
	var tests = []struct {
		query   []string
		results []string
	}{
		{[]string{"a"}, []string{"|a|", "|a c|", "|a b|", "|a b c|"}},
		{[]string{"a", "b"}, []string{"|a|", "|a c|", "|a b|", "|a b c|",
			"|b|", "|b c|"}},
		{[]string{"a", "b", "c"}, []string{"|a|", "|a c|", "|a b|", "|a b c|",
			"|b|", "|b c|", "|c|"}},
		{[]string{"d"}, []string{}},
		{[]string{"a", "d"}, []string{"|a|", "|a c|", "|a b|", "|a b c|"}},
		{[]string{"e"}, []string{"|e|"}},
		{[]string{"a", "e"}, []string{"|a|", "|a c|", "|a b|", "|a b c|", "|e|"}},
		{[]string{}, []string{}},
	}

	filesMap, filesNames, err := createTestFiles(simpleTestContents...)
	if err != nil {
		t.Error("test failed while creating files: ", err)
	}

	filesIndex, err := CreateFileIndex(filesNames...)
	if err != nil {
		t.Error("test failed while creating index: ", err)
	}

	for _, test := range tests {
		testName := fmt.Sprint("with any from ", test.query, ":")
		t.Run(testName, func(t *testing.T) {
			ans := filesIndex.WithAny(test.query...)
			if !equalStringSetsWithRWithoutDup(test.results, ans, filesMap) {
				t.Error("expected ", test.results, ", got ", ans, " at map:\n", filesMap)
			}
		})
	}

	err = removeFiles(filesNames)
	if err != nil {
		t.Error("test failed while removing files: ", err)
	}
}

func TestFileIndex_WithoutAny_Simple (t *testing.T) {
	var tests = []struct {
		query   []string
		results []string
	}{
		{[]string{"a"}, []string{"|b|", "|c|", "|b c|", "|e|"}},
		{[]string{"a", "b"}, []string{"|c|", "|e|"}},
		{[]string{"a", "b", "c"}, []string{"|e|"}},
		{[]string{"d"}, []string{"|a|", "|b|", "|c|",
			"|a b|", "|a c|", "|b c|", "|a b c|", "|e|"}},
		{[]string{"a", "d"}, []string{"|b|", "|c|", "|b c|", "|e|"}},
		{[]string{"e"}, []string{"|a|", "|b|", "|c|",
			"|a b|", "|a c|", "|b c|", "|a b c|"}},
		{[]string{"a", "e"}, []string{"|b|", "|c|", "|b c|"}},
		{[]string{}, []string{"|a|", "|b|", "|c|",
			"|a b|", "|a c|", "|b c|", "|a b c|", "|e|"}},
	}

	filesMap, filesNames, err := createTestFiles(simpleTestContents...)
	if err != nil {
		t.Error("test failed while creating files: ", err)
	}

	filesIndex, err := CreateFileIndex(filesNames...)
	if err != nil {
		t.Error("test failed while creating index: ", err)
	}

	for _, test := range tests {
		testName := fmt.Sprint("without any from ", test.query, ":")
		t.Run(testName, func(t *testing.T) {
			ans := filesIndex.WithoutAny(test.query...)
			if !equalStringSetsWithRWithoutDup(test.results, ans, filesMap) {
				t.Error("expected ", test.results, ", got ", ans, " at map:\n", filesMap)
			}
		})
	}

	err = removeFiles(filesNames)
	if err != nil {
		t.Error("test failed while removing files: ", err)
	}
}

func TestCreateFileIndex_error (t *testing.T) {
	_, err := CreateFileIndex("./not_existent_file")
	if err == nil {
		t.Error("no errors catched")
	}
}