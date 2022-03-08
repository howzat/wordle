package words

import (
	"regexp"

	"github.com/go-logr/logr"
)

func NewWordsConsumer(log *logr.Logger) *WordConsumer {
	return &WordConsumer{
		logger:     log,
		Words:      make([]string, 0),
		IngestChan: make(chan Result, 1),
	}
}

type WordConsumer struct {
	IngestChan chan Result
	Words      []string
	logger     *logr.Logger
}

/*
Start is always terminated on the side of the producer
*/
func (c *WordConsumer) Start() {
	for {
		select {
		case event := <-c.IngestChan:
			if event.Err != nil {
				c.logger.Error(event.Err, "error result received")
			} else if event.Done {
				return
			} else {
				c.Words = append(c.Words, event.Words...)
			}
		default:
		}
	}
}

func (c *WordConsumer) ListWords() []string {
	return c.Words
}

func (c *WordConsumer) Stop() {
	c.logger.Info("stopping consumer")
	c.IngestChan <- Done()
}

type Result struct {
	Done  bool
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

func Done() Result {
	return Result{
		Done: true,
	}
}

type FilterFunc = func(e string) bool

var filter FilterFunc = func(e string) bool {
	matches, _ := regexp.MatchString(`^[a-zA-Z]+$`, e)
	return matches
}

type WordsProducer struct {
	ingestChan chan Result
	logger     *logr.Logger
}

func NewWordsProducer(ingestChan chan Result, log *logr.Logger) WordsProducer {
	return WordsProducer{
		logger:     log,
		ingestChan: ingestChan,
	}
}

type ReadWordsFn = func() ([]string, error)
