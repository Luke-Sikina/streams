package streams

type Predicate func(element interface{}) bool
type Reducer func(first, second interface{}) interface{}
type Mapper func(element interface{}) interface{}
type Consumer func(element interface{})

// Used by Streams.Collect
// Add is used to add stream elements to the collection
// Complete returns the collection
type Collector interface {
	Add(subject interface{})
	Complete() interface{}
}

type Stream chan interface{}
type Streams struct {
	Streams       []Stream
	channelBuffer int
}

func toStream(collection []interface{}) Stream {
	ch := make(chan interface{}, len(collection))

	go func() {
		defer close(ch)
		for _, element := range collection {
			ch <- element
		}
	}()

	return ch
}

// Creates a Streams object from the given slice.
// The channel buffer size will be set to the size of the slice
// If the data being processed is large enough that a slice would be
// impractical, use FromStream instead
func FromCollection(collection []interface{}) Streams {
	startStream := toStream(collection)
	return FromStream(startStream, len(collection))
}

// Creates a Streams object from the given channel
// Future Stream objects in the Streams object will be created with
// a buffer size of bufferSize
func FromStream(stream Stream, bufferSize int) Streams {
	return Streams{[]Stream{stream}, bufferSize}
}

func addNewStream(streams *Streams) (current, next Stream) {
	current = streams.Streams[len(streams.Streams)-1]
	next = make(Stream, streams.channelBuffer)
	streams.Streams = append(streams.Streams, next)
	return
}

// Asynchronously filters the elements in the streams using the provided Predicate.
// Elements that cause the Predicate to evaluate to true are kept,
// elements that cause the Predicate to evaluate to false are discarded.
func (streams *Streams) Filter(predicate Predicate) *Streams {
	current, next := addNewStream(streams)

	go func() {
		defer close(next)
		for object := range current {
			if predicate(object) {
				next <- object
			}
		}
	}()

	return streams
}

// Asynchronously transforms the elements in the streams using the provided Mapper.
// Use this to turn the elements of the stream from one thing into another thing
func (streams *Streams) Map(mapper Mapper) *Streams {
	current, next := addNewStream(streams)

	go func() {
		defer close(next)
		for object := range current {
			next <- mapper(object)
		}
	}()

	return streams
}

func (streams *Streams) lastStream() Stream {
	return streams.Streams[len(streams.Streams)-1]
}

// Reduce the elements in the stream to a singe element
// The single element can be of a different element, but it should
// probably be the same type as initial, or else your Reducer function
// will be ugly.
// The first parameter will always be the reduction thus far, and for the first
// iteration, it will be the initial parameter passed into Reduce
func (streams *Streams) Reduce(initial interface{}, reducer Reducer) interface{} {
	for element := range streams.lastStream() {
		initial = reducer(initial, element)
	}
	return initial
}

// Collect the elements in the stream in a single collection like a slice or
// a map. This function is out of place because it accepts an interface
// instead of a function.
// Calls Collector.Add on each element of the stream, then returns Collector.Complete
func (streams *Streams) Collect(collector Collector) interface{} {
	lastStream := streams.lastStream()
	for element := range lastStream {
		collector.Add(element)
	}
	return collector.Complete()
}

// Calls consumer(element) on each element on the stream
func (streams *Streams) ForEach(consumer Consumer) {
	lastStream := streams.lastStream()
	for element := range lastStream {
		consumer(element)
	}
}

// Calls consumer(element) on each element of the stream, returns the streams
// for future use.
// As far as I can tell, there is no instance where it is more efficient
// to do ForEachThen().ForEach() rather than combining the consumer functions
// That said, this can be more readable, and it allows you to Collect / Reduce after
func (streams *Streams) ForEachThen(consumer Consumer) *Streams {
	current, next := addNewStream(streams)

	go func() {
		defer close(next)
		for element := range current {
			consumer(element)
			next <- element
		}
	}()

	return streams
}
