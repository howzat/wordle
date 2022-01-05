package wordset

import (
	"encoding/json"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/howzat/wordle/model"
)

type WordsetFile = map[string]Entry

// Entry wrapper struct to hold additional attributes if required
type Entry struct {
	Word string `json:"word"`
}

type WordsProducer struct {
	ingestChan chan model.Result
	baseDir    string
	logger     *zap.Logger
}

func NewWordsetWordsProducer(baseDir string, ingestChan chan model.Result, log *zap.Logger) WordsProducer {
	return WordsProducer{
		logger:     log,
		ingestChan: ingestChan,
		baseDir:    baseDir,
	}
}

func (w WordsProducer) ReadWordsetFile(filename fs.FileInfo) {
	filepath, err := w.fileLocation(filename)
	if err != nil {
		w.logger.Error("error reading Wordset file",
			zap.String("file", *filepath),
			zap.Error(err),
		)

		w.ingestChan <- model.Failure(err)
	}

	words, err := w.readWords(*filepath)
	if err != nil {
		w.ingestChan <- model.Failure(err)
	} else {
		w.ingestChan <- model.Success(words)
	}
}

func (w WordsProducer) fileLocation(filename fs.FileInfo) (*string, error) {
	filepath := w.baseDir + "/" + filename.Name()
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return nil, model.WrapErr(err, "file does not exist", filepath)
	}
	return &filepath, nil
}

func (w *WordsProducer) readWords(file string) ([]string, error) {
	wordsetFile, err := w.parseWordsetFile(file)
	if err != nil {
		return nil, model.WrapErr(err, "error adding words from file [%v]", file)
	}

	words := make([]string, len(*wordsetFile))
	i := 0
	for word := range *wordsetFile {
		words[i] = word
		i++
	}
	return words, nil
}

func (w *WordsProducer) parseWordsetFile(file string) (*WordsetFile, error) {
	f, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, model.WrapErr(err, "error reading file [%v]", file)
	}
	var ws WordsetFile
	err = json.Unmarshal(f, &ws)
	if err != nil {
		return nil, model.WrapErr(err, "error unmarshalling JSON from file [%v]", file)
	}

	return &ws, nil
}
