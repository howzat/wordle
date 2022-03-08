package words

import "regexp"

type FilterFunc = func(e string) bool

var filter FilterFunc = func(e string) bool {
	matches, _ := regexp.MatchString(`^[a-zA-Z]+$`, e)
	return matches
}

func NoFilter(e string) bool {
	return true
}

type ReadWordsFn = func(filter ...FilterFunc) ([]string, error)
