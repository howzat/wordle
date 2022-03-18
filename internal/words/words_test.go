package words

import (
	"context"
	"os"
	"testing"

	"github.com/howzat/wordle/internal/logging"
	"github.com/sethvargo/go-envconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Config struct {
	BaseDir string `env:"DICTIONARY_DIR"`
}

func TestCompileWordList(t *testing.T) {

	log, err := logging.NewProductionLogger("TestCompileWordList")
	require.NoError(t, err)

	require.NotEmpty(t, os.Getenv("DICTIONARY_DIR"))

	var c Config
	err = envconfig.Process(context.TODO(), &c)
	require.NoError(t, err)

	wordSource, err := NewWordSources(c.BaseDir)
	require.NoError(t, err)

	assert.Equal(t, 27, len(wordSource.WordSetFiles))

	compiled, err := wordSource.LoadWords(context.TODO(), log)

	assert.NoError(t, err)
	assert.Equal(t, 31941, compiled.Size)
}
