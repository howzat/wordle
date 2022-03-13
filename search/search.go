package search

import (
	"hash"
	"reflect"

	"blainsmith.com/go/seahash"
	"github.com/cespare/xxhash"
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

func NewSearchEngine(words []string) SearchEngine {
	return &LocalSearchEngine{
		words: words,
	}
}

type LocalSearchEngine struct {
	words []string
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

func NewIndexedDB(ingested []string, idFn IDFn) (*IndexedDB, error) {

	idx := newIndex()
	ridx := make(map[uint64]string, len(ingested))

	for _, w := range ingested {
		id, err := idFn(w)
		if err != nil {
			return nil, err
		}

		ridx[id] = w
		for _, letter := range []rune(w) {
			idx[letter] = append(idx[letter], id)
		}
	}

	return &IndexedDB{
		size:         len(ingested),
		reverseIndex: ridx,
		index:        idx,
	}, nil
}

type IndexedDB struct {
	size         int
	reverseIndex map[uint64]string
	index        map[rune][]uint64
}

type IDFn = func(string) (uint64, error)

func NewHashingIDFn(h hash.Hash64) IDFn {
	return func(s string) (uint64, error) {
		_, err := h.Write([]byte(s))
		return h.Sum64(), err
	}
}

var seaHashID IDFn = NewHashingIDFn(seahash.New())
var xxHashID IDFn = NewHashingIDFn(xxhash.New())

func newIndex() map[rune][]uint64 {
	return map[rune][]uint64{
		'a': {},
		'b': {},
		'c': {},
		'd': {},
		'f': {},
		'g': {},
		'h': {},
		'i': {},
		'j': {},
		'k': {},
		'l': {},
		'm': {},
		'n': {},
		'o': {},
		'p': {},
		'q': {},
		'r': {},
		's': {},
		't': {},
		'u': {},
		'v': {},
		'w': {},
		'x': {},
		'y': {},
		'z': {},
	}
}
