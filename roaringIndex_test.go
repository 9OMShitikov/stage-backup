package roaringIndex

import (
	"fmt"
	"testing"
)

//Check if slices l and r contain equal sets of strings and r doesn't contain duplicates
func equalStringSetsWithRAsInterfacesWithoutDup(l []string, r []interface{}) bool {
	set := make(map[string]struct{})

	for _, elementL := range l {
		set[elementL] = struct{}{}
	}

	for _, elementR := range r {
		elementRStr := elementR.(string)
		if _, ok := set[elementRStr]; ok {
			delete(set, elementRStr)
		} else {
			return false
		}
	}
	if len(set) != 0 {
		return false
	}
	return true
}

var simpleTestObjects = []IndexedObject{
	CreateIndexedObject("|a b c|", "a", "b", "c"),
	CreateIndexedObject("|a b|", "a", "b"),
	CreateIndexedObject("|a c|", "a", "c"),
	CreateIndexedObject("|b c|", "b", "c"),
	CreateIndexedObject("|a|", "a"),
	CreateIndexedObject("|b|", "b"),
	CreateIndexedObject("|c|", "c"),
	CreateIndexedObject("|e|", "e"),
}

func TestIndex_WithAll_Simple(t *testing.T) {
	test_index := CreateIndexFrom(simpleTestObjects...)

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

	for _, test := range tests {
		testname := fmt.Sprint("with all from ", test.query, ":")
		t.Run(testname, func(t *testing.T) {
			ans := test_index.WithAll(test.query...)
			if !equalStringSetsWithRAsInterfacesWithoutDup(test.results, ans) {
				t.Error("expected ", test.results, ", got ", ans)
			}
		})
	}
}

func TestIndex_WithAny_Simple(t *testing.T) {
	test_index := CreateIndexFrom(simpleTestObjects...)

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

	for _, test := range tests {
		testname := fmt.Sprint("with any from ", test.query, ":")
		t.Run(testname, func(t *testing.T) {
			ans := test_index.WithAny(test.query...)
			if !equalStringSetsWithRAsInterfacesWithoutDup(test.results, ans) {
				t.Error("expected ", test.results, ", got ", ans)
			}
		})
	}
}

func TestIndex_WithoutAny_Simple(t *testing.T) {
	test_index := CreateIndexFrom(simpleTestObjects...)

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

	for _, test := range tests {
		testname := fmt.Sprint("without any of ", test.query, ":")
		t.Run(testname, func(t *testing.T) {
			ans := test_index.WithoutAny(test.query...)
			if !equalStringSetsWithRAsInterfacesWithoutDup(test.results, ans) {
				t.Error("expected ", test.results, ", got ", ans)
			}
		})
	}
}
