package mappers

import (
	"fmt"
	"github.com/Luke-Sikina/streams"
	"github.com/Luke-Sikina/streams/collectors"
	"strconv"
)

// IntToStringMapper converts elements from int to string
func IntToStringMapper(element interface{}) interface{} {
	return strconv.Itoa(element.(int))
}

// FloatToStringMapper converts elements from float to string
func FloatToStringMapper(element interface{}) interface{} {
	return fmt.Sprintf("%f", element)
}

// StringToIntMapper converts elements from string to int
// Elements that cannot be converted result in the 0 value being added instead.
// To avoid this, it may be worthwhile to use streams.Filter first.
func StringToIntMapper(element interface{}) interface{} {
	asInt, err := strconv.Atoi(element.(string))

	if err == nil {
		return asInt
	} else {
		return 0
	}
}

// StringToFloatMapper returns a function that converts elements
// from string to float of floatSize (either 32 or 64). Elements that
// cannot be converted result in a mapping to the float 0.0.
// To avoid this, it may be worthwhile to use streams.Filter first.
func StringToFloatMapper(floatSize int) streams.Mapper {
	return func(element interface{}) interface{} {
		asFloat, err := strconv.ParseFloat(element.(string), floatSize)

		if err == nil {
			return asFloat
		} else {
			return 0.0
		}
	}
}

// Gets the key and value from an element
type EntryCreator func(element interface{}) (interface{}, interface{})

// KeyValueMapper takes a EntryCreator function. It uses this function to create
// a collectors.Entry, which it then returns. The function can be used to create
// a stream of collectors.Entry for collectors.GroupByCollector and collectors.MapCollector
func KeyValueMapper(toEntry EntryCreator) streams.Mapper {
	return func(element interface{}) interface{} {
		key, value := toEntry(element)
		return collectors.Entry{Key: key, Value: value}
	}
}
