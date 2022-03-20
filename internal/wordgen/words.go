package wordgen

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-logr/logr"
)

func (w *WordSources) LoadWords(ctx context.Context, log *logr.Logger, mutate MutatorFn, filter FilterFn) (*Words, error) {

	ctx, done := context.WithCancel(ctx)
	consumer := NewConsumer(log)
	go consumer.Consume(ctx)

	producer := NewProducer(log, consumer.AddWordsStream)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		englishWordsDictionary := ParseLineSeperatedDictionary(w.EnglishWordFile)
		producer.Produce(englishWordsDictionary, filter, mutate)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		localWordsList := ParseLineSeperatedDictionary(w.LocalWordFiles[0])
		producer.Produce(localWordsList, filter, mutate)
	}()

	for _, wordFile := range w.WordSetFiles {
		wg.Add(1)
		fp := wordFile
		go func() {
			defer wg.Done()
			producer.Produce(ParseWordsetDictionary(fp), filter, mutate)
		}()
	}

	wg.Wait()
	done()

	words := consumer.ListWords()

	return &Words{
		Size:  len(words),
		Words: words,
	}, nil
}

type Words struct {
	Size  int
	Words []string
}

func NewWordSources(config Config) (*WordSources, error) {
	wordSources := WordSources{baseDir: config.BaseDir}
	ewf, err := wordSources.filepath("english-words/words_alpha.txt")
	if err != nil {
		return nil, err
	}

	wordSources.EnglishWordFile = ewf

	ld, err := wordSources.filepath("wordlist.txt")
	if err != nil {
		return nil, err
	}

	wordSources.LocalWordFiles = []string{ld}

	var wsf []string
	directory, err := wordSources.filepath("wordset-dictionary/data/")
	err = filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}

		if filepath.Ext(d.Name()) == ".json" {
			wsf = append(wsf, path)
		}

		return nil
	})

	wordSources.WordSetFiles = wsf
	return &wordSources, nil
}

type WordSources struct {
	WordSetFiles    []string
	EnglishWordFile string
	LocalWordFiles  []string
	baseDir         string
}

func (w WordSources) filepath(path string) (string, error) {
	fullpath := w.baseDir + "/" + path
	_, err := os.Stat(fullpath)
	return fullpath, err
}
