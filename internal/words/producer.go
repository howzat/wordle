package words

import "github.com/go-logr/logr"

type Producer struct {
	wordsStream chan Result
	logger      *logr.Logger
}

func NewProducer(log *logr.Logger, wordsStream chan Result) Producer {
	return Producer{
		wordsStream: wordsStream,
		logger:      log,
	}
}

func (p *Producer) Produce(readWords ReadWordsFn, filterFunc ...FilterFunc) {
	words, err := readWords(filterFunc...)
	if err != nil {
		p.wordsStream <- Failure(err)
	} else {
		p.wordsStream <- Success(words)
	}
}
