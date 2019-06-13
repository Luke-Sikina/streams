# Streams - Higher Order Functions for Streams
[![Build Status](https://travis-ci.org/Luke-Sikina/streams.svg?branch=master)](https://travis-ci.org/Luke-Sikina/streams)

Functions for creating and using Streams. Good for performing map / filter / reduce
operations on large datasets.

## Purpose
I was messing around with large datasets, and I couldn't find a lot of good tools out there.
I found some good design patterns for making iterators with goroutines and channels, but nothing
on using those streams. This library is designed to help you manipulate large data sets in a an
asynchronous fashion
using well understood components of logic.
I borrowed a lot of nomenclature from Java's Streams API because that's what I am familiar with.
## Features
### FromCollection
Creates a Streams object from a slice
### FromStream
Creates a Streams object from a Stream, which is just a `chan interface{}`
### Filter
Filters the elements of a Stream according to the Predicate function.
### Map
Transforms the elements of a Stream according to the Mapper function.
### Reduce
Reduces the elements of a Stream according to the Reducer function.
### Collect
Collects the elements of a Stream according to the Collector object, which collects
elements using its `Add` function and returns the collection using the `Complete` function
### ForEach
Applies the Consumer function to each element of the stream.
### ForEachThen
Applies the Consumer function to each element of the stream. Returns the stream with all
the elements still in it for future use.
## Examples
I strongly recommend you look at the unit and integration tests, as those are examples that
you can trust to compile and work.

### Basic Trivial Use Case
``` Golang

func DivisibleByTwo(subject interface{}) bool {
	return subject.(int)%2 == 0
}

func TimesFive(element interface{}) interface{} {
	return element.(int)*5
}

func Sum(first, second interface{}) interface{} {
	return first.(int) + second.(int)
}

func TrivialExample() {
	stream := FromCollection([]interface{}{1,2,3,4})
	someSum := stream.
		Filter(DivisibleByTwo).
		Map(TimesFive).
		Reduce(0, Sum)

	log.Printf("Aaah! Numbers! %v", someSum)
}

```