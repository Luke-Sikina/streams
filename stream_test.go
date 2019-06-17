package streams

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

// a lot of tests in here depend on collect to slice working. not ideal
type sliceCollector struct {
	collection []interface{}
}

func (collector *sliceCollector) Add(subject interface{}) {
	collector.collection = append(collector.collection, subject)
}

func (collector *sliceCollector) Complete() interface{} {
	return collector.collection
}

func TestFromCollection(t *testing.T) {
	type SomeStruct struct {
		Foo int
		Bar string
	}
	cases := [][]interface{}{
		{},
		{1, 2, 3},
		{"a", "b", "c"},
		{SomeStruct{1, "a"}, SomeStruct{2, "b"}},
	}

	for _, caze := range cases {
		collector := sliceCollector{[]interface{}{}}
		streams := FromCollection(caze)
		actual := streams.Collect(&collector)

		assert.Equal(t, caze, actual)
	}
}

type FromScannerCase struct {
	Lines    [][]byte
	Expected []interface{}
}

// All this gross boilerplate is to avoid doing a full file io test
// the result is a faster, more focused test, but it comes at the
// expense of anyone reading this boiler plate, so sorry!
type TestReader struct {
	lines [][]byte
	index int
}

func (reader *TestReader) Read(p []byte) (n int, err error) {
	if reader.index >= len(reader.lines) {
		return 0, io.EOF
	}
	// I thought that every action here = 1 scan + text
	// NOT TRUE! Make sure your lines end with a \n
	// I think this is decided by the scan.Scanner.Split fn
	copy(p, reader.lines[reader.index])
	reader.index++
	return len(reader.lines[reader.index-1]), nil
}

func TestFromScanner(t *testing.T) {
	cases := []FromScannerCase{
		{
			[][]byte{},
			[]interface{}{},
		}, {
			[][]byte{[]byte("abc")},
			[]interface{}{"abc"},
		}, {
			[][]byte{[]byte("abc\n"), []byte("def")},
			[]interface{}{"abc", "def"},
		}, {
			[][]byte{[]byte("abc\n"), []byte("\t€")},
			[]interface{}{"abc", "\t€"},
		},
	}

	for _, caze := range cases {
		collector := sliceCollector{[]interface{}{}}
		reader := TestReader{caze.Lines, 0}
		scanner := bufio.NewScanner(&reader)

		subject := FromScanner(scanner, 10)
		actual := subject.Collect(&collector)

		assert.Equal(t, caze.Expected, actual)
	}
}

type StreamsFiltercase struct {
	Start        []interface{}
	FirstFilter  Predicate
	SecondFilter Predicate
	Expected     []interface{}
}

func EvenPredicate(subject interface{}) bool {
	return 0 == subject.(int)%2
}

func OddPredicate(subject interface{}) bool {
	return 1 == (2+(subject.(int)%2))%2
}

func AcceptAllPredicate(_ interface{}) bool {
	return true
}

func RejectAllPredicate(_ interface{}) bool {
	return false
}

func TestStreams_Filter(t *testing.T) {
	cases := []StreamsFiltercase{
		{[]interface{}{1, 2, 3, 4}, AcceptAllPredicate, AcceptAllPredicate, []interface{}{1, 2, 3, 4}},
		{[]interface{}{1, 2, 3, 4}, AcceptAllPredicate, RejectAllPredicate, []interface{}{}},
		{[]interface{}{1, 2, 3, 4}, RejectAllPredicate, RejectAllPredicate, []interface{}{}},
		{[]interface{}{1, 2, 3, 4}, EvenPredicate, AcceptAllPredicate, []interface{}{2, 4}},
		{[]interface{}{1, 2, 3, 4}, OddPredicate, AcceptAllPredicate, []interface{}{1, 3}},
		{[]interface{}{1, 2, 3, 4}, OddPredicate, EvenPredicate, []interface{}{}},
	}

	for _, caze := range cases {
		collector := sliceCollector{[]interface{}{}}
		stream := FromCollection(caze.Start)
		actual := stream.Filter(caze.FirstFilter).
			Filter(caze.SecondFilter).
			Collect(&collector)

		assert.Equal(t, caze.Expected, actual)
	}
}

func MapDoubleVal(subject interface{}) interface{} {
	asInt := subject.(int)
	return asInt * 2
}

func MapToSentance(subject interface{}) interface{} {
	asInt := subject.(int)
	return fmt.Sprintf("And a %d", asInt)
}

type StreamsMapcase struct {
	Before   []interface{}
	Mapper   Mapper
	Expected []interface{}
}

func TestStreams_Map(t *testing.T) {
	cases := []StreamsMapcase{
		{[]interface{}{}, MapDoubleVal, []interface{}{}},
		{[]interface{}{1, 2, 3}, MapDoubleVal, []interface{}{2, 4, 6}},
		{[]interface{}{1, 2, 3}, MapToSentance, []interface{}{"And a 1", "And a 2", "And a 3"}},
	}

	for _, caze := range cases {
		collector := sliceCollector{[]interface{}{}}
		stream := FromCollection(caze.Before)
		actual := stream.
			Map(caze.Mapper).
			Collect(&collector)

		assert.Equal(t, caze.Expected, actual)
	}
}

func ReduceToSum(first, second interface{}) interface{} {
	return first.(int) + second.(int)
}

func ReduceToMax(first, second interface{}) interface{} {
	firstInt := first.(int)
	secondInt := second.(int)

	if firstInt > secondInt {
		return firstInt
	}
	return secondInt
}

func ReduceToString(first, second interface{}) interface{} {
	return fmt.Sprintf("%v%v", first, second)
}

type StreamsReduceCase struct {
	Reducer  Reducer
	Before   []interface{}
	Initial  interface{}
	Expected interface{}
}

func TestStreams_Reduce(t *testing.T) {
	cases := []StreamsReduceCase{
		{ReduceToSum, []interface{}{}, 0, 0},
		{ReduceToSum, []interface{}{0}, 0, 0},
		{ReduceToSum, []interface{}{1, 2, 3}, 0, 6},
		{ReduceToSum, []interface{}{1, 2, 3}, 1, 7},
		{ReduceToMax, []interface{}{1, 2, 3}, -1, 3},
		{ReduceToString, []interface{}{1, 2, 3}, "", "123"},
	}

	for _, caze := range cases {
		stream := FromCollection(caze.Before)
		actual := stream.Reduce(caze.Initial, caze.Reducer)

		assert.Equal(t, actual, caze.Expected)
	}
}

// This is not a good example of how to use ForEach
// Just an easy way to verify that it works

func TestStreams_ForEach(t *testing.T) {
	cases := [][]interface{}{
		{},
		{1, 2, 3},
	}

	for _, caze := range cases {
		seen := make([]interface{}, 0, 0)
		consumer := func(element interface{}) {
			seen = append(seen, element)
		}

		stream := FromCollection(caze)
		stream.ForEach(consumer)

		assert.Equal(t, caze, seen)
	}
}

func TestStreams_ForEachThen(t *testing.T) {
	cases := [][]interface{}{
		{},
		{1, 2, 3},
	}

	for _, caze := range cases {
		collector := sliceCollector{[]interface{}{}}
		seen := make([]interface{}, 0, 0)
		consumer := func(element interface{}) {
			seen = append(seen, element)
		}

		streams := FromCollection(caze)
		actual := streams.
			ForEachThen(consumer).
			Collect(&collector)

		assert.Equal(t, caze, seen)
		assert.Equal(t, caze, actual)
	}
}
