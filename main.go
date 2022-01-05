package main

import (
	"context"
	"go.uber.org/zap"
	"io/ioutil"
	"runtime"
	"sync"

	"github.com/howzat/wordle/words"
)

const WordsetDataDir = "dictionaries/wordset-dictionary/data"
const EnglishWordsDataDir = "dictionaries/english-words"

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx, cancelContextFn := context.WithCancel(context.Background())

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer func(logger *zap.Logger) { _ = logger.Sync() }(log)

	consumer := words.NewWordsConsumer(log)
	go consumer.Start(ctx)

	var wg sync.WaitGroup

	files, err := ioutil.ReadDir(WordsetDataDir)
	if err != nil {
		panic(words.WrapErr(err, "could not list contents of directory [%v]", WordsetDataDir))
	}

	for _, file := range files {
		file := file
		go func() {
			wg.Add(1)
			producer := words.NewWordsProducer(consumer.IngestChan, log)
			producer.Work(words.ParseWordsetDictionary(log, filepath(file.Name(), WordsetDataDir)))
		}()
	}

	go func() {
		wg.Add(1)
		producer := words.NewWordsProducer(consumer.IngestChan, log)
		producer.Work(words.ParseEnglishWordsDictionary(log, filepath("words_alpha.txt", EnglishWordsDataDir)))
	}()

	go func() {
		consumer.Consume(&wg)
	}()

	wg.Wait()
	cancelContextFn() // Signal cancellation to context.Context and shutdown the consumer
	log.Info("complete", zap.Int("ingested", len(consumer.ListWords())))
}

func filepath(filename string, baseDir string) string {
	return baseDir + "/" + filename
}
