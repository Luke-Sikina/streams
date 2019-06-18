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

// ConsumeWithDelimitedWriter returns a stream.Consumer function that, when called
// writes the element to the writer, with the delimiter appended. This means you
// have an extra delimiter at the end of whatever you're writing to.
func ConsumeWithDelimitedWriter(writer *bufio.Writer, delimter string) streams.Consumer {
	return func(element interface{}) {
		// Ignoring the err here is not ideal. When the lib becomes more
		// fault tolerant, this should be changed
		_, _ = writer.WriteString(fmt.Sprintf("%v%v", element, delimter))
	}
}
