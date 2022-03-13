package search

import (
	"hash"
	"reflect"

	"blainsmith.com/go/seahash"
	"github.com/cespare/xxhash"
	"github.com/go-logr/logr"
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

func NewIndexedDB(log logr.Logger, words []string, idFn IDFn) (*IndexedDB, error) {
	log.WithName("IndexDB")
	index := newAlphaMap()
	reverseIndex := make(map[uint64]string, len(words))

	for _, w := range words {
		id, err := idFn(w)
		if err != nil {
			return nil, err
		}

		if a, ok := reverseIndex[id]; ok {
			if a != w {
				log.Info("possible collision between %v and %v\n", a, w)
			}
		}

		reverseIndex[id] = w
		for _, c := range []rune(w) {
			index[c] = append(index[c], id)
		}
	}

	return &IndexedDB{
		size:         len(reverseIndex),
		reverseIndex: reverseIndex,
		index:        index,
	}, nil
}

type IndexedDB struct {
	size         int
	reverseIndex map[uint64]string
	index        map[rune][]uint64
}

type IDFn = func(string) (uint64, error)

func NewHashingIDFn(hr func() hash.Hash64) IDFn {
	return func(s string) (uint64, error) {
		h := hr()
		_, err := h.Write([]byte(s))
		return h.Sum64(), err
	}
}

var xxHashID IDFn = NewHashingIDFn(xxhash.New)
var seaHashID IDFn = NewHashingIDFn(func() hash.Hash64 {
	var h hash.Hash64 = seahash.New()
	return h
})

func newAlphaMap() map[rune][]uint64 {
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
