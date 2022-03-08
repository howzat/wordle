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

	compiled, err := CompileWordList(context.TODO(), log, CompileConfig{
		Outfile:             "wordle-words.json",
		WordsetDataDir:      "../../dictionaries/wordset-dictionary/data",
		EnglishWordsDataDir: "../../dictionaries/english-words",
	})

	assert.NoError(t, err)
	assert.Equal(t, 478243, compiled.ingested)
}
