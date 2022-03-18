package db

import (
	"hash"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/cespare/xxhash"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

type Index struct {
	size         int
	reverseIndex map[uint64]string
	index        map[string][]uint64
}

type IDFn = func(string) (uint64, error)

func NewIndex(_ logr.Logger, words []string, idFn IDFn) (*Index, error) {
	index := map[string][]uint64{}
	reverseIndex := make(map[uint64]string, len(words))
	var recall = map[string]bool{}
	for _, lw := range words {
		w := strings.ToLower(lw)
		if _, ok := recall[w]; !ok {
			recall[w] = true
		} else {
			continue
		}

		id, err := idFn(w)
		if err != nil {
			return nil, err
		}

		if a, ok := reverseIndex[id]; ok {
			if a != w {
				return nil, errors.Errorf("hash collision between %v and %v", a, w)
			}
		}

		reverseIndex[id] = w
		for _, c := range w {
			index[string(c)] = append(index[string(c)], id)
		}
	}

	return &Index{
		size:         len(reverseIndex),
		reverseIndex: reverseIndex,
		index:        index,
	}, nil
}

func NewHashingIDFn(hr func() hash.Hash64) IDFn {
	return func(s string) (uint64, error) {
		h := hr()
		_, err := h.Write([]byte(s))
		return h.Sum64(), err
	}
}

var UseXXHashID IDFn = NewHashingIDFn(xxhash.New)

var Alphabet = []string{"a",
	"b",
	"c",
	"d",
	"f",
	"g",
	"h",
	"i",
	"j",
	"k",
	"l",
	"m",
	"n",
	"o",
	"p",
	"q",
	"r",
	"s",
	"t",
	"u",
	"v",
	"w",
	"x",
	"y",
	"z",
}

func (d Index) PickRandomWord() string {
	rand.Seed(time.Now().Unix())
	firstAlpha := Alphabet[rand.Intn(25)]
	ids := d.index[firstAlpha]
	id := ids[rand.Intn(len(ids))]
	return d.reverseIndex[id]
}

func (d Index) Search(guess Wordle) (*MatchResult, error) {

	letters := guess.AllKnownLetters()
	var candidateIds []uint64
	for _, letter := range letters {
		candidateIds = append(candidateIds, d.index[letter]...)
	}

	var recall = map[string]bool{}
	var candidateResults []string
	for _, id := range candidateIds {
		candidateWord := d.reverseIndex[id]
		if _, ok := recall[candidateWord]; !ok {
			recall[candidateWord] = true // we've processed this word before
			if containsAllKnownLetters(letters, candidateWord) &&
				fullyKnownLettersAreInCorrectPosition(guess, candidateWord) {
				candidateResults = append(candidateResults, candidateWord)
			}
		} else {
			continue
		}
	}

	sort.Strings(candidateResults)
	return &MatchResult{
		Items: candidateResults,
		Guess: guess,
	}, nil
}

func fullyKnownLettersAreInCorrectPosition(wordle Wordle, letters string) bool {
	if len(wordle.FullyKnownLetters()) == 0 {
		return false
	}

	for i, k := range wordle.knowledge {
		if k == Full {
			if letters[i] != wordle.letters[i] {
				return false
			}
		}
	}
	return true
}

func containsAllKnownLetters(letters []string, word string) bool {
	for _, letter := range letters {
		if !strings.Contains(word, letter) {
			return false
		}
	}

	return true
}

func (d Index) CandidateGuess(wordle string) (*Wordle, error) {
	var candidateGuess string
	for guess := ""; len(guess) == 0; guess = candidateGuess {
		candidate := d.PickRandomWord()
		var knowledge = BuildKnowledgeForGuess(wordle, candidate)
		search, err := NewWordleSearch(candidate, knowledge)
		if err != nil {
			return nil, err
		}

		if len(search.AllKnownLetters()) > 0 {
			return search, nil
		}
	}
	return nil, nil
}

func BuildKnowledgeForGuess(wordle string, guess string) []Knowlege {

	wb := []byte(wordle)
	var k = []Knowlege{None, None, None, None, None}
	for i, char := range guess {
		pos := findChar(wb, byte(char))
		if pos >= 0 {
			if pos == i {
				k[i] = Full
			} else {
				k[i] = Present
			}
		}
	}
	return k
}

func findChar(wordle []byte, char byte) int {
	for i, b := range wordle {
		if b == char {
			return i
		}
	}
	return -1
}
