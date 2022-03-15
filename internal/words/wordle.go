package words

import (
	"context"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-logr/logr"
)

func CompileWordList(ctx context.Context, log *logr.Logger, wordSource WordSources) (*WordList, error) {

	ctx, done := context.WithCancel(ctx)
	consumer := NewConsumer(log)
	go consumer.Consume(ctx)

	producer := NewProducer(log, consumer.AddWordsStream)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		englishWordsDictionary := ParseLineSeperatedDictionary(wordSource.EnglishWordFile)
		producer.Produce(englishWordsDictionary, WordleCandidate, ChangeToLowerCase)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		localWordsList := ParseLineSeperatedDictionary(wordSource.LocalWordFiles[0])
		producer.Produce(localWordsList, WordleCandidate, ChangeToLowerCase)
	}()

	for _, wordFile := range wordSource.WordSetFiles {
		wg.Add(1)
		fp := wordFile
		go func() {
			defer wg.Done()
			producer.Produce(ParseWordsetDictionary(fp), WordleCandidate, ChangeToLowerCase)
		}()
	}

	wg.Wait()
	done()

	words := consumer.ListWords()

	return &WordList{
		Ingested: len(words),
		Words:    words,
	}, nil
}

type WordList struct {
	Ingested int
	Words    []string
}

func NewWordSourceFiles(baseDir string) (*WordSources, error) {
	wordSources := WordSources{baseDir: baseDir}
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

	directory, err := wordSources.filepath("wordset-dictionary/data/")
	wordSetDirectory, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, WrapErr(err, "could not list contents of directory [%v]", directory)
	}

	wsf := make([]string, len(wordSetDirectory))
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

func (c WordSources) filepath(path string) (string, error) {
	fullpath := c.baseDir + "/" + path
	_, err := os.Stat(fullpath)
	return fullpath, err
}
