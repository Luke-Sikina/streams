package filters

// NotNilPredicate filters out all nil elements from the stream
func NotNilPredicate(element interface{}) bool {
	return element != nil
}
