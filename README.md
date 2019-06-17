# Streams - Higher Order Functions for Streams
[![Build Status](https://travis-ci.org/Luke-Sikina/streams.svg?branch=master)](https://travis-ci.org/Luke-Sikina/streams)
[![GoDoc](https://godoc.org/github.com/Luke-Sikina/streams?status.svg)](https://godoc.org/github.com/Luke-Sikina/streams)
[![Go Report Card](https://goreportcard.com/badge/github.com/Luke-Sikina/streams)](https://goreportcard.com/report/github.com/Luke-Sikina/streams)

Functions for creating and using Streams. Good for performing map / filter / reduce
operations on datasets you can't keep entirely in memory.

## Purpose
I was messing around with large datasets, and I couldn't find a lot of good tools out there.
I found some good design patterns for making iterators with goroutines and channels, but nothing
on using those streams. This library is designed to help you manipulate large data sets using
well understood components that you've seen in other languages.
I borrowed a lot of nomenclature from Java's Streams API because that's what I am familiar with.
## Features
### streams
This is the main package. It contains the key functions for manipulating streams:
`Filter`, `Map`, `Reduce`, `Collect`, `ForEach`, and `ForEachThen`. It also
has several functions for creating a Streams object: `FromCollection` `FromStream`,
and `FromScanner`

### mappers
This package contains some common helpful mappers. Mappers are functions that match
the `streams.Mapper` function signature and can be passed to the `streams.Map` function.
This is also a good place to look if you're trying to understand how to write
your own Mapper function.

### collectors
This package contains some common helpful collectors. Collectors are structs that
implement the `streams.Collector` interface and can be used in the `streams.Collect`
function. This is also a good place to look if you're trying to understand how to
write your own Collector.

### filters
This package has one helpful filter function. Filter functions are functions that
match the `streams.Predicate` function signature and can be passed to `streams.Filter`
to filter a stream. Since filters are easy to write, and often task specific, this
package is very sparse. This is also a good place to look if you're trying to understand
how to write your own Predicate.

### consumers
This package contains some helpful consumers. Consumers are functions that match the Consumer
function signature; they can be passed to streams.ForEach and streams.ForEachThen. This is
also a good place to look if you're trying to understand how to write your own Consumer.

## Examples
I strongly recommend you look at the unit and integration tests, as those are examples that
you can trust to compile and work. In addition to the tests, here are some examples to get you going.

### Basic Trivial Use Case
This example builds a stream from a collection of ints, filters out the ints not divisible by 2,
then multiplies them each 5 five, and sums the elements together. The example is designed to be compact
and readable so you can quickly understand what this library does.
``` Golang
func TrivialExample() {
	stream := FromCollection([]interface{}{1,2,3,4})
	someSum := stream.
		Filter(DivisibleByTwo).
		Map(TimesFive).
		Reduce(0, Sum)

	log.Printf("Aaah! Numbers! %v", someSum)
}

func DivisibleByTwo(subject interface{}) bool {
	return subject.(int)%2 == 0
}

func TimesFive(element interface{}) interface{} {
	return element.(int)*5
}

func Sum(first, second interface{}) interface{} {
	return first.(int) + second.(int)
}

```

### Read from File, Write to Different File
This example reads a stream of numbers from a file, filters out the positive even numbers, then writes the results
to a second file. There's a lot of boilerplate opening files and setting up bufio objects, so if you're looking
for the meat of the example, start at `stream := streams.FromScanner(scanner, 1023)`.
``` Golang
	// open the input and output files
	input, err := os.Open("input.txt")
	defer closeFileLogErr(input)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	output, err := os.Create("output.txt")
	defer closeFileLogErr(output)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	// set up the scanner and writer
	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)

	// actually use this library to process the data
	stream := streams.FromScanner(scanner, 1023)
	stream.
		Map(mappers.StringToIntMapper).
		Filter(func(e interface{}) bool {return e.(int) % 2 == 1}).
		Map(mappers.IntToStringMapper).
		ForEach(consumers.ConsumeWithDelimitedWriter(writer, "\n"))

	// make sure to flush the writer so you actually get output.
	err = writer.Flush()
	if err != nil {
		log.Fatalf("Error flushing buffer: %v", err)
	}
```