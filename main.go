package main

import (
	"context"
	"go.uber.org/zap"
	"io/ioutil"
	"runtime"
	"sync"

	"github.com/howzat/wordle/model"
	"github.com/howzat/wordle/wordset"
)

const WordsetDataDir = "dictionaries/wordset-dictionary/data"

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx, cancelContextFn := context.WithCancel(context.Background())

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer func(logger *zap.Logger) { _ = logger.Sync() }(logger)

	consumer := model.NewWordsConsumer(logger)
	go consumer.Start(ctx)

	logger.Info("reading words from Wordset dictionary", zap.String("path", WordsetDataDir))
	files, err := ioutil.ReadDir(WordsetDataDir)
	if err != nil {
		panic(model.WrapErr(err, "could not list contents of directory [%v]", WordsetDataDir))
	}

	var wg sync.WaitGroup

	for _, file := range files {
		file := file
		go func() {
			wg.Add(1)
			producer := wordset.NewWordsetWordsProducer(WordsetDataDir, consumer.IngestChan, logger)
			producer.ReadWordsetFile(file)
		}()
	}

	go func() {
		consumer.Consume(&wg)
	}()

	wg.Wait()
	cancelContextFn() // Signal cancellation to context.Context and shutdown the consumer
	logger.Info("complete", zap.Int("ingested", len(consumer.ListWords())))
}
