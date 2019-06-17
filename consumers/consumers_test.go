package consumers

import (
	"bufio"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Luke-Sikina/streams"
)

type TestWriter struct {
	Lines [][]byte
}

func (writer *TestWriter) Write(toWrite []byte) (int, error) {
	writer.Lines = append(writer.Lines, toWrite)
	return len(toWrite), nil
}

type ConsumeWithWriterCase struct {
	Start    []interface{}
	Expected [][]byte
}

func TestConsumeWithWriter(t *testing.T) {
	cases := []ConsumeWithWriterCase{
		{
			[]interface{}{},
			[][]byte{},
		}, {
			[]interface{}{"foo", "bar"},
			[][]byte{[]byte("foobar")},
		},
	}

	for _, caze := range cases {
		writer := TestWriter{[][]byte{}}
		bufWriter := bufio.NewWriter(&writer)
		stream := streams.FromCollection(caze.Start)

		stream.ForEach(ConsumeWithWriter(bufWriter))
		_ = bufWriter.Flush()

		assert.Equal(t, caze.Expected, writer.Lines)
	}
}

func TestConsumeWithDelimitedWriter(t *testing.T) {
	cases := []ConsumeWithWriterCase{
		{
			[]interface{}{},
			[][]byte{},
		}, {
			[]interface{}{""},
			[][]byte{[]byte("\n")},
		}, {
			[]interface{}{"foo", "bar"},
			[][]byte{[]byte("foo\nbar\n")},
		},
	}

	for _, caze := range cases {
		writer := TestWriter{[][]byte{}}
		bufWriter := bufio.NewWriter(&writer)
		stream := streams.FromCollection(caze.Start)

		stream.ForEach(ConsumeWithDelimitedWriter(bufWriter, "\n"))
		_ = bufWriter.Flush()

		assert.Equal(t, caze.Expected, writer.Lines)
	}
}
