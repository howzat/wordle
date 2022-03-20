package db

import (
	"fmt"
	"hash"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/howzat/wordle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashingIndexingWordDB(t *testing.T) {
	log, err := wordle.NewProductionLogger(t.Name())
	require.NoError(t, err)

	db, err := NewIndex(*log, []string{"chunk", "latch", "LATCH", "Latch"}, UseXXHashID)
	require.NoError(t, err)

	assert.Equal(t, 2, len(db.reverseIndex))
	assert.Equal(t, 8, len(db.index)) // c,h,u,n,k,l,a,t
}

func TestHashingIndexingCollisionsWordDB(t *testing.T) {
	log, err := wordle.NewProductionLogger(t.Name())
	require.NoError(t, err)

	_, err = NewIndex(*log, []string{"chunk", "latch", "LATCH", "Latch"}, UseFixedHasher)
	assert.EqualError(t, err, "hash collision between chunk and latch")
}

func TestHashingConsistencyForIndexedWordDB(t *testing.T) {
	t.Run("test indexing with xxHash", testIndexingWithHasher(UseXXHashID))
	t.Run("test indexing with seaHash", testIndexingWithHasher(UseXXHashID)) //TODO: use other hasher
}

func testIndexingWithHasher(id IDFn) func(b *testing.T) {
	return func(t *testing.T) {
		log, err := wordle.NewProductionLogger(t.Name())
		require.NoError(t, err)

		db, err := NewIndex(*log, []string{"chunk", "latch"}, id)
		require.NoError(t, err)

		chunkID, err := id("chunk")
		require.NoError(t, err)

		latchID, err := id("latch")
		require.NoError(t, err)

		index := map[string][]uint64{}
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
	b.Run("seaHash hasher", testHasher(UseXXHashID)) //TODO Use other hasher
}

func testHasher(id IDFn) func(b *testing.B) {
	return func(b *testing.B) {
		log, err := wordle.NewProductionLogger(b.Name())
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

		db, err := NewIndex(*log, words, id)
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

var UseFixedHasher IDFn = NewHashingIDFn(func() hash.Hash64 {
	var h hash.Hash64 = &fixedHasher{}
	return h
})

type fixedHasher struct{ p []byte }

func (f *fixedHasher) Sum64() uint64       { return uint64(1) }
func (f *fixedHasher) Sum(b []byte) []byte { return b }
func (f *fixedHasher) Reset()              {}
func (f *fixedHasher) BlockSize() int      { return 1 }
func (f *fixedHasher) Size() int           { return len(f.p) }
func (f *fixedHasher) Write(p []byte) (n int, err error) {
	f.p = p
	return len(p), nil
}
