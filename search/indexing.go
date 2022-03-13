package search

import (
	"hash"

	"blainsmith.com/go/seahash"
	"github.com/cespare/xxhash"
	"github.com/go-logr/logr"
)

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

var UseXXHashID IDFn = NewHashingIDFn(xxhash.New)
var UseSeaHashID IDFn = NewHashingIDFn(func() hash.Hash64 {
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
