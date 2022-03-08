package words

import (
	"context"
	"sync"

	"github.com/go-logr/logr"
)

func NewWordsConsumer(log *logr.Logger) *WordConsumer {
	return &WordConsumer{
		logger:         log,
		words:          make([]string, 0),
		AddWordsStream: make(chan Result, 1),
	}
}

type WordConsumer struct {
	AddWordsStream chan Result
	readWordsLock  sync.Mutex
	words          []string
	logger         *logr.Logger
}

/*
Consume is always terminated on the side of the producer
*/
func (c *WordConsumer) Consume(ctx context.Context) {

	c.readWordsLock.Lock()
	for {
		select {
		case event := <-c.AddWordsStream:
			if event.Err != nil {
				c.logger.Error(event.Err, "error result received")
			} else {
				c.words = append(c.words, event.Words...)
			}
		case _ = <-ctx.Done():
			c.readWordsLock.Unlock()
			return
		default:
		}
	}
}

func (c *WordConsumer) ListWords() []string {
	c.readWordsLock.Lock()
	defer c.readWordsLock.Unlock()
	return c.words
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
