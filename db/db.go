package db

import (
	"reflect"

	"github.com/pkg/errors"
)

type WordSearchEngine interface {
	Search(guess Wordle) (*MatchResult, error)
}

type Wordle struct {
	letters   string
	knowledge []Knowlege
}

func (w Wordle) FullyKnownLetters() []string {
	return w.filterKnowledgeBy(func(knowledge Knowlege) bool {
		return knowledge == Full
	})
}

func (w Wordle) KnownLetters() []string {
	return w.filterKnowledgeBy(func(knowledge Knowlege) bool {
		return knowledge != None
	})
}

func (w Wordle) filterKnowledgeBy(f func(knowledge Knowlege) bool) []string {
	var known []string
	for i, k := range w.knowledge {
		if f(k) {
			known = append(known, string(w.letters[i]))
		}
	}
	return known
}

type Knowlege int8

const (
	Full Knowlege = iota
	Present
	None
)

type MatchResult struct {
	Items []string
	Guess Wordle
}

func NewWordleSearch(letters string, knowledge []Knowlege) (*Wordle, error) {
	if len(letters) > 5 {
		return nil, errors.New("guesses must have exactly 5 characters")
	}
	if len(letters) > 5 {
		return nil, errors.New("knowledge must have exactly 5 items")
	}

	return &Wordle{
		letters:   letters,
		knowledge: knowledge,
	}, nil
}

func NewSearchEngine(db *Index) WordSearchEngine {
	return &LocalSearchEngine{
		words: db,
	}
}

type LocalSearchEngine struct {
	words *Index
}

var NoKnowledge = []Knowlege{None, None, None, None, None}

func (ws *LocalSearchEngine) Search(guess Wordle) (*MatchResult, error) {

	if guess.knowledge == nil || reflect.DeepEqual(guess.knowledge, NoKnowledge) {
		return nil, errors.New("searching without search will match the entire dictionary")
	}

	var results []string
	for i, fact := range guess.knowledge {
		if fact == Present {
			l := guess.letters[i]
			ids := ws.words.index[string(l)]
			for _, id := range ids {
				word := ws.words.reverseIndex[id]
				results = append(results, word)
			}
		}
	}

	return &MatchResult{
		Items: results,
		Guess: guess,
	}, nil
}
