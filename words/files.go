package words

import (
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"strings"
)

type WordsetFile = map[string]WordsetDictionaryEntry

// WordsetDictionaryEntry wrapper struct to hold additional attributes if required
type WordsetDictionaryEntry struct {
	Word string `json:"word"`
}

func ParseWordsetDictionary(log *zap.Logger, filepath string) ReadWordsFn {
	return func() ([]string, error) {
		log.Info("reading dictionary", zap.String("source", filepath))
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

func ParseEnglishWordsDictionary(log *zap.Logger, filepath string) ReadWordsFn {
	return func() ([]string, error) {
		log.Info("reading dictionary", zap.String("source", filepath))
		file, err := ioutil.ReadFile(filepath)
		if err != nil {
			return nil, WrapErr(err, "error reading file contents", filepath)
		}

		return strings.Split(string(file), "\n"), nil
	}
}
