package wordgen

import (
	"context"
	"testing"

	"github.com/howzat/wordle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompileWordList(t *testing.T) {

	ctx := context.TODO()

	log, err := wordle.NewProductionLogger("TestCompileWordList")
	require.NoError(t, err)

	config, err := NewDictionaryConfig(ctx)
	require.NoError(t, err)

	wordSource, err := NewWordSources(config)
	require.NoError(t, err)

	assert.Equal(t, 27, len(wordSource.WordSetFiles))

	compiled, err := wordSource.LoadWords(ctx, log, NormaliseWord, WordleCandidate)

	assert.NoError(t, err)
	assert.Equal(t, 31941, compiled.Size)
}
