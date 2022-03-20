package wordgen

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/howzat/wordle"
)

var WordleCandidate FilterFn = inOrder(Length(5), Alphabetical(), NoFilter())

func inOrder(fn FilterFn, fns ...FilterFn) FilterFn {
	if len(fns) == 0 {
		return fn
	}

	return inOrder(andThen(fn, fns[0]), fns[1:]...)
}

func andThen(y FilterFn, f FilterFn) FilterFn {
	return func(e string) bool {
		return y(e) && f(e)
	}
}

var NormaliseWord MutatorFn = mInOrder(TrimSurroundingWhitespace, LowercaseWord)

var TrimSurroundingWhitespace MutatorFn = func(s string) string {
	return strings.TrimSpace(s)
}

var LowercaseWord MutatorFn = func(s string) string {
	return strings.ToLower(s)
}

type FilterFn = func(e string) bool

type MutatorFn = func(e string) string

func mInOrder(fn MutatorFn, fns ...MutatorFn) MutatorFn {
	if len(fns) == 0 {
		return fn
	}

	return mInOrder(mAndThen(fn, fns[0]), fns[1:]...)
}

func mAndThen(y MutatorFn, f MutatorFn) MutatorFn {
	return func(e string) string {
		return y(f(e))
	}
}

func NoFilter() FilterFn {
	return func(e string) bool {
		return true
	}
}

func Alphabetical() FilterFn {
	return func(e string) bool {
		isLetter := regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
		return isLetter(e)
	}
}

func Length(l int) FilterFn {
	return func(e string) bool {
		return len(e) == l
	}
}

type ReadWordsFn = func(mutate MutatorFn, filter FilterFn) ([]string, error)

type WordsetFile = map[string]WordsetDictionaryEntry

// WordsetDictionaryEntry wrapper struct to hold additional attributes if required
type WordsetDictionaryEntry struct {
	Word string `json:"word"`
}

func ParseWordsetDictionary(filepath string) ReadWordsFn {
	return func(mutate MutatorFn, filter FilterFn) ([]string, error) {
		f, err := ioutil.ReadFile(filepath)

		if err != nil {
			return nil, wordle.WrapErr(err, "error reading filepath [%v]", filepath)
		}

		var ws WordsetFile
		err = json.Unmarshal(f, &ws)
		if err != nil {
			return nil, wordle.WrapErr(err, "error unmarshalling JSON from filepath [%v]", filepath)
		}

		if err != nil {
			return nil, wordle.WrapErr(err, "error adding words from filepath [%v]", filepath)
		}

		var words []string
		for word := range ws {
			normalised := mutate(word)
			if filter(normalised) {
				words = append(words, normalised)
			}
		}
		return words, nil
	}
}

func ParseLineSeperatedDictionary(filepath string) ReadWordsFn {
	return func(mutate MutatorFn, filter FilterFn) ([]string, error) {
		fileReader, err := os.Open(filepath)
		if err != nil {
			return nil, wordle.WrapErr(err, "error reading file contents [%v]", filepath)
		}

		scanner := bufio.NewScanner(fileReader)
		scanner.Split(bufio.ScanLines)
		var words []string
		for scanner.Scan() {
			word := scanner.Text()
			normalised := mutate(word)
			if filter(normalised) {
				words = append(words, normalised)
			}
		}

		return words, nil
	}
}
