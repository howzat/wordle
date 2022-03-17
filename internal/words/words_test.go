package words

import (
	"context"
	"testing"

	"github.com/howzat/wordle/internal/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompileWordList(t *testing.T) {

	log, err := logging.NewProductionLogger("TestCompileWordList")
	require.NoError(t, err)

	wordSource, err := NewWordSources("../../dictionaries")
	require.NoError(t, err)

	assert.Equal(t, 27, len(wordSource.WordSetFiles))

	compiled, err := wordSource.LoadWords(context.TODO(), log)

	assert.NoError(t, err)
	assert.Equal(t, 31941, compiled.Size)
}
