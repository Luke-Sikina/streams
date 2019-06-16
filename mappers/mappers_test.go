package mappers

import (
	"github.com/Luke-Sikina/streams"
	"github.com/Luke-Sikina/streams/collectors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MapperCase struct {
	Start    []interface{}
	Expected []interface{}
}

func TestIntToStringMapper(t *testing.T) {
	cases := []MapperCase{
		{
			[]interface{}{},
			[]interface{}{},
		}, {
			[]interface{}{1, 2, 3},
			[]interface{}{"1", "2", "3"},
		},
	}

	for _, caze := range cases {
		stream := streams.FromCollection(caze.Start)
		actual := stream.
			Map(IntToStringMapper).
			Collect(collectors.NewSliceCollector())

		assert.Equal(t, caze.Expected, actual)
	}
}

func TestFloatToStringMapper(t *testing.T) {
	cases := []MapperCase{
		{
			[]interface{}{},
			[]interface{}{},
		}, {
			[]interface{}{1.0, 2.0, 3.0},
			[]interface{}{"1.000000", "2.000000", "3.000000"},
		},
	}

	for _, caze := range cases {
		stream := streams.FromCollection(caze.Start)
		actual := stream.
			Map(FloatToStringMapper).
			Collect(collectors.NewSliceCollector())

		assert.Equal(t, caze.Expected, actual)
	}
}

func TestStringToIntMapper(t *testing.T) {
	cases := []MapperCase{
		{
			[]interface{}{},
			[]interface{}{},
		}, {
			[]interface{}{"1", "2", "3"},
			[]interface{}{1, 2, 3},
		}, {
			[]interface{}{"1", "2.0", "three"},
			[]interface{}{1, 0, 0},
		},
	}

	for _, caze := range cases {
		stream := streams.FromCollection(caze.Start)
		actual := stream.
			Map(StringToIntMapper).
			Collect(collectors.NewSliceCollector())

		assert.Equal(t, caze.Expected, actual)
	}
}

type FLoatToIntMapperCase struct {
	Start     []interface{}
	FloatSize int
	Expected  []interface{}
}

func TestFloatToIntMapper(t *testing.T) {
	cases := []FLoatToIntMapperCase{
		{
			[]interface{}{},
			64,
			[]interface{}{},
		}, {
			[]interface{}{"1.25", "2.25", "3.25"},
			64,
			[]interface{}{1.25, 2.25, 3.25},
		}, {
			[]interface{}{"1.25", "2.25", "3.25"},
			32,
			[]interface{}{1.25, 2.25, 3.25},
		}, {
			[]interface{}{"1.25", "foo"},
			64,
			[]interface{}{1.25, 0.0},
		},
		{
			[]interface{}{"1.25", "foo"},
			32,
			[]interface{}{1.25, 0.0},
		},
	}

	for _, caze := range cases {
		stream := streams.FromCollection(caze.Start)
		actual := stream.
			Map(StringToFloatMapper(caze.FloatSize)).
			Collect(collectors.NewSliceCollector())

		assert.Equal(t, caze.Expected, actual)
	}
}
