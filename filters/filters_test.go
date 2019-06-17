package filters

import (
	"github.com/Luke-Sikina/streams"
	"github.com/Luke-Sikina/streams/collectors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type NotNilPredicateCase struct {
	Start    []interface{}
	Expected []interface{}
}

func TestNotNilPredicate(t *testing.T) {
	cases := []NotNilPredicateCase{
		{
			[]interface{}{},
			[]interface{}{},
		}, {
			[]interface{}{nil},
			[]interface{}{},
		}, {
			[]interface{}{1, "a"},
			[]interface{}{1, "a"},
		}, {
			[]interface{}{1, nil, "a", nil},
			[]interface{}{1, "a"},
		},
	}

	for _, caze := range cases {
		stream := streams.FromCollection(caze.Start)
		actual := stream.
			Filter(NotNilPredicate).
			Collect(collectors.NewSliceCollector())

		assert.Equal(t, caze.Expected, actual)
	}
}
