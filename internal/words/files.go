package words

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

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

func ParseEnglishWordsDictionary(filepath string) ReadWordsFn {
	return func(ftrs ...FilterFunc) ([]string, error) {
		file, err := ioutil.ReadFile(filepath)
		if err != nil {
			return nil, WrapErr(err, "error reading file contents [%v]", filepath)
		}

		return strings.Split(string(file), "\n"), nil
	}
}
