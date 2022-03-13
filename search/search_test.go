package search

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
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

	chunkID, err := seaHashID("chunk")
	require.NoError(t, err)

	latchID, err := seaHashID("latch")
	require.NoError(t, err)

	index := newIndex()
	index['a'] = []uint64{latchID}
	index['c'] = []uint64{chunkID, latchID}
	index['h'] = []uint64{chunkID, latchID}
	index['u'] = []uint64{chunkID}
	index['n'] = []uint64{chunkID}
	index['k'] = []uint64{chunkID}
	index['l'] = []uint64{latchID}
	index['t'] = []uint64{latchID}

	assert.Equal(t, db.index, index)
	assert.Equal(t, db.reverseIndex[chunkID], "chunk")
	assert.Equal(t, db.reverseIndex[latchID], "latch")
}

func TestXXHashedIndexedWordDB(t *testing.T) {

	db, err := NewIndexedDB([]string{"chunk", "latch"}, xxHashID)
	require.NoError(t, err)

	chunkID, err := xxHashID("chunk")
	require.NoError(t, err)

	latchID, err := xxHashID("latch")
	require.NoError(t, err)

	index := newIndex()
	index['a'] = []uint64{latchID}
	index['c'] = []uint64{chunkID, latchID}
	index['h'] = []uint64{chunkID, latchID}
	index['u'] = []uint64{chunkID}
	index['n'] = []uint64{chunkID}
	index['k'] = []uint64{chunkID}
	index['l'] = []uint64{latchID}
	index['t'] = []uint64{latchID}

	assert.Equal(t, db.index, index)
	assert.Equal(t, db.reverseIndex[chunkID], "chunk")
	assert.Equal(t, db.reverseIndex[latchID], "latch")
}

// XXHashIndexedWordDB-8   	1000000000	         0.5508 ns/op
// SeaHashIndexedWordDB-8   1000000000	         0.5473 ns/op
func BenchmarkTestSeaHashIndexedWordDB(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	size := 1000000
	var words = make([]string, size)
	for i := 0; i < size; i++ {
		wordle := randomWordle()
		words[i] = wordle.letters
	}

	db, err := NewIndexedDB(words, xxHashID)
	require.NoError(b, err)
	b.Logf("index contains %v items", db.size)
	assert.Equal(b, db.size, size)
	PrintMemUsage()
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

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
