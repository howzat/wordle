package words

import (
	"context"
	"go.uber.org/zap"
	"regexp"
	"sync"
)

func NewWordsConsumer(log *zap.Logger) *WordConsumer {
	return &WordConsumer{
		logger:          log,
		Words:           make([]string, 0),
		AddWordsChannel: make(chan Result, 1),
		IngestChan:      make(chan Result, 1),
	}
}

type WordConsumer struct {
	IngestChan      chan Result
	AddWordsChannel chan Result
	Words           []string
	logger          *zap.Logger
}

/*
Start is always terminated on the side of the producer
*/
func (c *WordConsumer) Start(ctx context.Context) {
	for {
		select {
		case result := <-c.IngestChan: // intermediate channel that allows the producer to stop sending
			c.AddWordsChannel <- result
		case _ = <-ctx.Done():
			close(c.AddWordsChannel)
			return
		}
	}
}

func (c *WordConsumer) Consume(wg *sync.WaitGroup) {
	for event := range c.AddWordsChannel {
		wg.Done()
		if event.Err != nil {
			c.logger.Error("error result received", zap.Error(event.Err))
		} else {
			c.Words = append(c.Words, event.Words...)
		}
	}
}

func (c *WordConsumer) ListWords() []string {
	return c.Words
}

type Result struct {
	Err   error
	Words []string
}

func (r *Result) HasError() bool {
	return r.Err != nil
}

func Success(w []string) Result {
	return Result{
		Err:   nil,
		Words: w,
	}
}

func Failure(e error) Result {
	return Result{
		Err:   e,
		Words: nil,
	}
}

type FilterFunc = func(e string) bool

var filter = func(e string) bool {
	matches, _ := regexp.MatchString(`^[a-zA-Z]+$`, e)
	return matches
}

type WordsProducer struct {
	ingestChan chan Result
	logger     *zap.Logger
}

func NewWordsProducer(ingestChan chan Result, log *zap.Logger) WordsProducer {
	return WordsProducer{
		logger:     log,
		ingestChan: ingestChan,
	}
}

func (w WordsProducer) Work(readWords ReadWordsFn) {
	words, err := readWords()
	if err != nil {
		w.ingestChan <- Failure(err)
	} else {
		w.ingestChan <- Success(words)
	}
}

type ReadWordsFn = func() ([]string, error)
