package mappers

import (
	"fmt"
	"github.com/Luke-Sikina/streams"
	"strconv"
)

func IntToStringMapper(element interface{}) interface{} {
	return strconv.Itoa(element.(int))
}

func FloatToStringMapper(element interface{}) interface{} {
	return fmt.Sprintf("%f", element)
}

func StringToIntMapper(element interface{}) interface{} {
	asInt, err := strconv.Atoi(element.(string))

	if err == nil {
		return asInt
	} else {
		return 0
	}
}

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
