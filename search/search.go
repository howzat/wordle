package search

import (
	"reflect"

	"github.com/pkg/errors"
)

type SearchEngine interface {
	Search(guess Wordle) (*MatchResult, error)
}

type Wordle struct {
	letters   string
	knowledge []MatchType
}

type MatchType int8

const (
	Full MatchType = iota
	Part
	None
)

type MatchResult struct {
	Items         []string
	Guess         Wordle
	LetterMatches []MatchType
}

func NewWordle(letters string) (*Wordle, error) {
	if len(letters) > 5 {
		return nil, errors.New("guesses must have exactly 5 characters")
	}
	return &Wordle{
		letters: letters,
	}, nil
}

func NewSearchEngine(db *IndexedDB) SearchEngine {
	return &LocalSearchEngine{
		words: db,
	}
}

type LocalSearchEngine struct {
	words *IndexedDB
}

var NoKnowledge = []MatchType{None, None, None, None, None}

func (ws *LocalSearchEngine) Search(guess Wordle) (*MatchResult, error) {

	if guess.knowledge == nil || reflect.DeepEqual(guess.knowledge, NoKnowledge) {
		return nil, errors.New("searching without knowledge will match the entire dictionary")
	}

	return &MatchResult{
		LetterMatches: []MatchType{None, None, None, None, None},
	}, nil
}
