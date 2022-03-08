package words

import (
	"context"
	"io/ioutil"
	"runtime"
	"sync"

	"github.com/go-logr/logr"
)

type CompileConfig struct {
	Outfile             string
	WordsetDataDir      string
	EnglishWordsDataDir string
}

func CompileWordList(ctx context.Context, log *logr.Logger, config CompileConfig) (*WordList, error) {

	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Info("compiling words list")

	consumer := NewWordsConsumer(log)
	go consumer.Start()

	files, err := ioutil.ReadDir(config.WordsetDataDir)
	if err != nil {
		return nil, WrapErr(err, "could not list contents of directory [%v]", config.WordsetDataDir)
	}

	ingest := consumer.IngestChan
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		readWords(ingest, ParseEnglishWordsDictionary(log, filepath("words_alpha.txt", config.EnglishWordsDataDir)))
	}()

	for _, file := range files { //:465527
		wg.Add(1)
		f := file
		go func() {
			defer wg.Done()
			readWords(ingest, ParseWordsetDictionary(log, filepath(f.Name(), config.WordsetDataDir)))
		}()
	}

	wg.Wait()
	ingest <- Done()

	words := consumer.ListWords()
	log.Info("compiled 4", "ingested", len(words))

	return &WordList{
		ingested: len(words),
		words:    words,
	}, nil
}

func filepath(filename string, baseDir string) string {
	return baseDir + "/" + filename
}

type WordList struct {
	ingested int
	words    []string
}

func readWords(ingest chan Result, readWords ReadWordsFn) {
	words, err := readWords()
	if err != nil {
		ingest <- Failure(err)
	} else {
		ingest <- Success(words)
	}
}
