package search

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLetterMatches(t *testing.T) {

	tests := []struct {
		name       string
		search     string
		results    []string
		dictionary []string
		knowledge  []MatchType
		err        error
	}{
		{
			name:      "providing no knowledge of the wordle is an error",
			search:    "brink",
			knowledge: nil,
			err:       errors.New("searching without knowledge will match the entire dictionary"),
		}, {
			name:      "providing knowledge that is empty is an error",
			search:    "brink",
			knowledge: []MatchType{None, None, None, None, None},
			err:       errors.New("searching without knowledge will match the entire dictionary"),
		}, {
			name:       "single letter match in the wrong place",
			search:     "audio",
			dictionary: []string{"zebra", "bunks"},
			results:    []string{"zebra"},
			knowledge:  []MatchType{Part, None, None, None, None},
		}, {
			name:       "multiple letter match in the wrong place",
			search:     "later",
			dictionary: []string{"slate", "bunks"},
			results:    []string{"slate"},
			knowledge:  []MatchType{Part, Part, Part, Part, None},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			store := NewSearchEngine(tt.dictionary)
			search, err := NewWordle(tt.search)
			require.NoError(t, err)

			matchResult, err := store.Search(*search)
			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.Equal(t, tt.knowledge, matchResult.LetterMatches)
			}
		})
	}
}

func TestSeaHashedIndexedWordDB(t *testing.T) {

	db, err := NewIndexedDB([]string{"chunk", "latch"}, seaHashID)
	require.NoError(t, err)

	index := newIndex()
	index['a'] = []uint64{2}
	index['c'] = []uint64{1, 2}
	index['h'] = []uint64{1, 2}
	index['u'] = []uint64{1}
	index['n'] = []uint64{1}
	index['k'] = []uint64{1}
	index['l'] = []uint64{2}
	index['t'] = []uint64{2}

	assert.Equal(t, db.index, index)
	assert.Equal(t, db.reverseIndex[1], "chunk")
	assert.Equal(t, db.reverseIndex[2], "latch")
}

func TestIndexedWordDB(t *testing.T) {

	rand.Seed(time.Now().UnixNano())

	size := 1000000
	var words = make([]string, size)
	for i := 0; i < size; i++ {
		wordle := randomWordle()
		words[i] = wordle.letters
	}

	db, err := NewIndexedDB(words, seaHashID)
	require.NoError(t, err)

	assert.Equal(t, db.size, size)
}

var alphabet = []rune("abcdefghijklmnopqrstuvwxyz")

func randomWordle() Wordle {
	b := make([]rune, 5)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return Wordle{
		letters: string(b),
	}
}
