package collectors

type SliceCollector struct {
	collection []interface{}
}

func NewSliceCollector() *SliceCollector {
	collector := SliceCollector{[]interface{}{}}
	return &collector
}

func (collector *SliceCollector) Add(element interface{}) {
	collector.collection = append(collector.collection, element)
}

func (collector *SliceCollector) Complete() interface{} {
	return collector.collection
}

type MapCollector struct {
	collection map[interface{}]interface{}
}

func NewMapCollector() *MapCollector {
	collector := MapCollector{map[interface{}]interface{}{}}
	return &collector
}

type Entry struct {
	Key   interface{}
	Value interface{}
}

func (collector *MapCollector) Add(entry interface{}) {
	asEntry := entry.(Entry)
	collector.collection[asEntry.Key] = asEntry.Value
}

func (collector *MapCollector) Complete() interface{} {
	return collector.collection
}

type GroupByCollector struct {
	collection map[interface{}][]interface{}
}

func NewGroupByCollector() *GroupByCollector {
	collector := GroupByCollector{map[interface{}][]interface{}{}}
	return &collector
}

func (collector *GroupByCollector) Add(entry interface{}) {
	asEntry := entry.(Entry)
	if value, exists := collector.collection[asEntry.Key]; exists {
		collector.collection[asEntry.Key] = append(value, asEntry.Value)
	} else {
		collector.collection[asEntry.Key] = []interface{}{asEntry.Value}
	}
}

func (collector *GroupByCollector) Complete() interface{} {
	return collector.collection
}
