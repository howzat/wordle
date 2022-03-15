package search

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/howzat/wordle/internal/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashingConsistencyForIndexedWordDB(t *testing.T) {
	t.Run("test indexing with xxHash", testIndexingWithHasher(UseXXHashID))
	t.Run("test indexing with seaHash", testIndexingWithHasher(UseSeaHashID))
}

func testIndexingWithHasher(id IDFn) func(b *testing.T) {
	return func(t *testing.T) {
		log, err := logging.NewProductionLogger(t.Name())
		require.NoError(t, err)

		db, err := NewIndexedDB(*log, []string{"chunk", "latch"}, id)
		require.NoError(t, err)

		chunkID, err := id("chunk")
		require.NoError(t, err)

		latchID, err := id("latch")
		require.NoError(t, err)

		index := newAlphaMap()
		index["a"] = []uint64{latchID}
		index["c"] = []uint64{chunkID, latchID}
		index["h"] = []uint64{chunkID, latchID}
		index["u"] = []uint64{chunkID}
		index["n"] = []uint64{chunkID}
		index["k"] = []uint64{chunkID}
		index["l"] = []uint64{latchID}
		index["t"] = []uint64{latchID}

		assert.Equal(t, db.index, index)
		assert.Equal(t, db.reverseIndex[chunkID], "chunk")
		assert.Equal(t, db.reverseIndex[latchID], "latch")
	}
}

// XXHashIndexedWordDB-8   	1000000000	         0.5508 ns/op
// SeaHashIndexedWordDB-8   1000000000	         0.5473 ns/op
func BenchmarkTestHashingForIndexedWordDB(b *testing.B) {
	b.Run("xxHash hasher", testHasher(UseXXHashID))
	b.Run("seaHash hasher", testHasher(UseSeaHashID))
}

func testHasher(id IDFn) func(b *testing.B) {
	return func(b *testing.B) {
		log, err := logging.NewProductionLogger(b.Name())
		require.NoError(b, err)

		rand.Seed(time.Now().UnixNano())
		size := 1000000
		var words = make([]string, size)
		var deduped = map[string]bool{}
		for i := 0; i < size; i++ {
			letters := randomString()
			words[i] = letters
			deduped[letters] = true
		}

		db, err := NewIndexedDB(*log, words, id)
		require.NoError(b, err)
		b.Logf("index contains %v items", db.size)
		assert.Equal(b, db.size, len(deduped))
		PrintMemUsage()
	}
}

func randomString() string {
	b := make([]rune, 5)
	for i := range b {
		b[i] = rune(Alphabet[rand.Intn(len(Alphabet))][0])
	}

	return string(b)
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
