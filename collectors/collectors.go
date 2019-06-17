package collectors

// SliceCollector collects the elements of the stream in a slice
type SliceCollector struct {
	collection []interface{}
}

// NewSliceCollector creates a SliceCollector with an empty
// slice and returns a pointer to it.
func NewSliceCollector() *SliceCollector {
	collector := SliceCollector{[]interface{}{}}
	return &collector
}

// Add adds the element to the slice in SliceCollector,
// expanding the slice if necessary
func (collector *SliceCollector) Add(element interface{}) {
	collector.collection = append(collector.collection, element)
}

// Complete returns the slice that Add has been populating
func (collector *SliceCollector) Complete() interface{} {
	return collector.collection
}

// MapCollector collects Entries into a map. Use mappers.KeyValueMapper
// to map to a stream of collectors.Entry
type MapCollector struct {
	collection map[interface{}]interface{}
}

// NewMapCollector creates a new MapCollector with an empty map
// and returns a pointer to it.
func NewMapCollector() *MapCollector {
	collector := MapCollector{map[interface{}]interface{}{}}
	return &collector
}

// Entry is a basic key value pair used by MapCollector and GroupByCollector.
// It can be created using the mappers.KeyValueMapper function
type Entry struct {
	Key   interface{}
	Value interface{}
}

// Add adds an entry to the map. If there is an existing entry with the
// same key, it will be overwritten.
func (collector *MapCollector) Add(entry interface{}) {
	asEntry := entry.(Entry)
	collector.collection[asEntry.Key] = asEntry.Value
}

// Complete returns the map that Add has populated
func (collector *MapCollector) Complete() interface{} {
	return collector.collection
}

// GroupByCollector collects Entries into a map where the keys are
// the keys of the Entries and the values are a slice of all the values
// for that key. Use mappers.KeyValueMapper to map to a stream of
// collectors.Entry
type GroupByCollector struct {
	collection map[interface{}][]interface{}
}

// NewGroupByCollector creates a new GroupByCollector with an empty map
// and returns a pointer to it.
func NewGroupByCollector() *GroupByCollector {
	collector := GroupByCollector{map[interface{}][]interface{}{}}
	return &collector
}

// Add adds the entry to the GroupByCollector. If a matching key already exists,
// the value is appended to the existing slice.
func (collector *GroupByCollector) Add(entry interface{}) {
	asEntry := entry.(Entry)
	if value, exists := collector.collection[asEntry.Key]; exists {
		collector.collection[asEntry.Key] = append(value, asEntry.Value)
	} else {
		collector.collection[asEntry.Key] = []interface{}{asEntry.Value}
	}
}

// Complete returns the underlying map that Add has populated.
func (collector *GroupByCollector) Complete() interface{} {
	return collector.collection
}
