package consumers

import (
	"bufio"
	"fmt"

	"github.com/Luke-Sikina/streams"
)

// ConsumeWithWriter returns a stream.Consumer function that, when called
// writes the element to the writer.
func ConsumeWithWriter(writer *bufio.Writer) streams.Consumer {
	return ConsumeWithDelimitedWriter(writer, "")
}

func ConsumeWithDelimitedWriter(writer *bufio.Writer, delimter string) streams.Consumer {
	return func(element interface{}) {
		// Ignoring the err here is not ideal. When the lib becomes more
		// fault tolerant, this should be changed
		_, _ = writer.WriteString(fmt.Sprintf("%v%v", element, delimter))
	}
}
