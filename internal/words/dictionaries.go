package words

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"
)

var WordleCandidate FilterFunc = inOrder(Length(5), Alphabetical(), NoFilter())

type FilterFunc = func(e string) bool

func inOrder(fn FilterFunc, fns ...FilterFunc) FilterFunc {
	if len(fns) == 0 {
		return fn
	}

	return inOrder(andThen(fn, fns[0]), fns[1:]...)
}

func andThen(y FilterFunc, f FilterFunc) FilterFunc {
	return func(e string) bool {
		return y(e) && f(e)
	}
}

func NoFilter() FilterFunc {
	return func(e string) bool {
		return true
	}
}

func Alphabetical() FilterFunc {
	return func(e string) bool {
		matches, _ := regexp.MatchString(`^[a-zA-Z]+$`, e)
		return matches
	}
}

func Length(l int) FilterFunc {
	return func(e string) bool {
		return len(e) == l
	}
}

type ReadWordsFn = func(filter ...FilterFunc) ([]string, error)

type WordsetFile = map[string]WordsetDictionaryEntry

// WordsetDictionaryEntry wrapper struct to hold additional attributes if required
type WordsetDictionaryEntry struct {
	Word string `json:"word"`
}

func ParseWordsetDictionary(filepath string) ReadWordsFn {
	return func(ftrs ...FilterFunc) ([]string, error) {
		f, err := ioutil.ReadFile(filepath)

		if err != nil {
			return nil, WrapErr(err, "error reading filepath [%v]", filepath)
		}

		var ws WordsetFile
		err = json.Unmarshal(f, &ws)
		if err != nil {
			return nil, WrapErr(err, "error unmarshalling JSON from filepath [%v]", filepath)
		}

		if err != nil {
			return nil, WrapErr(err, "error adding words from filepath [%v]", filepath)
		}

		words := make([]string, len(ws))
		i := 0
		for word := range ws {
			words[i] = word
			i++
		}
		return words, nil
	}
}

func ParseLineSeperatedDictionary(filepath string) ReadWordsFn {
	return func(ftrs ...FilterFunc) ([]string, error) {
		file, err := ioutil.ReadFile(filepath)

		if err != nil {
			return nil, WrapErr(err, "error reading file contents [%v]", filepath)
		}

		return strings.Split(string(file), "\n"), nil
	}
}
