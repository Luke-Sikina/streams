package filters

func NotNilPredicate(element interface{}) bool {
	return element != nil
}
