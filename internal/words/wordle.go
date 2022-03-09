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
	ctx, done := context.WithCancel(ctx)
	consumer := NewConsumer(log)
	go consumer.Consume(ctx)

	producer := NewProducer(log, consumer.AddWordsStream)

	files, err := ioutil.ReadDir(config.WordsetDataDir)
	if err != nil {
		done()
		return nil, WrapErr(err, "could not list contents of directory [%v]", config.WordsetDataDir)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fp := filepath("words_alpha.txt", config.EnglishWordsDataDir)
		englishWordsDictionary := ParseEnglishWordsDictionary(fp)
		producer.Produce(englishWordsDictionary, WordleCandidate)
	}()

	for _, wordFile := range files {
		wg.Add(1)
		f := wordFile
		go func() {
			defer wg.Done()
			fp := filepath(f.Name(), config.WordsetDataDir)
			producer.Produce(ParseWordsetDictionary(fp), WordleCandidate)
		}()
	}

	wg.Wait()
	done()

	words := consumer.ListWords()

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
