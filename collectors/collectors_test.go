package collectors

import (
	"github.com/Luke-Sikina/streams"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSliceCollector(t *testing.T) {
	cases := [][]interface{}{
		{},
		{1, 2, 3},
		{nil},
		{'a', 'b', 'c'},
	}

	for _, caze := range cases {
		subject := streams.FromCollection(caze)
		actual := subject.Collect(NewSliceCollector())

		assert.Equal(t, caze, actual)
	}
}

type MapCollectorCase struct {
	Entries  []interface{}
	Expected map[interface{}]interface{}
}

func TestMapCollector(t *testing.T) {
	cases := []MapCollectorCase{
		{
			[]interface{}{},
			map[interface{}]interface{}{},
		}, {
			[]interface{}{Entry{"a", 1}, Entry{"b", 2}, Entry{"c", 3}},
			map[interface{}]interface{}{"a": 1, "b": 2, "c": 3},
		}, {
			// Not a good idea to do, but here's what will happen
			[]interface{}{Entry{"a", 1}, Entry{"a", 2}},
			map[interface{}]interface{}{"a": 2},
		},
	}

	for _, caze := range cases {
		subject := streams.FromCollection(caze.Entries)
		actual := subject.Collect(NewMapCollector())

		assert.Equal(t, caze.Expected, actual)
	}
}

type GroupByCollectorCase struct {
	Entries  []interface{}
	Expected map[interface{}][]interface{}
}

func TestGroupByCollector(t *testing.T) {
	cases := []GroupByCollectorCase{
		{
			[]interface{}{},
			map[interface{}][]interface{}{},
		}, {
			[]interface{}{Entry{"a", 1}, Entry{"b", 2}},
			map[interface{}][]interface{}{"a": {1}, "b": {2}},
		}, {
			[]interface{}{Entry{"a", 1}, Entry{"a", 2}},
			map[interface{}][]interface{}{"a": {1, 2}},
		},
	}

	for _, caze := range cases {
		subject := streams.FromCollection(caze.Entries)
		actual := subject.Collect(NewGroupByCollector())

		assert.Equal(t, caze.Expected, actual)
	}
}
